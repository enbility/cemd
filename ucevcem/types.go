package ucevcem

import "github.com/enbility/cemd/api"

const (
	// EV number of connected phases data updated
	//
	// Use `PhasesConnected` to get the current data
	//
	// Use Case EVCEM, Scenario 1
	DataUpdatePhasesConnected api.EventType = "ucevcem-DataUpdatePhasesConnected"

	// EV current measurement data updated
	//
	// Use `CurrentPerPhase` to get the current data
	//
	// Use Case EVCEM, Scenario 1
	DataUpdateCurrentPerPhase api.EventType = "ucevcem-DataUpdateCurrentPerPhase"

	// EV power measurement data updated
	//
	// Use `PowerPerPhase` to get the current data
	//
	// Use Case EVCEM, Scenario 2
	DataUpdatePowerPerPhase api.EventType = "ucevcem-DataUpdatePowerPerPhase"

	// EV charging energy measurement data updated
	//
	// Use `EnergyCharged` to get the current data
	//
	// Use Case EVCEM, Scenario 3
	DataUpdateEnergyCharged api.EventType = "ucevcem-DataUpdateEnergyCharged"
)
