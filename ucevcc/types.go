package ucevcc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/spine-go/model"
)

// value if the UCEVCC communication standard is unknown
const (
	UCEVCCCommunicationStandardUnknown model.DeviceConfigurationKeyValueStringType = "unknown"
)

const (

	// An EV was connected
	//
	// Use Case EVCC, Scenario 1
	EvConnected api.EventType = "ucevcc-EvConnected"

	// An EV was disconnected
	//
	// Note: The ev entity is no longer connected to the device!
	//
	// Use Case EVCC, Scenario 8
	EvDisconnected api.EventType = "ucevcc-EvDisconnected"

	// EV charge state data was updated
	//
	// Use `ChargeState` to get the current data
	DataUpdateChargeState api.EventType = "ucevcc-DataUpdateChargeState"

	// EV communication standard data was updated
	//
	// Use `CommunicationStandard` to get the current data
	//
	// Use Case EVCC, Scenario 2
	DataUpdateCommunicationStandard api.EventType = "ucevcc-DataUpdateCommunicationStandard"

	// EV asymmetric charging data was updated
	//
	// Use `AsymmetricChargingSupport` to get the current data
	DataUpdateAsymmetricChargingSupport api.EventType = "ucevcc-DataUpdateAsymmetricChargingSupport"

	// EV identificationdata was updated
	//
	// Use `Identifications` to get the current data
	//
	// Use Case EVCC, Scenario 4
	DataUpdateIdentifications api.EventType = "ucevcc-DataUpdateIdentifications"

	// EV manufacturer data was updated
	//
	// Use `ManufacturerData` to get the current data
	//
	// Use Case EVCC, Scenario 5
	DataUpdateManufacturerData api.EventType = "ucevcc-DataUpdateManufacturerData"

	// EV charging power limits
	//
	// Use `ChargingPowerLimits` to get the current data
	//
	// Use Case EVCC, Scenario 6
	DataUpdateCurrentLimits api.EventType = "ucevcc-DataUpdateCurrentLimits"

	// EV permitted power limits updated
	//
	// Use `IsInSleepMode` to get the current data
	//
	// Use Case EVCC, Scenario 7
	DataUpdateIsInSleepMode api.EventType = "ucevcc-DataUpdateIsInSleepMode"
)
