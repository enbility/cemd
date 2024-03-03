package ucoscev

import "github.com/enbility/cemd/api"

const (
	// EV load control recommendation limit data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	//
	// Use Case OSCEV, Scenario 1
	DataUpdateLimit api.EventType = "DataUpdateLimit"
)
