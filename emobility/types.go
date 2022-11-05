package emobility

type EVCommunicationStandardType string

const (
	EVCommunicationStandardTypeUnknown      EVCommunicationStandardType = "unknown"
	EVCommunicationStandardTypeISO151182ED1 EVCommunicationStandardType = "iso15118-2ed1"
	EVCommunicationStandardTypeISO151182ED2 EVCommunicationStandardType = "iso15118-2ed2"
	EVCommunicationStandardTypeIEC61851     EVCommunicationStandardType = "iec61851"
)

type EVChargeStateType string

const (
	EVChargeStateTypeUnkown    EVChargeStateType = "Unknown"
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
