package ucevsoc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEVSOC struct {
	service eebusapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType
}

var _ UCEVSOCInterface = (*UCEVSOC)(nil)

func NewUCEVSOC(service eebusapi.ServiceInterface, eventCB api.EntityEventCallback) *UCEVSOC {
	uc := &UCEVSOC{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEVSOC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeEVStateOfCharge
}

func (e *UCEVSOC) AddFeatures() {
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

func (e *UCEVSOC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.0"),
		"RC1",
		true,
		[]model.UseCaseScenarioSupportType{1})
}

func (e *UCEVSOC) UpdateUseCaseAvailability(available bool) {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.SetUseCaseAvailability(model.UseCaseActorTypeCEM, e.UseCaseName(), available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVSOC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1},
		[]model.FeatureTypeType{model.FeatureTypeTypeMeasurement},
	) {
		return false, nil
	}

	// check for required features
	evMeasurement, err := util.Measurement(e.service, entity)
	if err != nil || evMeasurement == nil {
		return false, eebusapi.ErrFunctionNotSupported
	}

	// check if measurement description contains an element with scope SOC
	if _, err = evMeasurement.GetDescriptionsForScope(model.ScopeTypeTypeStateOfCharge); err != nil {
		return false, err
	}

	return true, nil
}
