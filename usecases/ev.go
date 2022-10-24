package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Interface for the evCC use case for CEM device
type EVDelegate interface {
	// handle device state updates from the remote EV entity
	HandleEVEntityState(ski string, failure bool)
}

// EV Commissioning and Configuration Use Case implementation
type EV struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	Delegate EVDelegate
}

// Register the use case and features for handling EVs
// CEM will call this on startup
func NewEVCommissioningAndConfiguration(service *service.EEBUSService) *EV {
	// add the use case
	ev := &EV{
		service: service,
		entity:  service.LocalEntity(),
	}

	// subscribe to get incoming EV events
	spine.Events.Subscribe(ev)

	// add use cases
	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
	/*
		_ = spine.NewUseCase(
			ev.entity,
			model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3})

		_ = spine.NewUseCase(
			ev.entity,
			model.UseCaseNameTypeCoordinatedEVCharging,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
	*/

	return ev
}

// Invoke to remove an EV entity
// Called when an EV was disconnected
func (e *EV) UnregisterEV() {
	// remove the entity
	e.service.RemoveEntity(e.entity)
}
