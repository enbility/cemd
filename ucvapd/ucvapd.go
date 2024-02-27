package ucvapd

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCVAPD struct {
	service eebusapi.ServiceInterface

	eventCB api.EventHandlerCB

	validEntityTypes []model.EntityTypeType
}

var _ UCVAPDInterface = (*UCVAPD)(nil)

func NewUCVAPD(service eebusapi.ServiceInterface, eventCB api.EventHandlerCB) *UCVAPD {
	uc := &UCVAPD{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypePVSystem,
	}
	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCVAPD) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData
}

func (e *UCVAPD) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}

}

func (e *UCVAPD) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"RC1",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCVAPD) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypePVSystem,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeDeviceConfiguration,
			model.FeatureTypeTypeElectricalConnection,
			model.FeatureTypeTypeMeasurement,
		},
	) {
		return false, nil
	}

	// check for required features
	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	// check if device configuration descriptions contains a required key name
	if _, err = deviceConfiguration.GetDescriptionForKeyName(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem); err != nil {
		return false, err
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	// check if electrical connection descriptions and parameter descriptions are available name
	if _, err = electricalConnection.GetDescriptions(); err != nil {
		return false, err
	}
	if _, err = electricalConnection.GetParameterDescriptions(); err != nil {
		return false, err
	}

	// check for required features
	measurement, err := util.Measurement(e.service, entity)
	if err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	// check if measurement descriptions contains a required scope
	if _, err = measurement.GetDescriptionsForScope(model.ScopeTypeTypeACPowerTotal); err != nil {
		return false, err
	}
	if _, err = measurement.GetDescriptionsForScope(model.ScopeTypeTypeACYieldTotal); err != nil {
		return false, err
	}

	return true, nil
}
