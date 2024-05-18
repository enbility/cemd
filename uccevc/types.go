package uccevc

import "github.com/enbility/cemd/api"

const (
	// Scenario 1

	// EV provided an energy demand
	//
	// Use `EnergyDemand` to get the current data
	DataUpdateEnergyDemand api.EventType = "DataUpdateEnergyDemand"

	// Scenario 2

	// EV provided a charge plan constraints
	//
	// Use `TimeSlotConstraints` to get the current data
	DataUpdateTimeSlotConstraints api.EventType = "DataUpdateTimeSlotConstraints"

	// Scenario 3

	// EV incentive table data updated
	//
	// Use `IncentiveConstraints` to get the current data
	DataUpdateIncentiveTable api.EventType = "DataUpdateIncentiveTable"

	// EV requested an incentive table, call to WriteIncentiveTableDescriptions required
	DataRequestedIncentiveTableDescription api.EventType = "DataRequestedIncentiveTableDescription"

	// Scenario 2 & 3

	// EV requested power limits, call to WritePowerLimits and WriteIncentives required
	DataRequestedPowerLimitsAndIncentives api.EventType = "DataRequestedPowerLimitsAndIncentives"

	// Scenario 4

	// EV provided a charge plan
	//
	// Use `ChargePlanConstraints` to get the current data
	DataUpdateChargePlanConstraints api.EventType = "DataUpdateChargePlanConstraints"

	// EV provided a charge plan
	//
	// Use `ChargePlan` to get the current data
	DataUpdateChargePlan api.EventType = "DataUpdateChargePlan"
)
