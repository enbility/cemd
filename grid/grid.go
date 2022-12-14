package grid

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/util"
)

type GridI interface {
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
	ski := util.NormalizeSKI(details.SKI())

	grid := &GridImpl{
		service: service,
		entity:  service.LocalEntity(),
		ski:     ski,
	}
	spine.Events.Subscribe(grid)

	service.PairRemoteService(details)

	return grid
}
