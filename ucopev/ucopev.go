package ucopev

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	serviceapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCOPEV struct {
	service serviceapi.ServiceInterface

	reader api.EventReaderInterface

	validEntityTypes []model.EntityTypeType
}

var _ UCOPEVInterface = (*UCOPEV)(nil)

func NewUCOPEV(service serviceapi.ServiceInterface, reader api.EventReaderInterface) *UCOPEV {
	uc := &UCOPEV{
		service: service,
		reader:  reader,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCOPEV) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment
}

func (e *UCOPEV) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeElectricalConnection,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}

	// server features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, false, false)
	f.AddResultHandler(e)
}

func (e *UCOPEV) AddUseCase() {
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
func (e *UCOPEV) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
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

	// check for required features
	evLoadControl, err := util.LoadControl(e.service, entity)
	if err != nil {
		return false, features.ErrFunctionNotSupported
	}

	// check if loadcontrol limit descriptions contains a recommendation category
	if _, err = evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation); err != nil {
		return false, err
	}

	return true, nil
}
