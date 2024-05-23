package uclppserver

import "github.com/enbility/cemd/api"

const (
	// Load control obligation limit data update received
	//
	// Use `ProductionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "uclppserver-DataUpdateLimit"

	// An incoming load control obligation limit needs to be approved or denied
	//
	// Use `PendingProductionLimits` to get the currently pending write approval requests
	// and invoke `ApproveOrDenyProductionLimit` for each
	//
	// Use Case LPC, Scenario 1
	WriteApprovalRequired api.EventType = "uclppserver-WriteApprovalRequired"

	// Failsafe limit for the produced active (real) power of the
	// Controllable System data update received
	//
	// Use `FailsafeProductionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeProductionActivePowerLimit api.EventType = "uclppserver-DataUpdateFailsafeProductionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data update received
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "uclppserver-DataUpdateFailsafeDurationMinimum"

	// Indicates a notify heartbeat event the application should care of.
	// E.g. going into or out of the Failsafe state
	//
	// Use Case LPP, Scenario 3
	DataUpdateHeartbeat api.EventType = "uclpcserver-DataUpdateHeartbeat"
)
