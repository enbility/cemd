package ucevsoc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the EV State Of Charge UseCase
type UCEVSOCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the EVscurrent state of charge of the EV or an error it is unknown
	StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2 to 4 are not supported, as there is no EV supporting this as of today
}
