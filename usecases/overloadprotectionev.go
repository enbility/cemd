package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Overload Protection by EV Charging Current Curtailment Use Case implementation
type OverloadProtectionEV struct {
	*spine.UseCaseImpl
	service *service.EEBUSService
}

// Register the use case and features for overload protection
// CEM will call this on startup
func NewOverloadProtectionEV(service *service.EEBUSService) *OverloadProtectionEV {
	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &OverloadProtectionEV{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
			model.SpecificationVersionType("1.0.1b"),
			[]model.UseCaseScenarioSupportType{1, 2, 3}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	return useCase
}
