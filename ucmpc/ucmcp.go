package ucmpc

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

type UCMPC struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface

	validEntityTypes []model.EntityTypeType
}

var _ UCMCPInterface = (*UCMPC)(nil)

func NewUCMCP(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCMPC {
	uc := &UCMPC{
		service: service,
		reader:  reader,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeCompressor,
		model.EntityTypeTypeElectricalImmersionHeater,
		model.EntityTypeTypeEVSE,
		model.EntityTypeTypeHeatPumpAppliance,
		model.EntityTypeTypeInverter,
		model.EntityTypeTypeSmartEnergyAppliance,
		model.EntityTypeTypeSubMeterElectricity,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCMPC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeMonitoringOfPowerConsumption
}

func (e *UCMPC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

func (e *UCMPC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeMonitoringAppliance,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCMPC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeMonitoredUnit,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeElectricalConnection,
			model.FeatureTypeTypeMeasurement,
		},
	) {
		return false, nil
	}

	// check if measurement description contain data for the required scope
	measurement, err := util.Measurement(e.service, entity)
	if err != nil {
		return false, features.ErrFunctionNotSupported
	}

	if _, err := measurement.GetDescriptionsForScope(model.ScopeTypeTypeACPowerTotal); err != nil {
		return false, features.ErrDataNotAvailable
	}

	// check if electrical connection descriptions is provided
	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return false, features.ErrFunctionNotSupported
	}

	if _, err = electricalConnection.GetDescriptions(); err != nil {
		return false, err
	}

	if _, err = electricalConnection.GetParameterDescriptions(); err != nil {
		return false, err
	}

	return true, nil
}
