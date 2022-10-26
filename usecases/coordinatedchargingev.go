package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Coordinated EV Charging Use Case implementation
type CoordinatedChargingEV struct {
	*spine.UseCaseImpl
	service *service.EEBUSService
}

// Register the use case and features for coordinated ev charging
// CEM will call this on startup
func NewCoordinatedChargingEV(service *service.EEBUSService) *CoordinatedChargingEV {
	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &CoordinatedChargingEV{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeCoordinatedEVCharging,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	return useCase
}
