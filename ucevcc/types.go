package ucevcc

import "github.com/enbility/cemd/api"

const (

	// An EV was connected
	//
	// Use Case EVCC, Scenario 1
	EvConnected api.EventType = "EvConnected"

	// An EV was disconnected
	//
	// Use Case EVCC, Scenario 8
	EvDisconnected api.EventType = "EvDisconnected"

	// EV charge state data was updated
	DataUpdateChargeState api.EventType = "DataUpdateChargeState"

	// EV communication standard data was updated
	//
	// Use Case EVCC, Scenario 2
	// Note: the referred data may be updated together with all other configuration items of this use case
	DataUpdateCommunicationStandard api.EventType = "DataUpdateCommunicationStandard"

	// EV asymmetric charging data was updated
	//
	// Use Case EVCC, Scenario 3
	//
	// Note: the referred data may be updated together with all other configuration items of this use case
	AsymmetricChargingSupportDataUpdate api.EventType = "AsymmetricChargingSupportDataUpdate"

	// EV identificationdata was updated
	//
	// Use Case EVCC, Scenario 4
	DataUpdateIdentifications api.EventType = "DataUpdateIdentifications"

	// EV manufacturer data was updated
	//
	// Use Case EVCC, Scenario 5
	DataUpdateManufacturerData api.EventType = "DataUpdateManufacturerData"

	// EV charging power limits
	//
	// Use Case EVCC, Scenario 6
	DataUpdateCurrentLimits api.EventType = "DataUpdateCurrentLimits"

	// EV permitted power limits updated
	//
	// Use Case EVCC, Scenario 7
	DataUpdateIsInSleepMode api.EventType = "DataUpdateIsInSleepMode"
)
