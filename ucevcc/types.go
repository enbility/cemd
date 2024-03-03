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
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 1
	EvConnected api.EventType = "EvConnected"

	// An EV was disconnected
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 8
	EvDisconnected api.EventType = "EvDisconnected"

	// EV charge state data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateChargeState api.EventType = "DataUpdateChargeState"

	// EV communication standard data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 2
	// Note: the referred data may be updated together with all other configuration items of this use case
	DataUpdateCommunicationStandard api.EventType = "DataUpdateCommunicationStandard"

	// EV asymmetric charging data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Note: the referred data may be updated together with all other configuration items of this use case
	AsymmetricChargingSupportDataUpdate api.EventType = "AsymmetricChargingSupportDataUpdate"

	// EV identificationdata was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 4
	DataUpdateIdentifications api.EventType = "DataUpdateIdentifications"

	// EV manufacturer data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 5
	DataUpdateManufacturerData api.EventType = "DataUpdateManufacturerData"

	// EV charging power limits
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 6
	DataUpdateCurrentLimits api.EventType = "DataUpdateCurrentLimits"

	// EV permitted power limits updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVCC, Scenario 7
	DataUpdateIsInSleepMode api.EventType = "DataUpdateIsInSleepMode"
)
