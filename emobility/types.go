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
	EVChargeStrategyTypeTimedCharging  EVChargeStrategyType = "timedcharging"
)

// Contains details about power limits or incentives for a defined timeframe
type EVDurationSlotValue struct {
	Duration time.Duration // Duration of this slot
	Value    float64       // Energy Cost or Power Limit
}

var ErrEVDisconnected = errors.New("ev is disconnected")
var ErrNotSupported = errors.New("function is not supported")
