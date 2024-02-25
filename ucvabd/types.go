package ucvabd

import "github.com/enbility/cemd/api"

const (
	// Battery System (dis)charge power data updated
	//
	// Use Case VABD, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdatePower api.EventType = "d"

	// Battery System cumulated charge energy data updated
	//
	// Use Case VABD, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateEnergyCharged api.EventType = "DataUpdateEnergyCharged"

	// Battery System cumulated discharge energy data updated
	//
	// Use Case VABD, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateEnergyDischarged api.EventType = "DataUpdateEnergyDischarged"

	// Battery System state of charge data updated
	//
	// Use Case VABD, Scenario 4
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdateStateOfCharge api.EventType = "DataUpdateStateOfCharge"
)
