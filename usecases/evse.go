package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EVSECommissioningAndConfiguration struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService
}

// Add EVSE support
func NewEVSECommissioningAndConfiguration(service *service.EEBUSService) *EVSECommissioningAndConfiguration {
	// add the use case
	evse := &EVSECommissioningAndConfiguration{
		service: service,
		entity:  service.LocalEntity(),
	}
	spine.Events.Subscribe(evse)

	_ = spine.NewUseCase(
		evse.entity,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2})

	return evse
}
