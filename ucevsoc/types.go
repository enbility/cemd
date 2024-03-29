package ucevsoc

import "github.com/enbility/cemd/api"

const (
	// EV state of charge data was updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case EVSOC, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateStateOfCharge api.EventType = "DataUpdateStateOfCharge"
)
