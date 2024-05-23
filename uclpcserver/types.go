package uclpcserver

import "github.com/enbility/cemd/api"

const (
	// Load control obligation limit data update received
	//
	// Use `ConsumptionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "uclpcserver-DataUpdateLimit"

	// An incoming load control obligation limit needs to be approved or denied
	//
	// Use `PendingConsumptionLimits` to get the currently pending write approval requests
	// and invoke `ApproveOrDenyConsumptionLimit` for each
	//
	// Use Case LPC, Scenario 1
	WriteApprovalRequired api.EventType = "uclpcserver-WriteApprovalRequired"

	// Failsafe limit for the consumed active (real) power of the
	// Controllable System data update received
	//
	// Use `FailsafeConsumptionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeConsumptionActivePowerLimit api.EventType = "uclpcserver-DataUpdateFailsafeConsumptionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data update received
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "uclpcserver-DataUpdateFailsafeDurationMinimum"
)
