package ucevcem

import "github.com/enbility/cemd/api"

const (
	// EV number of connected phases data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCEM, Scenario 1
	DataUpdatePhasesConnected api.EventType = "DataUpdatePhasesConnected"

	// EV current measurement data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCEM, Scenario 1
	DataUpdateCurrentPerPhase api.EventType = "DataUpdateCurrentPerPhase"

	// EV power measurement data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCEM, Scenario 2
	DataUpdatePowerPerPhase api.EventType = "DataUpdatePowerPerPhase"

	// EV charging energy measurement data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCEM, Scenario 3
	DataUpdateEnergyCharged api.EventType = "DataUpdateEnergyCharged"
)
