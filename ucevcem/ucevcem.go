package ucevcem

import (
	"github.com/enbility/cemd/api"
	serviceapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEVCEM struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCEVCEMInterface = (*UCEVCEM)(nil)

func NewUCEVCEM(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCEVCEM {
	uc := &UCEVCEM{
		service: service,
		reader:  reader,
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
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	f.AddResultHandler(e)
}

func (e *UCEVCEM) AddUseCase() {
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
func (e *UCEVCEM) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false, api.ErrNoEvEntity
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
