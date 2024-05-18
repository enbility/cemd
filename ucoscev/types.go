package ucoscev

import "github.com/enbility/cemd/api"

const (
	// EV current limits
	//
	// Use `CurrentLimits` to get the current data
	DataUpdateCurrentLimits api.EventType = "DataUpdateCurrentLimits"

	// EV load control recommendation limit data updated
	//
	// Use `LoadControlLimits` to get the current data
	//
	// Use Case OSCEV, Scenario 1
	DataUpdateLimit api.EventType = "DataUpdateLimit"
)
