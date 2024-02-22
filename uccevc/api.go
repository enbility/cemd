package uccevc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCCEVCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// returns the current charging stratey
	//
	// returns EVChargeStrategyTypeUnknown if it could not be determined, e.g.
	// if the vehicle communication is via IEC61851 or the EV doesn't provide
	// any information about its charging mode or plan
	ChargeStrategy(remoteEntity spineapi.EntityRemoteInterface) api.EVChargeStrategyType

	// returns the current energy demand
	//   - EVDemand: details about the actual demands from the EV
	//   - error: if no data is available
	//
	// if duration is 0, direct charging is active, otherwise timed charging is active
	EnergyDemand(remoteEntity spineapi.EntityRemoteInterface) (api.Demand, error)

	// Scenario 2

	TimeSlotConstraints(entity spineapi.EntityRemoteInterface) (api.TimeSlotConstraints, error)

	// send power limits to the EV
	// if no data is provided, default power limits with the max possible value for 7 days will be sent
	WritePowerLimits(entity spineapi.EntityRemoteInterface, data []api.DurationSlotValue) error

	// Scenario 3

	IncentiveConstraints(entity spineapi.EntityRemoteInterface) (api.IncentiveSlotConstraints, error)

	WriteIncentiveTableDescriptions(entity spineapi.EntityRemoteInterface, data []api.IncentiveTariffDescription) error

	// send incentives to the EV
	// if no data is provided, default incentives with the same price for 7 days will be sent
	WriteIncentives(entity spineapi.EntityRemoteInterface, data []api.DurationSlotValue) error

	// Scenario 4

	ChargePlanConstraints(entity spineapi.EntityRemoteInterface) ([]api.DurationSlotValue, error)

	ChargePlan(entity spineapi.EntityRemoteInterface) (api.ChargePlan, error)

	// Scenario 5 & 6

	// this is automatically covered by the SPINE implementation
}
