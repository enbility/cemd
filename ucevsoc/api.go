package ucevsoc

import (
	"github.com/enbility/cemd/api"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCEVSOCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// Scenario 2

	// this is automatically covered by the SPINE implementation

	// Scenario 3

	// this is covered by the central CEM interface implementation
	// use that one to set the CEM's operation state which will inform all remote devices

	// Scenario 4
}
