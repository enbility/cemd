package ucvapd

import "github.com/enbility/cemd/api"

const (
	// PV System total power data updated
	//
	// Use `Power` to get the current data
	//
	// Use Case VAPD, Scenario 1
	DataUpdatePower api.EventType = "ucvapd-DataUpdatePower"

	// PV System nominal peak power data updated
	//
	// Use `PowerNominalPeak` to get the current data
	//
	// Use Case VAPD, Scenario 2
	DataUpdatePowerNominalPeak api.EventType = "ucvapd-DataUpdatePowerNominalPeak"

	// PV System total yield data updated
	//
	// Use `PVYieldTotal` to get the current data
	//
	// Use Case VAPD, Scenario 3
	DataUpdatePVYieldTotal api.EventType = "ucvapd-DataUpdatePVYieldTotal"
)
