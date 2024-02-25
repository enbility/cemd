package ucevcc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	serviceapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEVCC struct {
	service serviceapi.ServiceInterface

	reader api.EventReaderInterface

	validEntityTypes []model.EntityTypeType
}

var _ UCEVCCInterface = (*UCEVCC)(nil)

func NewUCEVCC(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.EventReaderInterface) *UCEVCC {
	uc := &UCEVCC{
		service: service,
		reader:  reader,
	}

	uc.validEntityTypes = []model.EntityTypeType{
		model.EntityTypeTypeEV,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEVCC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeEVCommissioningAndConfiguration
}

func (e *UCEVCC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeIdentification,
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeDeviceDiagnosis,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

func (e *UCEVCC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
}

// returns if the entity supports the usecase
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVCC) IsUseCaseSupported(entity spineapi.EntityRemoteInterface) (bool, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	// check if the usecase and mandatory scenarios are supported and
	// if the required server features are available
	if !entity.Device().VerifyUseCaseScenariosAndFeaturesSupport(
		model.UseCaseActorTypeEV,
		e.UseCaseName(),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 8},
		[]model.FeatureTypeType{model.FeatureTypeTypeDeviceConfiguration},
	) {
		return false, nil
	}

	return true, nil
}
