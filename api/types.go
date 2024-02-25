package api

import (
	"errors"
	"time"

	"github.com/enbility/spine-go/model"
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

// Defines a phase specific limit
type LoadLimitsPhase struct {
	Phase    model.ElectricalConnectionPhaseNameType
	IsActive bool
	Value    float64
}

// identification
type IdentificationItem struct {
	// the identification value
	Value string

	// the type of the identification value, e.g.
	ValueType model.IdentificationTypeType
}

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
type Demand struct {
	MinDemand          float64 // minimum demand in Wh to reach the minSoC setting, 0 if not set
	OptDemand          float64 // demand in Wh to reach the timer SoC setting
	MaxDemand          float64 // the maximum possible demand until the battery is full
	DurationUntilStart float64 // the duration in s from now until charging will start, this could be in the future but usualy is now
	DurationUntilEnd   float64 // the duration in s from now until minDemand or optDemand has to be reached, 0 if direct charge strategy is active
}

// Contains details about an EV generated charging plan
type ChargePlan struct {
	Slots []ChargePlanSlotValue // Individual charging slot details
}

// Contains details about a charging plan slot
type ChargePlanSlotValue struct {
	Start    time.Time // The start time of the slot
	End      time.Time // The duration of the slot
	Value    float64   // planned power value
	MinValue float64   // minimum power value
	MaxValue float64   // maximum power value
}

// Details about the time slot constraints
type TimeSlotConstraints struct {
	MinSlots             uint          // the minimum number of slots, no minimum if 0
	MaxSlots             uint          // the maximum number of slots, unlimited if 0
	MinSlotDuration      time.Duration // the minimum duration of a slot, no minimum if 0
	MaxSlotDuration      time.Duration // the maximum duration of a slot, unlimited if 0
	SlotDurationStepSize time.Duration // the duration has to be a multiple of this value if != 0
}

// Details about the incentive slot constraints
type IncentiveSlotConstraints struct {
	MinSlots uint // the minimum number of slots, no minimum if 0
	MaxSlots uint // the maximum number of slots, unlimited if 0
}

// details about the boundary
type TierBoundaryDescription struct {
	// the id of the boundary
	Id uint

	// the type of the boundary
	Type model.TierBoundaryTypeType

	// the unit of the boundary
	Unit model.UnitOfMeasurementType
}

// details about incentive
type IncentiveDescription struct {
	// the id of the incentive
	Id uint

	// the type of the incentive
	Type model.IncentiveTypeType

	// the currency of the incentive, if it is price based
	Currency model.CurrencyType
}

// Contains about one tier in a tariff
type IncentiveTableDescriptionTier struct {
	// the id of the tier
	Id uint

	// the tiers type
	Type model.TierTypeType

	// each tear has 1 to 3 boundaries
	// used for different power limits, e.g. 0-1kW x€, 1-3kW y€, ...
	Boundaries []TierBoundaryDescription

	// each tier has 1 to 3 incentives
	//   - price/costs (absolute or relative)
	//   - renewable energy percentage
	//   - CO2 emissions
	Incentives []IncentiveDescription
}

// Contains details about a tariff
type IncentiveTariffDescription struct {
	// each tariff can have 1 to 3 tiers
	Tiers []IncentiveTableDescriptionTier
}

// Contains details about power limits or incentives for a defined timeframe
type DurationSlotValue struct {
	Duration time.Duration // Duration of this slot
	Value    float64       // Energy Cost or Power Limit
}

// value if the UCEVCC communication standard is unknown
const (
	UCEVCCCommunicationStandardUnknown string = "unknown"
)

// type for cem and usecase specfic event names
type EventType string

const (
	// CEM

	// A paired remote device was connected
	DeviceConnected EventType = "deviceConnected"

	// A paired remote device was disconnected
	DeviceDisconnected EventType = "deviceDisconnected"

	// Visible remote eebus services list was updated
	VisibleRemoteServicesUpdated EventType = "visibleRemoteServicesUpdated"

	// UCCEVC

	// EV provided an energy demand
	UCCEVCEnergyDemandProvided EventType = "ucCEVCEnergyDemandProvided"

	// EV provided a charge plan
	UCCEVCChargePlanProvided EventType = "ucCEVCChargePlanProvided"

	// EV provided a charge plan constraints
	UCCEVCChargePlanConstraintsProvided EventType = "ucCEVCChargePlanConstraintsProvided"

	UCCEVCIncentiveDescriptionsRequired EventType = "ucCEVCIncentiveDescriptionsRequired"

	// EV incentive table data updated
	UCCEVCIncentiveTableDataUpdate EventType = "ucCEVCIncentiveTableDataUpdate"

	// EV requested power limits
	UCCEVPowerLimitsRequested EventType = "ucCEVPowerLimitsRequested"

	// EV requested incentives
	UCCEVCIncentivesRequested EventType = "ucCEVCIncentivesRequested"

	// UCEVCC

	// An EV was connected
	//
	// Use Case EVCC, Scenario 1
	UCEVCCEventConnected EventType = "ucEVCCEventConnected"

	// An EV was disconnected
	//
	// Use Case EVCC, Scenario 8
	UCEVCCEventDisconnected EventType = "ucEVCCEventDisconnected"

	// EV communication standard data was updated
	//
	// Use Case EVCC, Scenario 2
	// Note: the referred data may be updated together with all other configuration items of this use case
	UCEVCCCommunicationStandardConfigurationDataUpdate EventType = "ucEVCCCommunicationStandardConfigurationDataUpdate"

	// EV asymmetric charging data was updated
	//
	// Use Case EVCC, Scenario 3
	//
	// Note: the referred data may be updated together with all other configuration items of this use case
	UCEVCCAsymmetricChargingConfigurationDataUpdate EventType = "ucEVCCAsymmetricChargingConfigurationDataUpdate"

	// EV identificationdata was updated
	//
	// Use Case EVCC, Scenario 4
	UCEVCCIdentificationDataUpdate EventType = "ucEVCCIdentificationDataUpdate"

	// EV manufacturer data was updated
	//
	// Use Case EVCC, Scenario 5
	UCEVCCManufacturerDataUpdate EventType = "ucEVCCManufacturerDataUpdate"

	// EV charging power limits
	//
	// Use Case EVCC, Scenario 6
	UCEVCCChargingPowerLimitsDataUpdate EventType = "ucEVCCChargingPowerLimitsDataUpdate"

	// EV permitted power limits updated
	//
	// Use Case EVCC, Scenario 7
	UCEVCCSleepModeDataUpdate EventType = "ucEVCCSleepModeDataUpdate"

	// UCEVCEM

	// EV number of connected phases data updated
	//
	// Use Case EVCEM, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVCEMNumberOfConnectedPhasesDataUpdate EventType = "ucEVCEMNumberOfConnectedPhasesDataUpdate"

	// EV current measurement data updated
	//
	// Use Case EVCEM, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVCEMCurrentMeasurementDataUpdate EventType = "ucEVCEMCurrentMeasurementDataUpdate"

	// EV power measurement data updated
	//
	// Use Case EVCEM, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVCEMPowerMeasurementDataUpdate EventType = "ucEVCEMCurrentMeasurementDataUpdate"

	// EV charging energy measurement data updated
	//
	// Use Case EVCEM, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVCEMChargingEnergyMeasurementDataUpdate EventType = "UCEVCEMChargingEnergyMeasurementDataUpdate"

	// UCEVSECC

	// An EVSE was connected
	UCEVSECCEventConnected EventType = "ucEVSEConnected"

	// An EVSE was disconnected
	UCEVSECCEventDisconnected EventType = "ucEVSEDisonnected"

	// EVSE manufacturer data was updated
	//
	// Use Case EVSECC, Scenario 1
	UCEVSECCEventManufacturerUpdate EventType = "ucEVSEManufacturerUpdate"

	// EVSE operation state was updated
	//
	// Use Case EVSECC, Scenario 2
	UCEVSECCEventOperationStateUpdate EventType = "ucEVSEOperationStateUpdate"

	// UCEVSOC

	// EV state of charge data was updated
	//
	// Use Case EVSOC, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVSOCStateOfChargeMeasurementDataUpdate EventType = "ucEVSOCStateOfChargeMeasurementDataUpdate"

	// EV nominal capacity data was updated
	//
	// Use Case EVSOC, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVSOCNominalCapacityMeasurementDataUpdate EventType = "ucEVSOCNominalCapacityMeasurementDataUpdate"

	// EV state of health data was updated
	//
	// Use Case EVSOC, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	EVSOCStateOfHealthMeasurementDataUpdate EventType = "ucEVSOCStateOfHealthMeasurementDataUpdate"

	// EV actual range data was updated
	//
	// Use Case EVSOC, Scenario 4
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCEVSOCActualRangeMeasurementDataUpdate EventType = "ucEVSOCActualRangeMeasurementDataUpdate"

	// MGCP

	// Grid maximum allowed feed-in power as percentage value of the cumulated
	// nominal peak power of all electricity producting PV systems was updated
	//
	// Use Case MGCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPPVFeedInPowerLimitationFactorDataUpdate EventType = "ucMGCPPVFeedInPowerLimitationFactorDataUpdate"

	// Grid momentary power consumption/production data updated
	//
	// Use Case MGCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPPowerTotalMeasurementDataUpdate EventType = "ucMGCPPowerTotalMeasurementDataUpdate"

	// MTotal grid feed in energy data updated
	//
	// Use Case MGCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPGridFeedInMeasurementDataUpdate EventType = "ucMGCPGridFeedInMeasurementDataUpdate"

	// Total grid consumed energy data updated
	//
	// Use Case MGCP, Scenario 4
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPGridConsumptionMeasurementDataUpdate EventType = "ucMGCPGridConsumptionMeasurementDataUpdate"

	// Phase specific momentary current consumption/production phase detail data updated
	//
	// Use Case MGCP, Scenario 5
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPCurrentsMeasurementDataUpdate EventType = "ucMGCPCurrentsMeasurementDataUpdate"

	// Phase specific voltage at the grid connection point
	//
	// Use Case MGCP, Scenario 6
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPVoltagesMeasurementDataUpdate EventType = "ucMGCPVoltagesMeasurementDataUpdate"

	// Grid frequency data updated
	//
	// Use Case MGCP, Scenario 7
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMGCPFrequencyMeasurementDataUpdate EventType = "ucMGCPFrequencyMeasurementDataUpdate"

	// UCMPC

	// Total momentary active power consumption or production
	//
	// Use Case MCP, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCPowerTotalMeasurementDataUpdate EventType = "ucMPCPowerTotalMeasurementDataUpdate"

	// Phase specific momentary active power consumption or production
	//
	// Use Case MCP, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCPowerPerPhaseMeasurementDataUpdate EventType = "ucMPCPowerPerPhaseMeasurementDataUpdate"

	// Total energy consumed
	//
	// Use Case MCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCEnergyConsumedMeasurementDataUpdate EventType = "ucMPCEnergyConsumedMeasurementDataUpdate"

	// Total energy produced
	//
	// Use Case MCP, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCEnergyProcudedMeasurementDataUpdate EventType = "ucMPCEnergyProcudedMeasurementDataUpdate"

	// Phase specific momentary current consumption or production
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCCurrentsMeasurementDataUpdate EventType = "ucMPCCurrentsMeasurementDataUpdate"

	// Phase specific voltage
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCVoltagesMeasurementDataUpdate EventType = "ucMPCVoltagesMeasurementDataUpdate"

	// Power network frequency data updated
	//
	// Use Case MCP, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCMPCFrequencyMeasurementDataUpdate EventType = "ucMPCFrequencyMeasurementDataUpdate"

	// UCOPEV

	// EV load control obligation limit data updated
	UCOPEVLoadControlLimitDataUpdate EventType = "ucOPEVLoadControlLimitDataUpdate"

	// UCOSCEV

	// EV load control recommendation limit data updated
	//
	// Use Case OSCEV, Scenario 1
	UCOSCEVLoadControlLimitDataUpdate EventType = "ucOSCEVLoadControlLimitDataUpdate"

	// UCVABD

	// Battery System (dis)charge power data updated
	//
	// Use Case VABD, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVABDPowerTotalMeasurementDataUpdate EventType = "ucVABDPowerTotalMeasurementDataUpdate"

	// Battery System cumulated charge energy data updated
	//
	// Use Case VABD, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVABDChargeMeasurementDataUpdate EventType = "ucVABDChargeMeasurementDataUpdate"

	// Battery System cumulated discharge energy data updated
	//
	// Use Case VABD, Scenario 2
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVABDDischargeMeasurementDataUpdate EventType = "ucVABDDischargeMeasurementDataUpdate"

	// Battery System state of charge data updated
	//
	// Use Case VABD, Scenario 4
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVABDStateOfChargeMeasurementDataUpdate EventType = "ucVABDStateOfChargeMeasurementDataUpdate"

	// UCVAPD

	// PV System total power data updated
	//
	// Use Case VAPD, Scenario 1
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVAPDPowerTotalMeasurementDataUpdate EventType = "ucVAPDPowerTotalMeasurementDataUpdate"

	// PV System nominal peak power data updated
	//
	// Use Case VAPD, Scenario 2
	UCVAPDPeakPowerDataUpdate EventType = "ucVAPDPeakPowerDataUpdate"

	// PV System total yield data updated
	//
	// Use Case VAPD, Scenario 3
	//
	// Note: the referred data may be updated together with all other measurement items of this use case
	UCVAPDYieldTotalMeasurementDataUpdate EventType = "ucVAPDYieldTotalMeasurementDataUpdate"
)

var ErrNoCompatibleEntity = errors.New("entity is not an compatible entity")
