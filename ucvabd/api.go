package ucvabd

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the Visualization of Aggregated Battery Data UseCase
type UCVABDInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current (dis)charging power
	CurrentChargePower(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2

	// return the cumulated battery system charge energy
	TotalChargeEnergy(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 3

	// return the cumulated battery system discharge energy
	TotalDischargeEnergy(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 4

	// return the current state of charge of the battery system
	CurrentStateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error)
}
