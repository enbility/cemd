package inverterbatteryvis

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type InverterBatteryVis struct {
	entity spineapi.EntityLocalInterface

	service api.ServiceInterface

	inverterEntity               spineapi.EntityRemoteInterface
	inverterElectricalConnection *features.ElectricalConnection
	inverterMeasurement          *features.Measurement

	ski string
}

var _ InverterBatteryVisInterface = (*InverterBatteryVis)(nil)

// Add InverterBatteryVis support
func NewInverterBatteryVis(service api.ServiceInterface, details *shipapi.ServiceDetails) *InverterBatteryVis {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	inverter := &InverterBatteryVis{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(inverter)

	service.RegisterRemoteSKI(ski, true)

	return inverter
}
