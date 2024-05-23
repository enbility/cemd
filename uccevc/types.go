package uccevc

import "github.com/enbility/cemd/api"

const (
	// Scenario 1

	// EV provided an energy demand
	//
	// Use `EnergyDemand` to get the current data
	DataUpdateEnergyDemand api.EventType = "uccevc-DataUpdateEnergyDemand"

	// Scenario 2

	// EV provided a charge plan constraints
	//
	// Use `TimeSlotConstraints` to get the current data
	DataUpdateTimeSlotConstraints api.EventType = "uccevc-DataUpdateTimeSlotConstraints"

	// Scenario 3

	// EV incentive table data updated
	//
	// Use `IncentiveConstraints` to get the current data
	DataUpdateIncentiveTable api.EventType = "uccevc-DataUpdateIncentiveTable"

	// EV requested an incentive table, call to WriteIncentiveTableDescriptions required
	DataRequestedIncentiveTableDescription api.EventType = "uccevc-DataRequestedIncentiveTableDescription"

	// Scenario 2 & 3

	// EV requested power limits, call to WritePowerLimits and WriteIncentives required
	DataRequestedPowerLimitsAndIncentives api.EventType = "uccevc-DataRequestedPowerLimitsAndIncentives"

	// Scenario 4

	// EV provided a charge plan
	//
	// Use `ChargePlanConstraints` to get the current data
	DataUpdateChargePlanConstraints api.EventType = "uccevc-DataUpdateChargePlanConstraints"

	// EV provided a charge plan
	//
	// Use `ChargePlan` to get the current data
	DataUpdateChargePlan api.EventType = "uccevc-DataUpdateChargePlan"
)
