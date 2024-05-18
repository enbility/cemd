package uclppserver

import (
	"errors"
	"sync"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCLPPServer struct {
	service eebusapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType

	pendingMux    sync.Mutex
	pendingLimits map[model.MsgCounterType]*spineapi.Message

	heartbeatKeoWorkaround bool // required because KEO Stack uses multiple identical entities for the same functionality, and it is not clear which to use
}

var _ UCLPPServerInterface = (*UCLPPServer)(nil)

func NewUCLPP(service eebusapi.ServiceInterface, eventCB api.EntityEventCallback) *UCLPPServer {
	uc := &UCLPPServer{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeGridGuard,
		model.EntityTypeTypeCEM, // KEO uses this entity type for an SMGW whysoever
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCLPPServer) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeLimitationOfPowerProduction
}

func (e *UCLPPServer) loadControlLimitId() (limitid model.LoadControlLimitIdType, err error) {
	limitid = model.LoadControlLimitIdType(0)
	err = errors.New("not found")

	descriptions := util.GetLocalLimitDescriptionsForTypeCategoryDirectionScope(
		e.service,
		model.LoadControlLimitTypeTypeSignDependentAbsValueLimit,
		model.LoadControlCategoryTypeObligation,
		model.EnergyDirectionTypeProduce,
		model.ScopeTypeTypeActivePowerLimit,
	)
	if len(descriptions) != 1 || descriptions[0].LimitId == nil {
		return
	}
	description := descriptions[0]

	if description.LimitId == nil {
		return
	}

	return *description.LimitId, nil
}

// callback invoked on incoming write messages to this
// loadcontrol server feature.
// the implementation only considers write messages for this use case and
// approves all others
func (e *UCLPPServer) loadControlWriteCB(msg *spineapi.Message) {
	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	if msg.RequestHeader == nil || msg.RequestHeader.MsgCounter == nil ||
		msg.Cmd.LoadControlLimitListData == nil {
		return
	}

	limitId, err := e.loadControlLimitId()
	if err != nil {
		return
	}

	data := msg.Cmd.LoadControlLimitListData

	// we assume there is always only one limit
	if data == nil || data.LoadControlLimitData == nil ||
		len(data.LoadControlLimitData) == 0 {
		return
	}

	// check if there is a matching limitId in the data
	for _, item := range data.LoadControlLimitData {
		if item.LimitId == nil ||
			limitId != *item.LimitId {
			continue
		}

		if _, ok := e.pendingLimits[*msg.RequestHeader.MsgCounter]; !ok {
			e.pendingLimits[*msg.RequestHeader.MsgCounter] = msg
			e.eventCB(msg.DeviceRemote.Ski(), msg.DeviceRemote, msg.EntityRemote, WriteApprovalRequired)
			return
		}
	}

	// approve, because this is no request for this usecase
	e.ApproveOrDenyProductionLimit(*msg.RequestHeader.MsgCounter, true, "")
}

func (e *UCLPPServer) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)

	// server features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	_ = f.AddWriteApprovalCallback(e.loadControlWriteCB)

	var limitId model.LoadControlLimitIdType = 0
	// get the highest limitId
	loadControlDesc, err := spine.LocalFeatureDataCopyOfType[*model.LoadControlLimitDescriptionListDataType](
		f, model.FunctionTypeLoadControlLimitDescriptionListData)
	if err == nil && loadControlDesc.LoadControlLimitDescriptionData != nil {
		for _, desc := range loadControlDesc.LoadControlLimitDescriptionData {
			if desc.LimitId != nil && *desc.LimitId >= limitId {
				limitId++
			}
		}
	}

	if loadControlDesc == nil || len(loadControlDesc.LoadControlLimitDescriptionData) == 0 {
		loadControlDesc = &model.LoadControlLimitDescriptionListDataType{
			LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{},
		}
	}

	newLimitDesc := model.LoadControlLimitDescriptionDataType{
		LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(limitId)),
		LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeProduce),
		MeasurementId:  eebusutil.Ptr(model.MeasurementIdType(0)), // This is a fake Measurement ID, as there is no Electrical Connection server defined, it can't provide any meaningful. But KEO requires this to be set :(
		Unit:           eebusutil.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}
	loadControlDesc.LoadControlLimitDescriptionData = append(loadControlDesc.LoadControlLimitDescriptionData, newLimitDesc)
	f.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, loadControlDesc)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)

	var configId model.DeviceConfigurationKeyIdType = 0
	// get the highest keyId
	deviceConfigDesc, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueDescriptionListDataType](
		f, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if err == nil && deviceConfigDesc.DeviceConfigurationKeyValueDescriptionData != nil {
		for _, desc := range deviceConfigDesc.DeviceConfigurationKeyValueDescriptionData {
			if desc.KeyId != nil && *desc.KeyId >= configId {
				configId++
			}
		}
	}

	if deviceConfigDesc == nil || len(deviceConfigDesc.DeviceConfigurationKeyValueDescriptionData) == 0 {
		deviceConfigDesc = &model.DeviceConfigurationKeyValueDescriptionListDataType{
			DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{},
		}
	}

	newConfigs := []model.DeviceConfigurationKeyValueDescriptionDataType{
		{
			KeyId:     eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId)),
			KeyName:   eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			ValueType: eebusutil.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
			Unit:      eebusutil.Ptr(model.UnitOfMeasurementTypeW),
		},
		{
			KeyId:     eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId + 1)),
			KeyName:   eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			ValueType: eebusutil.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
		},
	}
	deviceConfigDesc.DeviceConfigurationKeyValueDescriptionData = append(deviceConfigDesc.DeviceConfigurationKeyValueDescriptionData, newConfigs...)
	f.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, deviceConfigDesc)

	configData, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueListDataType](f, model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil || configData == nil || len(configData.DeviceConfigurationKeyValueData) == 0 {
		configData = &model.DeviceConfigurationKeyValueListDataType{
			DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
		}
	}

	newConfigData := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId:             eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId)),
			IsValueChangeable: eebusutil.Ptr(true),
		},
		{
			KeyId:             eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId + 1)),
			IsValueChangeable: eebusutil.Ptr(true),
		},
	}

	configData.DeviceConfigurationKeyValueData = append(configData.DeviceConfigurationKeyValueData, newConfigData...)
	f.SetData(model.FunctionTypeDeviceConfigurationKeyValueListData, configData)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, false)

	var elCharId model.ElectricalConnectionCharacteristicIdType = 0
	// get the highest CharacteristicId
	elCharData, err := spine.LocalFeatureDataCopyOfType[*model.ElectricalConnectionCharacteristicListDataType](
		f, model.FunctionTypeElectricalConnectionCharacteristicListData)
	if err == nil && elCharData.ElectricalConnectionCharacteristicData != nil {
		for _, desc := range elCharData.ElectricalConnectionCharacteristicData {
			if desc.CharacteristicId != nil && *desc.CharacteristicId >= elCharId {
				elCharId++
			}
		}
	}

	if err != nil || configData == nil || len(configData.DeviceConfigurationKeyValueData) == 0 {
		elCharData = &model.ElectricalConnectionCharacteristicListDataType{
			ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{},
		}
	}

	// ElectricalConnectionId and ParameterId should be identical to the ones used
	// in a MPC Server role implementation, which is not done here (yet)
	newCharData := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicId:       eebusutil.Ptr(elCharId),
		CharacteristicContext:  eebusutil.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     eebusutil.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualProductionNominalMax),
		Unit:                   eebusutil.Ptr(model.UnitOfMeasurementTypeW),
	}
	elCharData.ElectricalConnectionCharacteristicData = append(elCharData.ElectricalConnectionCharacteristicData, newCharData)
	f.SetData(model.FunctionTypeElectricalConnectionCharacteristicListData, elCharData)
}

func (e *UCLPPServer) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeControllableSystem,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (e *UCLPPServer) UpdateUseCaseAvailability(available bool) {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.SetUseCaseAvailability(model.UseCaseActorTypeControllableSystem, e.UseCaseName(), available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCLPPServer) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEnergyGuard,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeDeviceDiagnosis,
		},
	) {
		return false, nil
	}

	if _, err := util.DeviceDiagnosis(e.service, entity); err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	return true, nil
}
