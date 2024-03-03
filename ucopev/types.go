package ucopev

import "github.com/enbility/cemd/api"

const (
	// EV load control obligation limit data updated
	//
	// The callback with this message provides:
	//   - the device of the EVSE the EV is connected to
	//   - the entity of the EV
	DataUpdateLimit api.EventType = "DataUpdateLimit"
)
