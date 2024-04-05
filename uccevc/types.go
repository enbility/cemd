package uccevc

import "github.com/enbility/cemd/api"

const (
	// Scenario 1

	// EV provided an energy demand
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateEnergyDemand api.EventType = "DataUpdateEnergyDemand"

	// Scenario 2

	// EV provided a charge plan constraints
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateTimeSlotConstraints api.EventType = "DataUpdateTimeSlotConstraints"

	// Scenario 3

	// EV incentive table data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateIncentiveTable api.EventType = "DataUpdateIncentiveTable"

	// EV incentive table data constraints updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateIncentiveTableConstraints api.EventType = "DataUpdateIncentiveTableConstraints"

	// EV requested an incentive table, call to WriteIncentiveTableDescriptions required
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataRequestedIncentiveTableDescription api.EventType = "DataRequestedIncentiveTableDescription"

	// Scenario 2 & 3

	// EV requested power limits, call to WritePowerLimits and WriteIncentives required
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataRequestedPowerLimitsAndIncentives api.EventType = "DataRequestedPowerLimitsAndIncentives"

	// Scenario 4

	// EV provided a charge plan
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateChargePlanConstraints api.EventType = "DataUpdateChargePlanConstraints"

	// EV provided a charge plan
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateChargePlan api.EventType = "DataUpdateChargePlan"
)
