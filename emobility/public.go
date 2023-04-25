package emobility

import (
	"errors"
	"time"

	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
	eebusUtil "github.com/enbility/eebus-go/util"
)

// return if an EV is connected
//
// this includes all required features and
// minimal data being available
func (e *EMobilityImpl) EVConnected() bool {
	// To report an EV as being connected, also consider all required
	// features to be available and assigned
	if e.evEntity == nil ||
		e.evDeviceDiagnosis == nil ||
		e.evElectricalConnection == nil ||
		e.evMeasurement == nil ||
		e.evLoadControl == nil ||
		e.evDeviceConfiguration == nil {
		return false
	}

	// getting current charge state should work
	if _, err := e.EVCurrentChargeState(); err != nil {
		return false
	}

	// the communication standard needs to be available
	if _, err := e.EVCommunicationStandard(); err != nil {
		return false
	}

	// getting currents measurements should work
	if _, err := e.EVCurrentsPerPhase(); err != nil {
		return false
	}

	// getting limits should work
	if _, err := e.EVLoadControlObligationLimits(); err != nil {
		return false
	}

	return true
}

// return the current charge state of the EV
func (e *EMobilityImpl) EVCurrentChargeState() (EVChargeStateType, error) {
	if e.evEntity == nil || e.evDeviceDiagnosis == nil {
		return EVChargeStateTypeUnplugged, nil
	}

	diagnosisState, err := e.evDeviceDiagnosis.GetState()
	if err != nil {
		return EVChargeStateTypeUnknown, err
	}

	operatingState := diagnosisState.OperatingState
	if operatingState == nil {
		return EVChargeStateTypeUnknown, features.ErrDataNotAvailable
	}

	switch *operatingState {
	case model.DeviceDiagnosisOperatingStateTypeNormalOperation:
		return EVChargeStateTypeActive, nil
	case model.DeviceDiagnosisOperatingStateTypeStandby:
		return EVChargeStateTypePaused, nil
	case model.DeviceDiagnosisOperatingStateTypeFailure:
		return EVChargeStateTypeError, nil
	case model.DeviceDiagnosisOperatingStateTypeFinished:
		return EVChargeStateTypeFinished, nil
	}

	return EVChargeStateTypeUnknown, nil
}

// return the number of ac connected phases of the EV or 0 if it is unknown
func (e *EMobilityImpl) EVConnectedPhases() (uint, error) {
	if e.evEntity == nil || e.evElectricalConnection == nil {
		return 0, ErrEVDisconnected
	}

	data, err := e.evElectricalConnection.GetDescriptions()
	if err != nil {
		return 0, features.ErrDataNotAvailable
	}

	for _, item := range data {
		if item.ElectricalConnectionId == nil {
			continue
		}

		if item.AcConnectedPhases != nil {
			return *item.AcConnectedPhases, nil
		}
	}

	// default to 3 if the value is not available
	return 3, nil
}

// return the charged energy measurement in Wh of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVChargedEnergy() (float64, error) {
	if e.evEntity == nil || e.evMeasurement == nil {
		return 0, ErrEVDisconnected
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeCharge
	data, err := e.evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume there is only one result
	value := data[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), err
}

// return the last power measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVPowerPerPhase() ([]float64, error) {
	if e.evEntity == nil || e.evMeasurement == nil {
		return nil, ErrEVDisconnected
	}

	var data []model.MeasurementDataType

	powerAvailable := true
	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPower
	data, err := e.evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil || len(data) == 0 {
		powerAvailable = false

		// If power is not provided, fall back to power calculations via currents
		measurement = model.MeasurementTypeTypeCurrent
		scope = model.ScopeTypeTypeACCurrent
		data, err = e.evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
		if err != nil {
			return nil, err
		}
	}

	var result []float64

	for _, phase := range util.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			elParam, err := e.evElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil || elParam.AcMeasuredPhases == nil || *elParam.AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			if !powerAvailable {
				phaseValue *= e.service.Configuration.Voltage()
			}

			result = append(result, phaseValue)
		}
	}

	return result, nil
}

// return the last current measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVCurrentsPerPhase() ([]float64, error) {
	if e.evEntity == nil || e.evElectricalConnection == nil {
		return nil, ErrEVDisconnected
	}

	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	data, err := e.evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	var result []float64

	for _, phase := range util.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			elParam, err := e.evElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil || elParam.AcMeasuredPhases == nil || *elParam.AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			result = append(result, phaseValue)
		}
	}

	return result, nil
}

// return the min, max, default limits for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVCurrentLimits() ([]float64, []float64, []float64, error) {
	if e.evEntity == nil || e.evElectricalConnection == nil {
		return nil, nil, nil, ErrEVDisconnected
	}

	var resultMin, resultMax, resultDefault []float64

	for _, phaseName := range util.PhaseNameMapping {
		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := e.evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.ParameterId == nil {
			continue
		}

		dataMin, dataMax, dataDefault, err := e.evElectricalConnection.GetLimitsForParameterId(*elParamDesc.ParameterId)
		if err != nil {
			continue
		}

		// Min current data should be derived from min power data
		// but as this value is only properly provided via VAS the
		// currrent min values can not be trusted.
		// Min current for 3-phase should be at least 2.2A (ISO)
		if dataMin < 2.2 {
			dataMin = 2.2
		}

		resultMin = append(resultMin, dataMin)
		resultMax = append(resultMax, dataMax)
		resultDefault = append(resultDefault, dataDefault)
	}

	if len(resultMin) == 0 {
		return nil, nil, nil, features.ErrDataNotAvailable
	}
	return resultMin, resultMax, resultDefault, nil
}

// return the current loadcontrol obligation limits
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVLoadControlObligationLimits() ([]float64, error) {
	if e.evEntity == nil || e.evElectricalConnection == nil || e.evLoadControl == nil {
		return nil, ErrEVDisconnected
	}

	// find out the appropriate limitId for each phase value
	// limitDescription contains the measurementId for each limitId
	limitDescriptions, err := e.evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation)
	if err != nil {
		return nil, features.ErrDataNotAvailable
	}

	var result []float64

	for i := 0; i < 3; i++ {
		phaseName := util.PhaseNameMapping[i]

		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := e.evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.MeasurementId == nil {
			// there is no data for this phase, the phase may not exit
			result = append(result, 0)
			continue
		}

		var limitDesc *model.LoadControlLimitDescriptionDataType
		for _, desc := range limitDescriptions {
			if desc.MeasurementId != nil && *desc.MeasurementId == *elParamDesc.MeasurementId {
				limitDesc = &desc
				break
			}
		}

		if limitDesc == nil || limitDesc.LimitId == nil {
			return nil, features.ErrDataNotAvailable
		}

		limitIdData, err := e.evLoadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
		if err != nil {
			return nil, features.ErrDataNotAvailable
		}

		if limitIdData.Value == nil {
			return nil, features.ErrDataNotAvailable
		}

		result = append(result, limitIdData.Value.GetValue())
	}

	return result, nil
}

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
// notes:
//   - For obligations to work for optimizing solar excess power, the EV needs to have an energy demand.
//   - Recommendations work even if the EV does not have an active energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
//   - In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific and needs to have specific EVSE support for the specific EV brand.
//   - In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
//   - Min power data is only provided via IEC61851 or using VAS in ISO15118-2.
func (e *EMobilityImpl) EVWriteLoadControlLimits(obligations, recommendations []float64) error {
	if e.evEntity == nil {
		return ErrEVDisconnected
	}

	if e.evElectricalConnection == nil || e.evLoadControl == nil {
		return features.ErrDataNotAvailable
	}

	var limitData []model.LoadControlLimitDataType

	for scopeTypes := 0; scopeTypes < 2; scopeTypes++ {
		category := model.LoadControlCategoryTypeObligation
		currentsPerPhase := obligations
		if scopeTypes == 1 {
			category = model.LoadControlCategoryTypeRecommendation
			currentsPerPhase = recommendations
		}

		for index, phaseLimit := range currentsPerPhase {
			phaseName := util.PhaseNameMapping[index]

			// find out the appropriate limitId for each phase value
			// limitDescription contains the measurementId for each limitId
			limitDescriptions, err := e.evLoadControl.GetLimitDescriptionsForCategory(category)
			if err != nil {
				continue
			}

			// electricalParameterDescription contains the measured phase for each measurementId
			elParamDesc, err := e.evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
			if err != nil || elParamDesc.MeasurementId == nil {
				continue
			}

			var limitDesc *model.LoadControlLimitDescriptionDataType
			for _, desc := range limitDescriptions {
				if desc.MeasurementId != nil && *desc.MeasurementId == *elParamDesc.MeasurementId {
					limitDesc = &desc
					break
				}
			}

			if limitDesc == nil || limitDesc.LimitId == nil {
				continue
			}

			limitIdData, err := e.evLoadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
			if err != nil {
				continue
			}

			// EEBus_UC_TS_OverloadProtectionByEvChargingCurrentCurtailment V1.01b 3.2.1.2.2.2
			// If omitted or set to "true", the timePeriod, value and isLimitActive element SHALL be writeable by a client.
			if limitIdData.IsLimitChangeable != nil && !*limitIdData.IsLimitChangeable {
				continue
			}

			// electricalPermittedValueSet contains the allowed min, max and the default values per phase
			phaseLimit = e.evElectricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(phaseLimit, *elParamDesc.ParameterId)

			newLimit := model.LoadControlLimitDataType{
				LimitId:       limitDesc.LimitId,
				IsLimitActive: eebusUtil.Ptr(true),
				Value:         model.NewScaledNumberType(phaseLimit),
			}
			limitData = append(limitData, newLimit)
		}
	}

	_, err := e.evLoadControl.WriteLimitValues(limitData)

	return err
}

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
func (e *EMobilityImpl) EVCommunicationStandard() (EVCommunicationStandardType, error) {
	if e.evEntity == nil || e.evDeviceConfiguration == nil {
		return EVCommunicationStandardTypeUnknown, ErrEVDisconnected
	}

	// check if device configuration descriptions has an communication standard key name
	_, err := e.evDeviceConfiguration.GetDescriptionForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}

	data, err := e.evDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard, model.DeviceConfigurationKeyValueTypeTypeString)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}

	if data == nil {
		return EVCommunicationStandardTypeUnknown, features.ErrDataNotAvailable
	}

	value := data.(*model.DeviceConfigurationKeyValueStringType)
	return EVCommunicationStandardType(*value), nil
}

// returns the identification of the currently connected EV or nil if not available
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EMobilityImpl) EVIdentification() (string, error) {
	if e.evEntity == nil {
		return "", ErrEVDisconnected
	}

	if e.evIdentification == nil {
		return "", features.ErrDataNotAvailable
	}

	identifications, err := e.evIdentification.GetValues()
	if err != nil {
		return "", err
	}

	for _, identification := range identifications {
		value := identification.IdentificationValue
		if value == nil {
			continue
		}

		return string(*value), nil
	}
	return "", nil
}

// returns if the EVSE and EV combination support optimzation of self consumption
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EMobilityImpl) EVOptimizationOfSelfConsumptionSupported() (bool, error) {
	if e.evEntity == nil || e.evLoadControl == nil {
		return false, ErrEVDisconnected
	}

	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the Optimization of self consumption usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	// check if loadcontrol limit descriptions contains a recommendation category
	if _, err = e.evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeRecommendation); err != nil {
		return false, err
	}

	return true, nil
}

// return if the EVSE and EV combination support providing an SoC
//
// requires EVSoCSupported to return true
// only works with a current ISO15118-2 with VAS or ISO15118-20
// communication between EVSE and EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVSoCSupported() (bool, error) {
	if e.evEntity == nil || e.evMeasurement == nil {
		return false, ErrEVDisconnected
	}

	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the SoC usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeEVStateOfCharge, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	// check if measurement descriptions has an SoC scope type
	desc, err := e.evMeasurement.GetDescriptionsForScope(model.ScopeTypeTypeStateOfCharge)
	if err != nil {
		return false, err
	}
	if len(desc) == 0 {
		return false, features.ErrDataNotAvailable
	}

	return true, nil
}

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
func (e *EMobilityImpl) EVSoC() (float64, error) {
	if e.evEntity == nil || e.evMeasurement == nil {
		return 0, ErrEVDisconnected
	}

	// check if the SoC is supported
	support, err := e.EVSoCSupported()
	if err != nil {
		return 0, err
	}
	if !support {
		return 0, features.ErrNotSupported
	}

	data, err := e.evMeasurement.GetValuesForTypeCommodityScope(model.MeasurementTypeTypePercentage, model.CommodityTypeTypeElectricity, model.ScopeTypeTypeStateOfCharge)
	if err != nil {
		return 0, err
	}

	// we assume there is only one value, nil is already checked
	value := data[0].Value

	return value.GetValue(), nil
}

// returns if the EVSE and EV combination support coordinated charging
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EMobilityImpl) EVCoordinatedChargingSupported() (bool, error) {
	if e.evEntity == nil {
		return false, ErrEVDisconnected
	}

	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the Coordinated charging usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeCoordinatedEVCharging, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	return true, nil
}

// returns the current charging strategy
func (e *EMobilityImpl) EVChargeStrategy() EVChargeStrategyType {
	if e.evEntity == nil || e.evTimeSeries == nil {
		return EVChargeStrategyTypeUnknown
	}

	// only ISO communication can provide a charging strategy information
	com, err := e.EVCommunicationStandard()
	if err != nil || com == EVCommunicationStandardTypeUnknown || com == EVCommunicationStandardTypeIEC61851 {
		return EVChargeStrategyTypeUnknown
	}

	// only the time series data for singledemand is relevant for detecting the charging strategy
	data, err := e.evTimeSeries.GetValueForType(model.TimeSeriesTypeTypeSingleDemand)
	if err != nil {
		return EVChargeStrategyTypeUnknown
	}

	// without time series slots, there is no known strategy
	if data.TimeSeriesSlot == nil || len(data.TimeSeriesSlot) == 0 {
		return EVChargeStrategyTypeUnknown
	}

	// get the value for the first slot
	firstSlot := data.TimeSeriesSlot[0]

	switch {
	case firstSlot.Duration == nil:
		// if value is > 0 and duration does not exist, the EV is direct charging
		if firstSlot.Value != nil {
			return EVChargeStrategyTypeDirectCharging
		}

	case firstSlot.Duration != nil:
		if _, err := firstSlot.Duration.GetTimeDuration(); err != nil {
			// we got an invalid duration
			return EVChargeStrategyTypeUnknown
		}

		if firstSlot.MinValue != nil && firstSlot.MinValue.GetValue() > 0 {
			return EVChargeStrategyTypeMinSoC
		}

		if firstSlot.Value != nil {
			if firstSlot.Value.GetValue() > 0 {
				// there is demand and a duration
				return EVChargeStrategyTypeTimedCharging
			}

			return EVChargeStrategyTypeNoDemand
		}

	}

	return EVChargeStrategyTypeUnknown
}

// returns the current energy demand in Wh and the duration
func (e *EMobilityImpl) EVEnergyDemand() (EVDemand, error) {
	demand := EVDemand{}

	if e.evEntity == nil {
		return demand, ErrEVDisconnected
	}

	if e.evTimeSeries == nil {
		return demand, features.ErrDataNotAvailable
	}

	data, err := e.evTimeSeries.GetValueForType(model.TimeSeriesTypeTypeSingleDemand)
	if err != nil {
		return demand, features.ErrDataNotAvailable
	}

	// we need at a time series slot
	if data.TimeSeriesSlot == nil {
		return demand, features.ErrDataNotAvailable
	}

	// get the value for the first slot, ignore all others, which
	// in the tests so far always have min/max/value 0
	firstSlot := data.TimeSeriesSlot[0]
	if firstSlot.MinValue != nil {
		demand.MinDemand = firstSlot.MinValue.GetValue()
	}
	if firstSlot.Value != nil {
		demand.OptDemand = firstSlot.Value.GetValue()
	}
	if firstSlot.MaxValue != nil {
		demand.MaxDemand = firstSlot.MaxValue.GetValue()
	}
	if firstSlot.Duration != nil {
		if tempDuration, err := firstSlot.Duration.GetTimeDuration(); err == nil {
			demand.DurationUntilEnd = tempDuration
		}
	}

	// start time has to be defined either in TimePeriod or the first slot
	relStartTime := time.Duration(0)

	startTimeSet := false
	if data.TimePeriod != nil && data.TimePeriod.StartTime != nil {
		if temp, err := data.TimePeriod.StartTime.GetTimeDuration(); err == nil {
			relStartTime = temp
			startTimeSet = true
		}
	}

	if !startTimeSet {
		if firstSlot.TimePeriod != nil && firstSlot.TimePeriod.StartTime != nil {
			if temp, err := firstSlot.TimePeriod.StartTime.GetTimeDuration(); err == nil {
				relStartTime = temp
			}
		}
	}

	demand.DurationUntilStart = relStartTime

	return demand, nil
}

// returns the constraints for the power slots
func (e *EMobilityImpl) EVGetPowerConstraints() EVTimeSlotConstraints {
	result := EVTimeSlotConstraints{}

	if e.evEntity == nil || e.evTimeSeries == nil {
		return result
	}

	constraints, err := e.evTimeSeries.GetConstraints()
	if err != nil {
		return result
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.SlotCountMin)
	}
	if constraint.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.SlotCountMax)
	}
	if constraint.SlotDurationMin != nil {
		if duration, err := constraint.SlotDurationMin.GetTimeDuration(); err == nil {
			result.MinSlotDuration = duration
		}
	}
	if constraint.SlotDurationMax != nil {
		if duration, err := constraint.SlotDurationMax.GetTimeDuration(); err == nil {
			result.MaxSlotDuration = duration
		}
	}
	if constraint.SlotDurationStepSize != nil {
		if duration, err := constraint.SlotDurationStepSize.GetTimeDuration(); err == nil {
			result.SlotDurationStepSize = duration
		}
	}

	return result
}

// send power limits to the EV
func (e *EMobilityImpl) EVWritePowerLimits(data []EVDurationSlotValue) error {
	if e.evEntity == nil || e.evTimeSeries == nil {
		return ErrNotSupported
	}

	if len(data) == 0 {
		return errors.New("missing power limit data")
	}

	constraints := e.EVGetPowerConstraints()

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	desc, err := e.evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints)
	if err != nil {
		return ErrNotSupported
	}

	timeSeriesSlots := []model.TimeSeriesSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		timeSeriesSlot := model.TimeSeriesSlotType{
			TimeSeriesSlotId: eebusUtil.Ptr(model.TimeSeriesSlotIdType(index)),
			TimePeriod: &model.TimePeriodType{
				StartTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeStart),
			},
			MaxValue: model.NewScaledNumberType(slot.Value),
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeSeriesSlot.TimePeriod.EndTime = model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeEndTime)
		}
		timeSeriesSlots = append(timeSeriesSlots, timeSeriesSlot)

		totalDuration += slot.Duration
	}

	timeSeriesData := model.TimeSeriesDataType{
		TimeSeriesId: desc.TimeSeriesId,
		TimePeriod: &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromDuration(totalDuration),
		},
		TimeSeriesSlot: timeSeriesSlots,
	}

	_, err = e.evTimeSeries.WriteValues([]model.TimeSeriesDataType{timeSeriesData})

	return err
}

// returns the minimum and maximum number of incentive slots allowed
func (e *EMobilityImpl) EVGetIncentiveConstraints() EVIncentiveSlotConstraints {
	result := EVIncentiveSlotConstraints{}

	if e.evEntity == nil || e.evIncentiveTable == nil {
		return result
	}

	constraints, err := e.evIncentiveTable.GetConstraints()
	if err != nil {
		return result
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.IncentiveSlotConstraints.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMin)
	}
	if constraint.IncentiveSlotConstraints.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMax)
	}

	return result
}

// send incentives to the EV
func (e *EMobilityImpl) EVWriteIncentives(data []EVDurationSlotValue) error {
	if e.evEntity == nil || e.evIncentiveTable == nil {
		return features.ErrDataNotAvailable
	}

	if len(data) == 0 {
		return errors.New("missing incentive data")
	}

	constraints := e.EVGetIncentiveConstraints()

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	incentiveSlots := []model.IncentiveTableIncentiveSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		timeInterval := &model.TimeTableDataType{
			StartTime: &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeStart),
			},
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeInterval.EndTime = &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeEndTime),
			}
		}

		incentiveSlot := model.IncentiveTableIncentiveSlotType{
			TimeInterval: timeInterval,
			Tier: []model.IncentiveTableTierType{
				{
					Tier: &model.TierDataType{
						TierId: eebusUtil.Ptr(model.TierIdType(1)),
					},
					Boundary: []model.TierBoundaryDataType{
						{
							BoundaryId:         eebusUtil.Ptr(model.TierBoundaryIdType(1)), // only 1 boundary exists
							LowerBoundaryValue: model.NewScaledNumberType(0),
						},
					},
					Incentive: []model.IncentiveDataType{
						{
							IncentiveId: eebusUtil.Ptr(model.IncentiveIdType(1)), // always use price
							Value:       model.NewScaledNumberType(slot.Value),
						},
					},
				},
			},
		}
		incentiveSlots = append(incentiveSlots, incentiveSlot)

		totalDuration += slot.Duration
	}

	incentiveData := model.IncentiveTableType{
		Tariff: &model.TariffDataType{
			TariffId: eebusUtil.Ptr(model.TariffIdType(0)),
		},
		IncentiveSlot: incentiveSlots,
	}

	_, err := e.evIncentiveTable.WriteValues([]model.IncentiveTableType{incentiveData})

	return err
}
