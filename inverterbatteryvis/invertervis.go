package inverterbatteryvis

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type InverterBatteryVisI interface {
	CurrentDisChargePower() (float64, error)
	TotalChargeEnergy() (float64, error)
	TotalDischargeEnergy() (float64, error)
	CurrentStateOfCharge() (float64, error)
}

type InverterBatteryVisImpl struct {
	entity spineapi.EntityLocal

	service api.EEBUSService

	inverterEntity               spineapi.EntityRemote
	inverterElectricalConnection *features.ElectricalConnection
	inverterMeasurement          *features.Measurement

	ski string
}

var _ InverterBatteryVisI = (*InverterBatteryVisImpl)(nil)

// Add InverterBatteryVis support
func NewInverterBatteryVis(service api.EEBUSService, details *api.ServiceDetails) *InverterBatteryVisImpl {
	ski := util.NormalizeSKI(details.SKI)

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	inverter := &InverterBatteryVisImpl{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}
	_ = spine.Events.Subscribe(inverter)

	service.RegisterRemoteSKI(ski, true)

	return inverter
}
