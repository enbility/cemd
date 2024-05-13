package uclppserver

import "github.com/enbility/cemd/api"

const (
	// Load control obligation limit data update received
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "DataUpdateLimit"

	// Failsafe limit for the produced active (real) power of the
	// Controllable System data update received
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeProductionActivePowerLimit api.EventType = "DataUpdateFailsafeProductionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data update received
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "DataUpdateFailsafeDurationMinimum"
)
