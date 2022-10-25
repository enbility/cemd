package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Optimization of Self Consumption During EV Charging Use Case implementation
type OptimizationOfSelfConsumptionEV struct {
	*spine.UseCaseImpl
	service *service.EEBUSService
}

// Register the use case and features for optimization of self consumption
// CEM will call this on startup
func NewOptimizationOfSelfConsumptionEV(service *service.EEBUSService) *OptimizationOfSelfConsumptionEV {
	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &OptimizationOfSelfConsumptionEV{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
			model.SpecificationVersionType("1.0.1b"),
			[]model.UseCaseScenarioSupportType{1, 2, 3}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	return useCase
}
