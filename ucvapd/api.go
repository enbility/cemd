package ucvapd

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the Visualization of Aggregated Photovoltaic Data UseCase
type UCVAPDInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current production power
	//
	// parameters:
	//   - entity: the entity of the inverter
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the nominal peak power
	//
	// parameters:
	//   - entity: the entity of the inverter
	PowerNominalPeak(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return total PV yield
	//
	// parameters:
	//   - entity: the entity of the inverter
	PVYieldTotal(entity spineapi.EntityRemoteInterface) (float64, error)
}
