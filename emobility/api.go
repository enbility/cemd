package emobility

import "github.com/enbility/spine-go/api"

//go:generate mockgen -source emobility.go -destination mock_emobility_test.go -package emobility

// used by emobility and implemented by the CEM
type EmobilityDataProvider interface {
	// The EV provided a charge strategy
	EVProvidedChargeStrategy(strategy EVChargeStrategyType)

	// EV provided an energy demand
	//
	// Parameters:
	//   - demand: Contains details about the actual demands from the EV
	EVProvidedEnergyDemand(demand EVDemand)

	// Energy demand and duration is provided by the EV which requires the CEM
	// to respond with time slots containing power limits for each slot
	//
	// `EVWritePowerLimits` must be invoked within <55s, idealy <15s, after receiving this call
	//
	// Parameters:
	//   - demand: Contains details about the actual demands from the EV
	//   - constraints: Contains details about the time slot constraints
	EVRequestPowerLimits(demand EVDemand, constraints EVTimeSlotConstraints)

	// Energy demand and duration is provided by the EV which requires the CEM
	// to respond with time slots containing incentives for each slot
	//
	// `EVWriteIncentives` must be invoked within <20s after receiving this call
	//
	// Parameters:
	//   - demand: Contains details about the actual demands from the EV
	//   - constraints: Contains details about the incentive slot constraints
	EVRequestIncentives(demand EVDemand, constraints EVIncentiveSlotConstraints)

	// The EV provided a charge plan
	EVProvidedChargePlan(plan EVChargePlan)

	// The EV provided charge plan constraints
	EVProvidedChargePlanConstraints(constraints []EVDurationSlotValue)
}

// used by the CEM and implemented by emobility
type EMobilityInterface interface {
	// return if an EV is connected
	EVConnected(remoteEntity api.EntityRemoteInterface) bool

	// return the current charge state of the EV
	EVCurrentChargeState(remoteEntity api.EntityRemoteInterface) (EVChargeStateType, error)

	// return the current loadcontrol obligation limits
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVLoadControlObligationLimits(remoteEntity api.EntityRemoteInterface) ([]float64, error)

	// send new LoadControlLimits to the remote EV
	//
	// parameters:
	//   - obligations: Overload Protection Limits per phase in A
	//   - recommendations: Self Consumption recommendations per phase in A
	//
	// obligations:
	// Sets a maximum A limit for each phase that the EV may not exceed.
	// Mainly used for implementing overload protection of the site or limiting the
	// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
	// and with ISO15118 if the EV does not support the Optimization of Self Consumption
	// usecase.
	//
	// recommendations:
	// Sets a recommended charge power in A for each phase. This is mainly
	// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
	// The EV either needs to support the Optimization of Self Consumption usecase or
	// the EVSE needs to be able map the recommendations into oligation limits which then
	// works for all EVs communication either via IEC61851 or ISO15118.
	//
	// note:
	// For obligations to work for optimizing solar excess power, the EV needs to
	// have an energy demand. Recommendations work even if the EV does not have an active
	// energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
	// In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific
	// and needs to have specific EVSE support for the specific EV brand.
	// In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
	EVWriteLoadControlLimits(remoteEntity api.EntityRemoteInterface, limits []EVLoadLimits) error

	// returns if the EVSE and EV combination support optimzation of self consumption
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVOptimizationOfSelfConsumptionSupported(remoteEntity api.EntityRemoteInterface) (bool, error)

	// return if the EVSE and EV combination support providing an SoC
	//
	// requires EVSoCSupported to return true
	// only works with a current ISO15118-2 with VAS or ISO15118-20
	// communication between EVSE and EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVSoCSupported(remoteEntity api.EntityRemoteInterface) (bool, error)

	// return the last known SoC of the connected EV
	//
	// requires EVSoCSupported to return true
	// only works with a current ISO15118-2 with VAS or ISO15118-20
	// communication between EVSE and EV
	//
	// possible errors:
	//   - ErrNotSupported if support for SoC is not possible
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVSoC(remoteEntity api.EntityRemoteInterface) (float64, error)

	// returns if the EVSE and EV combination support coordinated charging
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVCoordinatedChargingSupported(remoteEntity api.EntityRemoteInterface) (bool, error)

	// returns the current charging stratey
	//
	// returns EVChargeStrategyTypeUnknown if it could not be determined, e.g.
	// if the vehicle communication is via IEC61851 or the EV doesn't provide
	// any information about its charging mode or plan
	EVChargeStrategy(remoteEntity api.EntityRemoteInterface) EVChargeStrategyType

	// returns the current energy demand
	//   - EVDemand: details about the actual demands from the EV
	//   - error: if no data is available
	//
	// if duration is 0, direct charging is active, otherwise timed charging is active
	EVEnergyDemand(remoteEntity api.EntityRemoteInterface) (EVDemand, error)

	// returns the current charge plan
	//   - EVChargePlan: details about the actual charge plan provided by the EV
	//   - error: if no data is available
	EVChargePlan(remoteEntity api.EntityRemoteInterface) (EVChargePlan, error)

	// returns the constraints for the time slots
	//   - EVTimeSlotConstraints: details about the time slot constraints
	//   - error: if no data is available
	EVTimeSlotConstraints(remoteEntity api.EntityRemoteInterface) (EVTimeSlotConstraints, error)

	// send power limits data to the EV
	//
	// returns an error if sending failed or charge slot count do not meet requirements
	//
	// this needs to be invoked either <55s, idealy <15s, of receiving a call to EVRequestPowerLimits
	// or if the CEM requires the EV to change its charge plan
	EVWritePowerLimits(remoteEntity api.EntityRemoteInterface, data []EVDurationSlotValue) error

	// returns the constraints for incentive slots
	//   - EVIncentiveConstraints: details about the incentive slot constraints
	//   - error: if no data is available
	EVIncentiveConstraints(remoteEntity api.EntityRemoteInterface) (EVIncentiveSlotConstraints, error)

	// send price slots data to the EV
	//
	// returns an error if sending failed or charge slot count do not meet requirements
	//
	// this needs to be invoked either within 20s of receiving a call to EVRequestIncentives
	// or if the CEM requires the EV to change its charge plan
	EVWriteIncentives(remoteEntity api.EntityRemoteInterface, data []EVDurationSlotValue) error
}
