package ucmpc

import "github.com/enbility/cemd/api"

const (
	// Total momentary active power consumption or production
	//
	// Use `Power` to get the current data
	//
	// Use Case MCP, Scenario 1
	DataUpdatePower api.EventType = "DataUpdatePower"

	// Phase specific momentary active power consumption or production
	//
	// Use `PowerPerPhase` to get the current data
	//
	// Use Case MCP, Scenario 1
	DataUpdatePowerPerPhase api.EventType = "DataUpdatePowerPerPhase"

	// Total energy consumed
	//
	// Use `EnergyConsumed` to get the current data
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyConsumed api.EventType = "DataUpdateEnergyConsumed"

	// Total energy produced
	//
	// Use `EnergyProduced` to get the current data
	//
	// Use Case MCP, Scenario 2
	DataUpdateEnergyProduced api.EventType = "DataUpdateEnergyProduced"

	// Phase specific momentary current consumption or production
	//
	// Use `CurrentPerPhase` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateCurrentsPerPhase api.EventType = "DataUpdateCurrentsPerPhase"

	// Phase specific voltage
	//
	// Use `VoltagePerPhase` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateVoltagePerPhase api.EventType = "DataUpdateVoltagePerPhase"

	// Power network frequency data updated
	//
	// Use `Frequency` to get the current data
	//
	// Use Case MCP, Scenario 3
	DataUpdateFrequency api.EventType = "DataUpdateFrequency"
)
