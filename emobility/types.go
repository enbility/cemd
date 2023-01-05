package emobility

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/spine/model"
)

type EVCommunicationStandardType model.DeviceConfigurationKeyValueStringType

const (
	EVCommunicationStandardTypeUnknown      EVCommunicationStandardType = "unknown"
	EVCommunicationStandardTypeISO151182ED1 EVCommunicationStandardType = "iso15118-2ed1"
	EVCommunicationStandardTypeISO151182ED2 EVCommunicationStandardType = "iso15118-2ed2"
	EVCommunicationStandardTypeIEC61851     EVCommunicationStandardType = "iec61851"
)

type EVChargeStateType string

const (
	EVChargeStateTypeUnknown   EVChargeStateType = "Unknown"
	EVChargeStateTypeUnplugged EVChargeStateType = "unplugged"
	EVChargeStateTypeError     EVChargeStateType = "error"
	EVChargeStateTypePaused    EVChargeStateType = "paused"
	EVChargeStateTypeActive    EVChargeStateType = "active"
	EVChargeStateTypeFinished  EVChargeStateType = "finished"
)

type EVChargeStrategyType string

const (
	EVChargeStrategyTypeUnknown        EVChargeStrategyType = "unknown"
	EVChargeStrategyTypeNoDemand       EVChargeStrategyType = "nodemand"
	EVChargeStrategyTypeDirectCharging EVChargeStrategyType = "directcharging"
	EVChargeStrategyTypeMinSoC         EVChargeStrategyType = "minsoc"
	EVChargeStrategyTypeTimedCharging  EVChargeStrategyType = "timedcharging"
)

// Contains details about the actual demands from the EV
//
// General:
//   - If duration and energy is 0, charge mode is EVChargeStrategyTypeNoDemand
//   - If duration is 0, charge mode is EVChargeStrategyTypeDirectCharging and the slots should cover at least 48h
//   - If both are != 0, charge mode is EVChargeStrategyTypeTimedCharging and the slots should cover at least the duration, but at max 168h (7d)
type EVDemand struct {
	MinDemand          float64       // minimum demand in Wh to reach the minSoC setting, 0 if not set
	OptDemand          float64       // demand in Wh to reach the timer SoC setting
	MaxDemand          float64       // the maximum possible demand until the battery is full
	DurationUntilStart time.Duration // the duration from now until charging will start, this could be in the future but usualy is now
	DurationUntilEnd   time.Duration // the duration from now until minDemand or optDemand has to be reached, 0 if direct charge strategy is active
}

// Details about the time slot constraints
type EVTimeSlotConstraints struct {
	MinSlots             uint          // the minimum number of slots, no minimum if 0
	MaxSlots             uint          // the maximum number of slots, unlimited if 0
	MinSlotDuration      time.Duration // the minimum duration of a slot, no minimum if 0
	MaxSlotDuration      time.Duration // the maximum duration of a slot, unlimited if 0
	SlotDurationStepSize time.Duration // the duration has to be a multiple of this value if != 0
}

// Details about the incentive slot constraints
type EVIncentiveSlotConstraints struct {
	MinSlots uint // the minimum number of slots, no minimum if 0
	MaxSlots uint // the maximum number of slots, unlimited if 0
}

// Contains details about power limits or incentives for a defined timeframe
type EVDurationSlotValue struct {
	Duration time.Duration // Duration of this slot
	Value    float64       // Energy Cost or Power Limit
}

var ErrEVDisconnected = errors.New("ev is disconnected")
var ErrNotSupported = errors.New("function is not supported")

// Allows to exclude some features
type EmobilityConfiguration struct {
	CoordinatedChargingDisabled bool
}
