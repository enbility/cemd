package ucvapd

import "github.com/enbility/cemd/api"

const (
	// PV System total power data updated
	//
	// The callback with this message provides:
	//   - the device of the inverter
	//   - the entity of the inverter
	//
	// Use Case VAPD, Scenario 1
	DataUpdatePower api.EventType = "DataUpdatePower"

	// PV System nominal peak power data updated
	//
	// The callback with this message provides:
	//   - the device of the inverter
	//   - the entity of the inverter
	//
	// Use Case VAPD, Scenario 2
	DataUpdatePowerNominalPeak api.EventType = "DataUpdatePowerNominalPeak"

	// PV System total yield data updated
	//
	// The callback with this message provides:
	//   - the device of the inverter
	//   - the entity of the inverter
	//
	// Use Case VAPD, Scenario 3
	DataUpdatePVYieldTotal api.EventType = "DataUpdatePVYieldTotal"
)
