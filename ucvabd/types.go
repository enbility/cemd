package ucvabd

import "github.com/enbility/cemd/api"

const (
	// Battery System (dis)charge power data updated
	//
	// Use `Power` to get the current data
	//
	// Use Case VABD, Scenario 1
	DataUpdatePower api.EventType = "d"

	// Battery System cumulated charge energy data updated
	//
	// Use `EnergyCharged` to get the current data
	//
	// Use Case VABD, Scenario 2
	DataUpdateEnergyCharged api.EventType = "DataUpdateEnergyCharged"

	// Battery System cumulated discharge energy data updated
	//
	// Use `EnergyDischarged` to get the current data
	//
	// Use Case VABD, Scenario 2
	DataUpdateEnergyDischarged api.EventType = "DataUpdateEnergyDischarged"

	// Battery System state of charge data updated
	//
	// Use `StateOfCharge` to get the current data
	//
	// Use Case VABD, Scenario 4
	DataUpdateStateOfCharge api.EventType = "DataUpdateStateOfCharge"
)
