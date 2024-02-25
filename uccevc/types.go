package uccevc

import "github.com/enbility/cemd/api"

const (
	// Scenario 1

	// EV provided an energy demand
	DataUpdateEnergyDemand api.EventType = "DataUpdateEnergyDemand"

	// Scenario 2

	// EV provided a charge plan constraints
	DataUpdateTimeSlotConstraints api.EventType = "DataUpdateTimeSlotConstraints"

	// Scenario 3

	// EV incentive table data updated
	DataUpdateIncentiveTable api.EventType = "DataUpdateIncentiveTable"

	// EV requested an incentive table, call to WriteIncentiveTableDescriptions required
	DataRequestedIncentiveTableDescription api.EventType = "DataRequestedIncentiveTableDescription"

	// Scenario 2 & 3

	// EV requested power limits, call to WritePowerLimits and WriteIncentives required
	DataRequestedPowerLimitsAndIncentives api.EventType = "DataRequestedPowerLimitsAndIncentives"

	// Scenario 4

	// EV provided a charge plan
	DataUpdateChargePlanConstraints api.EventType = "DataUpdateChargePlanConstraints"

	// EV provided a charge plan
	DataUpdateChargePlan api.EventType = "DataUpdateChargePlan"
)
