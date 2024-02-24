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
	CurrentProductionPower(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the nominal peak power
	NominalPeakPower(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return total PV yield
	TotalPVYield(entity spineapi.EntityRemoteInterface) (float64, error)
}
