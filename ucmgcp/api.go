package ucmgcp

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the Monitoring of Grid Connection Point UseCase
type UCMGCPInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current power limitation factor
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	PowerLimitationFactor(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the momentary power consumption or production at the grid connection point
	//
	//   - positive values are used for consumption
	//   - negative values are used for production
	Power(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return the total feed in energy at the grid connection point
	//
	//   - negative values are used for production
	EnergyFeedIn(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 4

	// return the total consumption energy at the grid connection point
	//
	//   - positive values are used for consumption
	EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 5

	// return the momentary current consumption or production at the grid connection point
	//
	//   - positive values are used for consumption
	//   - negative values are used for production
	CurrentsPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 6

	// return the voltage phase details at the grid connection point
	//
	VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 7

	// return frequency at the grid connection point
	//
	Frequency(entity spineapi.EntityRemoteInterface) (float64, error)
}
