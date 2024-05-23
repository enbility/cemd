package ucopev

import "github.com/enbility/cemd/api"

const (
	// EV current limits
	//
	// Use `CurrentLimits` to get the current data
	DataUpdateCurrentLimits api.EventType = "ucopev-DataUpdateCurrentLimits"

	// EV load control obligation limit data updated
	//
	// Use `LoadControlLimits` to get the current data
	DataUpdateLimit api.EventType = "ucopev-DataUpdateLimit"
)
