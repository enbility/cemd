package ucmpc

import "github.com/enbility/cemd/api"

const (
	// Total momentary active power consumption or production
	//
	// Use Case MCP, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdatePower api.EventType = "DataUpdatePower"

	// Phase specific momentary active power consumption or production
	//
	// Use Case MCP, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdatePowerPerPhase api.EventType = "DataUpdatePowerPerPhase"

	// Total energy consumed
	//
	// Use Case MCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateEnergyConsumed api.EventType = "DataUpdateEnergyConsumed"

	// Total energy produced
	//
	// Use Case MCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateEnergyProduced api.EventType = "DataUpdateEnergyProduced"

	// Phase specific momentary current consumption or production
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateCurrentsPerPhase api.EventType = "DataUpdateCurrentsPerPhase"

	// Phase specific voltage
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateVoltagePerPhase api.EventType = "DataUpdateVoltagePerPhase"

	// Power network frequency data updated
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateFrequency api.EventType = "DataUpdateFrequency"
)
