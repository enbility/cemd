package ucvapd

import "github.com/enbility/cemd/api"

const (
	// PV System total power data updated
	//
	// Use Case VAPD, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdatePower api.EventType = "DataUpdatePower"

	// PV System nominal peak power data updated
	//
	// Use Case VAPD, Scenario 2
	DataUpdatePowerNominalPeak api.EventType = "DataUpdatePowerNominalPeak"

	// PV System total yield data updated
	//
	// Use Case VAPD, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	DataUpdatePVYieldTotal api.EventType = "DataUpdatePVYieldTotal"
)
