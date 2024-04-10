package ucevcem

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	serviceapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEVCEM struct {
	service serviceapi.ServiceInterface

	eventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType
}

var _ UCEVCEMInterface = (*UCEVCEM)(nil)

func NewUCEVCEM(service serviceapi.ServiceInterface, eventCB api.EntityEventCallback) *UCEVCEM {
	uc := &UCEVCEM{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEVCEM) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging
}

func (e *UCEVCEM) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		_ = localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
	}
}

func (e *UCEVCEM) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"release",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}

func (e *UCEVCEM) UpdateUseCaseAvailability(available bool) {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.SetUseCaseAvailability(model.UseCaseActorTypeCEM, e.UseCaseName(), available)
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVCEM) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName(),
		nil,
		nil,
	) {
		return false, nil
	}

	return true, nil
}
