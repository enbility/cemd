package ucevsecc

import "github.com/enbility/cemd/api"

const (
	// An EVSE was connected
	EvseConnected api.EventType = "EvseConnected"

	// An EVSE was disconnected
	EvseDisconnected api.EventType = "EvseDisconnected"

	// EVSE manufacturer data was updated
	//
	// Use `ManufacturerData` to get the current data
	//
	// Use Case EVSECC, Scenario 1
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateManufacturerData api.EventType = "DataUpdateManufacturerData"

	// EVSE operation state was updated
	//
	// Use `OperatingState` to get the current data
	//
	// Use Case EVSECC, Scenario 2
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateOperatingState api.EventType = "DataUpdateOperatingState"
)
