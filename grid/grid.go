package grid

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type Grid struct {
	entity spineapi.EntityLocalInterface

	service api.ServiceInterface

	gridEntity spineapi.EntityRemoteInterface

	gridDeviceConfiguration  *features.DeviceConfiguration
	gridElectricalConnection *features.ElectricalConnection
	gridMeasurement          *features.Measurement

	ski string
}

var _ GridInterface = (*Grid)(nil)

// Add Grid support
func NewGrid(service api.ServiceInterface, details *shipapi.ServiceDetails) *Grid {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	grid := &Grid{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(grid)

	service.RegisterRemoteSKI(ski, true)

	return grid
}
