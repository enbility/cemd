package inverterpvvis

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type InverterPVVisI interface {
	CurrentProductionPower() (float64, error)
	NominalPeakPower() (float64, error)
	TotalPVYield() (float64, error)
}

type InverterPVVisImpl struct {
	entity spineapi.EntityLocal

	service api.EEBUSService

	inverterEntity               spineapi.EntityRemote
	inverterDeviceConfiguration  *features.DeviceConfiguration
	inverterElectricalConnection *features.ElectricalConnection
	inverterMeasurement          *features.Measurement

	ski string
}

var _ InverterPVVisI = (*InverterPVVisImpl)(nil)

// Add InverterPVVis support
func NewInverterPVVis(service api.EEBUSService, details *shipapi.ServiceDetails) *InverterPVVisImpl {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	inverter := &InverterPVVisImpl{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(inverter)

	service.RegisterRemoteSKI(ski, true)

	return inverter
}
