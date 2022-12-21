package emobility

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
	eebusUtil "github.com/enbility/eebus-go/util"
)

// return the current charge sate of the EV
func (e *EMobilityImpl) EVCurrentChargeState() (EVChargeStateType, error) {
	if e.evEntity == nil {
		return EVChargeStateTypeUnplugged, nil
	}

	if e.evDeviceDiagnosis == nil {
		return EVChargeStateTypeUnknown, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return 0, ErrEVDisconnected
	}

	if e.evElectricalConnection == nil {
		return 0, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return 0, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return 0, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return nil, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return nil, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return nil, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return nil, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return nil, nil, nil, ErrEVDisconnected
	}

	if e.evElectricalConnection == nil {
		return nil, nil, nil, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return EVCommunicationStandardTypeUnknown, ErrEVDisconnected
	}

	if e.evDeviceConfiguration == nil {
		return EVCommunicationStandardTypeUnknown, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return false, ErrEVDisconnected
	}

	if e.evLoadControl == nil {
		return false, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return false, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return false, features.ErrDataNotAvailable
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
	if e.evEntity == nil {
		return 0, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return 0, features.ErrDataNotAvailable
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
