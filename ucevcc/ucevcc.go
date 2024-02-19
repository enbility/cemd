package ucevcc

import (
	"github.com/enbility/cemd/api"
	serviceapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEvCC struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCEvCCInterface = (*UCEvCC)(nil)

func NewUCEvCC(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCEvCC {
	uc := &UCEvCC{
		service: service,
		reader:  reader,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEvCC) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeEVCommissioningAndConfiguration
}

func (e *UCEvCC) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeIdentification,
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
		model.FeatureTypeTypeLoadControl,
	}

	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

func (e *UCEvCC) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
}
