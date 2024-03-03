package ucevsecc

import "github.com/enbility/cemd/api"

const (
	// An EVSE was connected
	//
	// The callback with this message provides:
	//   - the device of the EVSE
	//   - the entity of the EVSE
	EvseConnected api.EventType = "EvseConnected"

	// An EVSE was disconnected
	//
	// The callback with this message provides:
	//   - the device of the EVSE
	//   - the entity of the EVSE
	EvseDisconnected api.EventType = "EvseDisconnected"

	// EVSE manufacturer data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE
	//   - the entity of the EVSE
	//
	// Use Case EVSECC, Scenario 1
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateManufacturerData api.EventType = "DataUpdateManufacturerData"

	// EVSE operation state was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE
	//   - the entity of the EVSE
	//
	// Use Case EVSECC, Scenario 2
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateOperatingState api.EventType = "DataUpdateOperatingState"
)
