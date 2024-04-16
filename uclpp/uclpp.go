package uclpp

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCLPP struct {
	service eebusapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType
}

var _ UCLPPInterface = (*UCLPP)(nil)

func NewUCLPP(service eebusapi.ServiceInterface, eventCB api.EntityEventCallback) *UCLPP {
	uc := &UCLPP{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCLPP) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeLimitationOfPowerProduction
}

func (e *UCLPP) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceDiagnosis,
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
	}
	for _, feature := range clientFeatures {
		_ = localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}

	// server features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
}

func (e *UCLPP) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeEnergyGuard,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (e *UCLPP) UpdateUseCaseAvailability(available bool) {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.SetUseCaseAvailability(model.UseCaseActorTypeEnergyGuard, e.UseCaseName(), available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCLPP) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
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

	return true, nil
}
