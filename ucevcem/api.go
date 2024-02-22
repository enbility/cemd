package ucevcem

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCEVCEMInterface interface {
	api.UseCaseInterface

	// return the number of ac connected phases of the EV or 0 if it is unknown
	ConnectedPhases(entity spineapi.EntityRemoteInterface) (uint, error)

	// Scenario 1

	// return the last current measurement for each phase of the connected EV
	CurrentsPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 2

	// return the last power measurement for each phase of the connected EV
	PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 3

	// return the charged energy measurement in Wh of the connected EV
	ChargedEnergy(entity spineapi.EntityRemoteInterface) (float64, error)
}
