package grid

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
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
	entity spineapi.EntityLocal

	service api.EEBUSService

	gridEntity spineapi.EntityRemote

	gridDeviceConfiguration  *features.DeviceConfiguration
	gridElectricalConnection *features.ElectricalConnection
	gridMeasurement          *features.Measurement

	ski string
}

var _ GridI = (*GridImpl)(nil)

// Add Grid support
func NewGrid(service api.EEBUSService, details *api.ServiceDetails) *GridImpl {
	ski := util.NormalizeSKI(details.SKI)

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	grid := &GridImpl{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(grid)

	service.RegisterRemoteSKI(ski, true)

	return grid
}
