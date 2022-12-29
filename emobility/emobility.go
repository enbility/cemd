package emobility

import (
	"time"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

// used by emobility and implemented by the CEM
type EmobilityDataProvider interface {
	// Energy demand and duration is provided by the EV which requires the CEM
	// to respond with time slots containing power limits for each slot
	//
	// `EVWritePowerLimits` must be invoked within <55s, idealy <15s, after receiving this call
	//
	// Parameters:
	//   - minDemand: minimum demand in Wh to reach the minSoC setting, 0 if not set
	//   - optDemand: demand in Wh to reach the timer SoC setting
	//   - maxDemand: the maximum possible demand until the battery is full
	//   - durationUntilStart: duration until charging will start (usually 0)
	//   - durationUntilEnd: duration until charging has to end
	//   - minSlots: the minimum number of slots, no minimum if 0
	//   - maxSlots: the maximum number of slots, unlimited if 0
	//   - minSlotDuration: the minimum duration of a slot, no minimum if 0
	//   - maxSlotDuration: the maximum duration of a slot, unlimited if 0
	//   - slotDurationStepSize: the duration has to be a multiple of this value if != 0
	//
	// General:
	//  - If duration and energy is 0, charge mode is EVChargeStrategyTypeNoDemand
	//  - If duration is 0, charge mode is EVChargeStrategyTypeDirectCharging and the slots should cover at least 48h
	//  - If both are != 0, charge mode is EVChargeStrategyTypeTimedCharging and the slots should cover at least the duration, but at max 168h (7d)
	EVRequestPowerLimits(minDemand, optDemand, maxDemand float64, durationUntilStart, durationUntilEnd time.Duration, minSlots, maxSlots uint, minSlotDuration, maxSlotDuration, slotDurationStepSize time.Duration)

	// Energy demand and duration is provided by the EV which requires the CEM
	// to respond with time slots containing incentives for each slot
	//
	// `EVWriteIncentives` must be invoked within <20s after receiving this call
	//
	// Parameters:
	//   - duration: timeframe in which the energy demand is required
	//   - minSlots: minimum amount of slots required
	//   - maxSlots: maximum amount of slots allowed
	//
	// General:
	//  - If duration and energy is 0, charge mode is EVChargeStrategyTypeNoDemand
	//  - If duration is 0, charge mode is EVChargeStrategyTypeDirectCharging and the slots should cover at least 48h
	//  - If both are != 0, charge mode is EVChargeStrategyTypeTimedCharging and the slots should cover at least the duration, but at max 168h (7d)
	EVRequestIncentives(duration time.Duration, minSlots, maxSlots uint)

	// The EV provided a charge plan
	EVProvideChargePlan(data []EVDurationSlotValue)
}

// used by the CEM and implemented by emobility
type EmobilityI interface {
	// return the current charge sate of the EV
	EVCurrentChargeState() (EVChargeStateType, error)

	// return the number of ac connected phases of the EV or 0 if it is unknown
	EVConnectedPhases() (uint, error)

	// return the charged energy measurement in Wh of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVChargedEnergy() (float64, error)

	// return the last power measurement for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVPowerPerPhase() ([]float64, error)

	// return the last current measurement for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVCurrentsPerPhase() ([]float64, error)

	// return the min, max, default limits for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVCurrentLimits() ([]float64, []float64, []float64, error)

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
	EVWriteLoadControlLimits(obligations, recommendations []float64) error

	// return the current communication standard type used to communicate between EVSE and EV
	//
	// if an EV is connected via IEC61851, no ISO15118 specific data can be provided!
	// sometimes the connection starts with IEC61851 before it switches
	// to ISO15118, and sometimes it falls back again. so the error return is
	// never absolut for the whole connection time, except if the use case
	// is not supported
	//
	// the values are not constant and can change due to communication problems, bugs, and
	// sometimes communication starts with IEC61851 before it switches to ISO
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - ErrNotSupported if getting the communication standard is not supported
	//   - and others
	EVCommunicationStandard() (EVCommunicationStandardType, error)

	// returns the identification of the currently connected EV or nil if not available
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVIdentification() (string, error)

	// returns if the EVSE and EV combination support optimzation of self consumption
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVOptimizationOfSelfConsumptionSupported() (bool, error)

	// return if the EVSE and EV combination support providing an SoC
	//
	// requires EVSoCSupported to return true
	// only works with a current ISO15118-2 with VAS or ISO15118-20
	// communication between EVSE and EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVSoCSupported() (bool, error)

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
	EVSoC() (float64, error)

	// returns if the EVSE and EV combination support coordinated charging
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVCoordinatedChargingSupported() (bool, error)

	// returns the current charging stratey
	//
	// returns EVChargeStrategyTypeUnknown if it could not be determined, e.g.
	// if the vehicle communication is via IEC61851 or the EV doesn't provide
	// any information about its charging mode or plan
	EVChargeStrategy() EVChargeStrategyType

	// returns the current energy demand
	//   - minDemand: minimum demand in Wh to reach the minSoC setting, 0 if not set
	//   - optDemand: demand in Wh to reach the timer SoC setting
	//   - maxDemand: the maximum possible demand until the battery is full
	//   - durationUntilStart: the duration from now until charging will start, this could be in the future but usualy is now
	//   - durationUntilEnd: the duration from now until minDemand or optDemand has to be reached, 0 if direct charge strategy is active
	//   - error: if no data is available
	//
	// if duration is 0, direct charging is active, otherwise timed charging is active
	EVEnergyDemand() (float64, float64, float64, time.Duration, time.Duration, error)

	// returns the constraints for the power slots
	//   - minSlots: the minimum number of slots, no minimum if 0
	//   - maxSlots: the maximum number of slots, unlimited if 0
	//   - minSlotDuration: the minimum duration of a slot, no minimum if 0
	//   - maxSlotDuration: the maximum duration of a slot, unlimited if 0
	//   - slotDurationStepSize: the duration has to be a multiple of this value if != 0
	EVGetPowerConstraints() (uint, uint, time.Duration, time.Duration, time.Duration)

	// send power limits data to the EV
	//
	// returns an error if sending failed or charge slot count do not meet requirements
	//
	// this needs to be invoked either <55s, idealy <15s, of receiving a call to EVRequestPowerLimits
	// or if the CEM requires the EV to change its charge plan
	EVWritePowerLimits(data []EVDurationSlotValue) error

	// returns the constraints for incentive slots
	//   - minimum number of incentive slots, no minimum if 0
	//   - maximum number of incentive slots, unlimited if 0
	EVGetIncentiveConstraints() (uint, uint)

	// send price slots data to the EV
	//
	// returns an error if sending failed or charge slot count do not meet requirements
	//
	// this needs to be invoked either within 20s of receiving a call to EVRequestIncentives
	// or if the CEM requires the EV to change its charge plan
	EVWriteIncentives(data []EVDurationSlotValue) error
}

type EMobilityImpl struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	evseEntity *spine.EntityRemoteImpl
	evEntity   *spine.EntityRemoteImpl

	evseDeviceClassification *features.DeviceClassification
	evseDeviceDiagnosis      *features.DeviceDiagnosis

	evDeviceClassification *features.DeviceClassification
	evDeviceDiagnosis      *features.DeviceDiagnosis
	evDeviceConfiguration  *features.DeviceConfiguration
	evElectricalConnection *features.ElectricalConnection
	evMeasurement          *features.Measurement
	evIdentification       *features.Identification
	evLoadControl          *features.LoadControl
	evTimeSeries           *features.TimeSeries
	evIncentiveTable       *features.IncentiveTable

	ski      string
	currency model.CurrencyType

	dataProvider EmobilityDataProvider
}

var _ EmobilityI = (*EMobilityImpl)(nil)

// Add E-Mobility support
func NewEMobility(service *service.EEBUSService, details *service.ServiceDetails, currency model.CurrencyType, dataProvider EmobilityDataProvider) *EMobilityImpl {
	ski := util.NormalizeSKI(details.SKI())

	emobility := &EMobilityImpl{
		service:      service,
		entity:       service.LocalEntity(),
		ski:          ski,
		currency:     currency,
		dataProvider: dataProvider,
	}
	spine.Events.Subscribe(emobility)

	service.PairRemoteService(details)

	return emobility
}
