package invertervis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/util"
)

type InverterVisI interface {
	CurrentDisChargePower() (float64, error)
	TotalChargeEnergy() (float64, error)
	TotalDischargeEnergy() (float64, error)
	CurrentStateOfCharge() (float64, error)
}

type InverterVisImpl struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	inverterEntity *spine.EntityRemoteImpl

	inverterElectricalConnection *features.ElectricalConnection
	inverterMeasurement          *features.Measurement

	ski string
}

var _ InverterVisI = (*InverterVisImpl)(nil)

// Add InverterVis support
func NewInverterVis(service *service.EEBUSService, details *service.ServiceDetails) *InverterVisImpl {
	ski := util.NormalizeSKI(details.SKI())

	inverter := &InverterVisImpl{
		service: service,
		entity:  service.LocalEntity(),
		ski:     ski,
	}
	spine.Events.Subscribe(inverter)

	service.PairRemoteService(details)

	return inverter
}
