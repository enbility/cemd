package uccevc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCCEVC struct {
	service eebusapi.ServiceInterface

	eventCB api.EventHandlerCB

	validEntityTypes []model.EntityTypeType
}

var _ UCCEVCInterface = (*UCCEVC)(nil)

func NewUCCEVC(service eebusapi.ServiceInterface, eventCB api.EventHandlerCB) *UCCEVC {
	uc := &UCCEVC{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCCEVC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeCoordinatedEVCharging
}

func (e *UCCEVC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeTimeSeries,
		model.FeatureTypeTypeIncentiveTable,
		model.FeatureTypeTypeElectricalConnection,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

func (e *UCCEVC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCCEVC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{2, 3, 4, 5, 6, 7, 8},
		[]model.FeatureTypeType{
			model.FeatureTypeTypeTimeSeries,
			model.FeatureTypeTypeIncentiveTable,
		},
	) {
		return false, nil
	}

	// check for required features
	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}
	evIncentiveTable, err := util.IncentiveTable(e.service, entity)
	if err != nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	// check if timeseries descriptions contains constraints data
	if _, err = evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints); err != nil {
		return false, err
	}

	// check if incentive table descriptions contains data for the required scope
	if _, err = evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable); err != nil {
		return false, err
	}

	return true, nil
}
