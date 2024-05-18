package ucevsoc

import "github.com/enbility/cemd/api"

const (
	// EV state of charge data was updated
	//
	// Use `StateOfCharge` to get the current data
	//
	// Use Case EVSOC, Scenario 1
	DataUpdateStateOfCharge api.EventType = "DataUpdateStateOfCharge"
)
