package uclpc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCLPC struct {
	service eebusapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType
}

var _ UCLCPInterface = (*UCLPC)(nil)

func NewUCLPC(service eebusapi.ServiceInterface, eventCB api.EntityEventCallback) *UCLPC {
	uc := &UCLPC{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeHeatPumpAppliance,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCLPC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeLimitationOfPowerConsumption
}

func (e *UCLPC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	f.AddResultHandler(e)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	f.AddResultHandler(e)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	f.AddResultHandler(e)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	f.AddResultHandler(e)

	// server features
	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	f.AddResultHandler(e)
}

func (e *UCLPC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (e *UCLPC) InitializeDataStructures() {}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCLPC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
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
			model.FeatureTypeTypeLoadControl,
			model.FeatureTypeTypeDeviceConfiguration,
			model.FeatureTypeTypeElectricalConnection,
		},
	) {
		return false, nil
	}

	if _, err := util.DeviceDiagnosis(e.service, entity); err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	if _, err := util.LoadControl(e.service, entity); err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	if _, err := util.DeviceConfiguration(e.service, entity); err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	if _, err := util.ElectricalConnection(e.service, entity); err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	return true, nil
}
