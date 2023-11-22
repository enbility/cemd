package grid

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/util"
)

type GridI interface {
	PowerLimitationFactor() (float64, error)
	MomentaryPowerConsumptionOrProduction() (float64, error)
	TotalFeedInEnergy() (float64, error)
	TotalConsumedEnergy() (float64, error)
	MomentaryCurrentConsumptionOrProduction() ([]float64, error)
	Voltage() ([]float64, error)
	Frequency() (float64, error)
}

type GridImpl struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	gridEntity *spine.EntityRemoteImpl

	gridDeviceConfiguration  *features.DeviceConfiguration
	gridElectricalConnection *features.ElectricalConnection
	gridMeasurement          *features.Measurement

	ski string
}

var _ GridI = (*GridImpl)(nil)

// Add Grid support
func NewGrid(service *service.EEBUSService, details *service.ServiceDetails) *GridImpl {
	ski := util.NormalizeSKI(details.SKI)

	grid := &GridImpl{
		service: service,
		entity:  service.LocalEntity(),
		ski:     ski,
	}
	_ = spine.Events.Subscribe(grid)

	service.RegisterRemoteSKI(ski, true)

	return grid
}
