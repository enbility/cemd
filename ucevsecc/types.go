package ucevsecc

import "github.com/enbility/cemd/api"

const (
	// An EVSE was connected
	EvseConnected api.EventType = "EvseConnected"

	// An EVSE was disconnected
	EvseDisconnected api.EventType = "EvseDisconnected"

	// EVSE manufacturer data was updated
	//
	// Use Case EVSECC, Scenario 1
	DataUpdateManufacturerData api.EventType = "DataUpdateManufacturerData"

	// EVSE operation state was updated
	//
	// Use Case EVSECC, Scenario 2
	DataUpdateOperatingState api.EventType = "DataUpdateOperatingState"
)
