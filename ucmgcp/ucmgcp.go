package ucmgcp

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

type UCMGCP struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCMGCPInterface = (*UCMGCP)(nil)

func NewUCMGCP(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCMGCP {
	uc := &UCMGCP{
		service: service,
		reader:  reader,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCMGCP) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeMonitoringOfGridConnectionPoint
}

func (e *UCMGCP) AddFeatures() {
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

func (e *UCMGCP) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeMonitoringAppliance,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCMGCP) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !e.isCompatibleEntity(entity) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeGridConnectionPoint,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{2, 3, 4},
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

	_, err1 := measurement.GetDescriptionsForScope(model.ScopeTypeTypeACPower)
	_, err2 := measurement.GetDescriptionsForScope(model.ScopeTypeTypeGridFeedIn)
	_, err3 := measurement.GetDescriptionsForScope(model.ScopeTypeTypeGridConsumption)
	if err1 != nil || err2 != nil || err3 != nil {
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

	return true, nil
}