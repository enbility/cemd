package ucevcem

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCEvCEMInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the last current measurement for each phase of the connected EV
	EVCurrentsPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 2

	// return the last power measurement for each phase of the connected EV
	EVPowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error)

	// Scenario 3

	// return the charged energy measurement in Wh of the connected EV
	EVChargedEnergy(entity spineapi.EntityRemoteInterface) (float64, error)
}

const (
	// EV measurement data updated
	UCEvCEMMeasurementDataUpdate api.UseCaseEventType = "ucEvCEMMeasurementDataUpdate"
)
