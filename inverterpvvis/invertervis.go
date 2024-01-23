package inverterpvvis

import (
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type InverterPVVis struct {
	entity spineapi.EntityLocalInterface

	service eebusapi.ServiceInterface

	inverterEntity               spineapi.EntityRemoteInterface
	inverterDeviceConfiguration  *features.DeviceConfiguration
	inverterElectricalConnection *features.ElectricalConnection
	inverterMeasurement          *features.Measurement

	ski string
}

var _ InverterPVVisInterface = (*InverterPVVis)(nil)

// Add InverterPVVis support
func NewInverterPVVis(service eebusapi.ServiceInterface, details *shipapi.ServiceDetails) *InverterPVVis {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	inverter := &InverterPVVis{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(inverter)

	service.RegisterRemoteSKI(ski, true)

	return inverter
}
