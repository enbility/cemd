package ucmgcp

import "github.com/enbility/cemd/api"

const (
	// Grid maximum allowed feed-in power as percentage value of the cumulated
	// nominal peak power of all electricity producting PV systems was updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 2
	DataUpdatePowerLimitationFactor api.EventType = "DataUpdatePowerLimitationFactor"

	// Grid momentary power consumption/production data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 2
	DataUpdatePower api.EventType = "DataUpdatePower"

	// Total grid feed in energy data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 3
	DataUpdateEnergyFeedIn api.EventType = "DataUpdateEnergyFeedIn"

	// Total grid consumed energy data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 4
	DataUpdateEnergyConsumed api.EventType = "DataUpdateEnergyConsumed"

	// Phase specific momentary current consumption/production phase detail data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 5
	DataUpdateCurrentPerPhase api.EventType = "DataUpdateCurrentPerPhase"

	// Phase specific voltage at the grid connection point
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 6
	DataUpdateVoltagePerPhase api.EventType = "DataUpdateVoltagePerPhase"

	// Grid frequency data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. SMGW
	//   - the entity of the e.g. SMGW
	//
	// Use Case MGCP, Scenario 7
	DataUpdateFrequency api.EventType = "DataUpdateFrequency"
)
