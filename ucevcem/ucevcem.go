package ucevcem

import (
	"github.com/enbility/cemd/api"
	serviceapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type UCEvCEM struct {
	service serviceapi.ServiceInterface

	reader api.UseCaseEventReaderInterface
}

var _ UCEvCEMInterface = (*UCEvCEM)(nil)

func NewUCEvCEM(service serviceapi.ServiceInterface, details *shipapi.ServiceDetails, reader api.UseCaseEventReaderInterface) *UCEvCEM {
	uc := &UCEvCEM{
		service: service,
		reader:  reader,
	}

	_ = spine.Events.Subscribe(uc)

	return uc
}

func (c *UCEvCEM) UseCaseName() model.UseCaseNameType {
	return model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging
}

func (e *UCEvCEM) AddFeatures() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	f := localEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	f.AddResultHandler(e)
}

func (e *UCEvCEM) AddUseCase() {
	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeCEM,
		e.UseCaseName(),
		model.SpecificationVersionType("1.0.1"),
		"",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}
