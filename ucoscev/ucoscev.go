package ucoscev

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

type UCOSCEV struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCOSCEVInterface = (*UCOSCEV)(nil)

func NewUCOSCEV(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCOSCEV {
	uc := &UCOSCEV{
		service: service,
		reader:  reader,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCOSCEV) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging
}

func (e *UCOSCEV) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	f.AddResultHandler(e)

	f = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddResultHandler(e)
}

func (e *UCOSCEV) AddUseCase() {
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
func (e *UCOSCEV) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1, 2, 3},
		[]model.FeatureTypeType{model.FeatureTypeTypeLoadControl},
	) {
		return false, nil
	}

	// check if loadcontrol limit descriptions contains a recommendation category
	evLoadControl, err := util.LoadControl(e.service, entity)
	if err != nil {
		return false, features.ErrFunctionNotSupported
	}

	if _, err = evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeRecommendation); err != nil {
		return false, err
	}

	return true, nil
}
