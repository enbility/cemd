package uclpcserver

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCLPCServer struct {
	service eebusapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType

	heartbeatKeoWorkaround bool // required because KEO Stack uses multiple identical entities for the same functionality, and it is not clear which to use
}

var _ UCLCPServerInterface = (*UCLPCServer)(nil)

func NewUCLPC(service eebusapi.ServiceInterface, eventCB api.EntityEventCallback) *UCLPCServer {
	uc := &UCLPCServer{
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

func (c *UCLPCServer) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeLimitationOfPowerConsumption
}

func (e *UCLPCServer) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)

	// server features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)

	var limitId model.LoadControlLimitIdType = 0
	// get the highest limitId
	if desc, err := spine.LocalFeatureDataCopyOfType[*model.LoadControlLimitDescriptionListDataType](
		f, model.FunctionTypeLoadControlLimitDescriptionListData); err == nil && desc.LoadControlLimitDescriptionData != nil {
		for _, desc := range desc.LoadControlLimitDescriptionData {
			if desc.LimitId != nil && *desc.LimitId >= limitId {
				limitId++
			}
		}
	}

	loadControlDesc := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(limitId)),
				LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeConsume),
				Unit:           eebusutil.Ptr(model.UnitOfMeasurementTypeW),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}
	f.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, loadControlDesc)

	loadControl := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           eebusutil.Ptr(model.LoadControlLimitIdType(limitId)),
				IsLimitChangeable: eebusutil.Ptr(true),
				IsLimitActive:     eebusutil.Ptr(false),
			},
		},
	}
	f.SetData(model.FunctionTypeLoadControlLimitListData, loadControl)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)

	var configId model.DeviceConfigurationKeyIdType = 0
	// get the heighest keyId
	if desc, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueDescriptionListDataType](
		f, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData); err == nil && desc.DeviceConfigurationKeyValueDescriptionData != nil {
		for _, desc := range desc.DeviceConfigurationKeyValueDescriptionData {
			if desc.KeyId != nil && *desc.KeyId >= configId {
				configId++
			}
		}
	}

	deviceConfigDesc := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:     eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId)),
				KeyName:   eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
				ValueType: eebusutil.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
				Unit:      eebusutil.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				KeyId:     eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId + 1)),
				KeyName:   eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
				ValueType: eebusutil.Ptr(model.DeviceConfigurationKeyValueTypeTypeDuration),
			},
		},
	}
	f.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, deviceConfigDesc)

	deviceConfig := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId:             eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId)),
				IsValueChangeable: eebusutil.Ptr(true),
			},
			{
				KeyId:             eebusutil.Ptr(model.DeviceConfigurationKeyIdType(configId + 1)),
				IsValueChangeable: eebusutil.Ptr(true),
			},
		},
	}
	f.SetData(model.FunctionTypeDeviceConfigurationKeyValueListData, deviceConfig)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, true)

	var elCharId model.ElectricalConnectionCharacteristicIdType = 0
	// get the heighest CharacteristicId
	if desc, err := spine.LocalFeatureDataCopyOfType[*model.ElectricalConnectionCharacteristicListDataType](
		f, model.FunctionTypeElectricalConnectionCharacteristicListData); err == nil && desc.ElectricalConnectionCharacteristicData != nil {
		for _, desc := range desc.ElectricalConnectionCharacteristicData {
			if desc.CharacteristicId != nil && *desc.CharacteristicId >= elCharId {
				elCharId++
			}
		}
	}

	// ElectricalConnectionId and ParameterId should be identical to the ones used
	// in a MPC Server role implementation, which is not done here (yet)
	elCharData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       eebusutil.Ptr(elCharId),
				CharacteristicContext:  eebusutil.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     eebusutil.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax),
				Unit:                   eebusutil.Ptr(model.UnitOfMeasurementTypeW),
			},
		},
	}
	f.SetData(model.FunctionTypeElectricalConnectionCharacteristicListData, elCharData)
}

func (e *UCLPCServer) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (e *UCLPCServer) UpdateUseCaseAvailability(available bool) {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.SetUseCaseAvailability(model.UseCaseActorTypeCEM, e.UseCaseName(), available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCLPCServer) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
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
