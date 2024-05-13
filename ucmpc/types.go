package ucmpc

import "github.com/enbility/cemd/api"

const (
	// Total momentary active power consumption or production
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 1
	DataUpdatePower api.EventType = "DataUpdatePower"

	// Phase specific momentary active power consumption or production
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 1
	DataUpdatePowerPerPhase api.EventType = "DataUpdatePowerPerPhase"

	// Total energy consumed
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyConsumed api.EventType = "DataUpdateEnergyConsumed"

	// Total energy produced
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyProduced api.EventType = "DataUpdateEnergyProduced"

	// Phase specific momentary current consumption or production
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 3
	DataUpdateCurrentsPerPhase api.EventType = "DataUpdateCurrentsPerPhase"

	// Phase specific voltage
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 3
	DataUpdateVoltagePerPhase api.EventType = "DataUpdateVoltagePerPhase"

	// Power network frequency data updated
	//
	// The callback with this message provides:
	//   - the device of the e.g. EVSE
	//   - the entity of the e.g. EVSE
	//
	// Use Case MCP, Scenario 3
	DataUpdateFrequency api.EventType = "DataUpdateFrequency"
)
