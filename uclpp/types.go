package uclpp

import "github.com/enbility/cemd/api"

const (
	// Load control obligation limit data updated
	//
	// Use `ProductionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "uclpp-DataUpdateLimit"

	// Failsafe limit for the produced active (real) power of the
	// Controllable System data updated
	//
	// Use `FailsafeProductionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeProductionActivePowerLimit api.EventType = "uclpp-DataUpdateFailsafeProductionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data updated
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "uclpp-DataUpdateFailsafeDurationMinimum"
)
