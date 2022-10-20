package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// EV Commissioning and Configuration Use Case implementation
type EVSoC struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService
}

// Register the use case and features for handling EV SoC
// CEM will call this on startup
func NewEVStateOfCharge(service *service.EEBUSService) *EVSoC {
	// add the use case
	evsoc := &EVSoC{
		service: service,
		entity:  service.LocalEntity(),
	}

	// subscribe to get incoming EV events
	spine.Events.Subscribe(evsoc)

	// add use cases
	_ = spine.NewUseCase(
		evsoc.entity,
		model.UseCaseNameTypeEVStateOfCharge,
		model.SpecificationVersionType("1.0.0"),
		[]model.UseCaseScenarioSupportType{1})

	return evsoc
}
