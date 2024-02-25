package ucevsecc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	serviceapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEVSECC struct {
	service serviceapi.ServiceInterface

	eventCB api.EventHandlerCB

	validEntityTypes []model.EntityTypeType
}

var _ UCEVSECCInterface = (*UCEVSECC)(nil)

func NewUCEVSECC(service serviceapi.ServiceInterface, eventCB api.EventHandlerCB) *UCEVSECC {
	uc := &UCEVSECC{
		service: service,
		eventCB: eventCB,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEVSE,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEVSECC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeEVSECommissioningAndConfiguration
}

func (e *UCEVSECC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeDeviceDiagnosis,
	}

	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

func (e *UCEVSECC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1, 2})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVSECC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEVSE,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{2},
		[]model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
	) {
		// Workaround for the Porsche Mobile Charger Connect that falsely reports
		// the usecase to be on the EV actor
		if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
			model.UseCaseActorTypeEV,
			e.UseCaseName(),
			[]model.UseCaseScenarioSupportType{2},
			[]model.FeatureTypeType{model.FeatureTypeTypeDeviceDiagnosis},
		) {
			return false, nil
		}

	}

	return true, nil
}
