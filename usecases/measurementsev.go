package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// MeasurementOfElectricityDuringEVCharging use case
type MeasurementOfElectricityDuringEVCharging struct {
	*spine.UseCaseImpl
	service *service.EEBUSService
}

// Register the features for electrical connection
// CEM will call this on startup
func NewMeasurementOfElectricityDuringEVCharging(service *service.EEBUSService) *MeasurementOfElectricityDuringEVCharging {
	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &MeasurementOfElectricityDuringEVCharging{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	// add the features
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
	}

	return useCase
}
