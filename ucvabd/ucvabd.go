package ucvabd

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	serviceapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCVABD struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCVABDInterface = (*UCVABD)(nil)

func NewUCVABD(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCVABD {
	uc := &UCVABD{
		service: service,
		reader:  reader,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCVABD) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeVisualizationOfAggregatedBatteryData
}

func (e *UCVABD) AddFeatures() {
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

func (e *UCVABD) AddUseCase() {
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
func (e *UCVABD) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.isCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypePVSystem,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1, 4},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeElectricalConnection,
			model.FeatureTypeTypeMeasurement,
		},
	) {
		return false, nil
	}

	// check for required features
	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return false, features.ErrFunctionNotSupported
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
		return false, features.ErrFunctionNotSupported
	}

	// check if measurement descriptions contains a required scope
	if _, err = measurement.GetDescriptionsForScope(model.ScopeTypeTypeACPowerTotal); err != nil {
		return false, err
	}
	if _, err = measurement.GetDescriptionsForScope(model.ScopeTypeTypeStateOfCharge); err != nil {
		return false, err
	}

	return true, nil
}
