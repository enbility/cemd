package ucmpc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the Monitoring of Power Consumption UseCase
type UCMCPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the momentary active power consumption or production
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the momentary active phase specific power consumption or production per phase
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 2

	// return the total consumption energy
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	//   - positive values are used for consumption
	EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error)

	// return the total feed in energy
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// return values:
	//   - negative values are used for production
	EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return the momentary phase specific current consumption or production
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	//
	// return values
	//   - positive values are used for consumption
	//   - negative values are used for production
	CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 4

	// return the phase specific voltage details
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 5

	// return frequency
	//
	// parameters:
	//   - entity: the entity of the device (e.g. EVSE)
	Frequency(entity spineapi.EntityRemoteInterface) (float64, error)
}
