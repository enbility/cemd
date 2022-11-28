package emobility

import (
	"github.com/DerAndereAndi/eebus-go-cem/util"
	"github.com/DerAndereAndi/eebus-go/features"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

var phaseMapping = []string{"a", "b", "c"}

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

	switch diagnosisState.OperatingState {
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

	return e.evElectricalConnection.GetConnectedPhases()
}

// return the charged energy measurement in Wh of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVChargedEnergy() (float64, error) {
	if e.evEntity == nil {
		return 0.0, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	return e.evMeasurement.GetValueForScope(model.ScopeTypeTypeCharge, e.evElectricalConnection)
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

	data, err := e.evMeasurement.GetValuesPerPhaseForScope(model.ScopeTypeTypeACPower, e.evElectricalConnection)
	if err != nil {
		return nil, err
	}

	var result []float64

	for _, phase := range phaseMapping {
		value := 0.0
		if theValue, exists := data[phase]; exists {
			value = theValue
		}
		result = append(result, value)
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

	data, err := e.evMeasurement.GetValuesPerPhaseForScope(model.ScopeTypeTypeACCurrent, e.evElectricalConnection)
	if err != nil {
		return nil, err
	}

	var result []float64

	for _, phase := range phaseMapping {
		value := 0.0
		if theValue, exists := data[phase]; exists {
			value = theValue
		}
		result = append(result, value)
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

	dataMin, dataMax, dataDefault, err := e.evElectricalConnection.GetCurrentsLimits()
	if err != nil {
		return nil, nil, nil, err
	}

	var resultMin, resultMax, resultDefault []float64

	for _, phase := range phaseMapping {
		value := 0.0
		if theValue, exists := dataMin[phase]; exists {
			value = theValue
		}
		resultMin = append(resultMin, value)

		value = 0.0
		if theValue, exists := dataMax[phase]; exists {
			value = theValue
		}
		resultMax = append(resultMax, value)

		value = 0.0
		if theValue, exists := dataDefault[phase]; exists {
			value = theValue
		}
		resultDefault = append(resultDefault, value)
	}

	// Min current data should be derived from min power data
	// but as this value is only properly provided via VAS the
	// currrent min values can not be trusted.
	// Min current for 3-phase should be at least 2.2A (ISO)
	for index, item := range resultMin {
		if item < 2.2 {
			resultMin[index] = 2.2
		}
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
// note:
// For obligations to work for optimizing solar excess power, the EV needs to
// have an energy demand. Recommendations work even if the EV does not have an active
// energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
// In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific
// and needs to have specific EVSE support for the specific EV brand.
// In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
func (e *EMobilityImpl) EVWriteLoadControlLimits(obligations, recommendations []float64) error {
	if e.evEntity == nil {
		return ErrEVDisconnected
	}

	if e.evElectricalConnection == nil || e.evLoadControl == nil {
		return features.ErrDataNotAvailable
	}

	electricalDesc, _, err := e.evElectricalConnection.GetParamDescriptionListData()
	if err != nil {
		return features.ErrMetadataNotAvailable
	}

	elLimits, err := e.evElectricalConnection.GetEVLimitValues()
	if err != nil {
		return features.ErrMetadataNotAvailable
	}

	limitDesc, err := e.evLoadControl.GetLimitDescription()
	if err != nil {
		return err
	}

	currentLimits, err := e.evLoadControl.GetLimitValues()
	if err != nil {
		return err
	}

	var limitData []model.LoadControlLimitDataType

	for scopeTypes := 0; scopeTypes < 2; scopeTypes++ {
		category := model.LoadControlCategoryTypeObligation
		currentsPerPhase := obligations
		if scopeTypes == 1 {
			category = model.LoadControlCategoryTypeRecommendation
			currentsPerPhase = recommendations
		}

		for index, limit := range currentsPerPhase {
			phase := phaseMapping[index]

			var limitId *model.LoadControlLimitIdType
			var elConnectionid *model.ElectricalConnectionIdType

			for _, lDesc := range limitDesc {
				if lDesc.LimitCategory == nil || lDesc.MeasurementId == nil {
					continue
				}

				if *lDesc.LimitCategory != category {
					continue
				}

				elDesc, exists := electricalDesc[*lDesc.MeasurementId]
				if !exists {
					continue
				}
				if elDesc.ElectricalConnectionId == nil || elDesc.AcMeasuredPhases == nil || string(*elDesc.AcMeasuredPhases) != phase {
					continue
				}

				limitId = lDesc.LimitId
				elConnectionid = elDesc.ElectricalConnectionId
				break
			}

			if limitId == nil || elConnectionid == nil {
				continue
			}

			var currentLimitsForID features.LoadControlLimitType
			var found bool
			for _, item := range currentLimits {
				if uint(*limitId) != item.LimitId {
					continue
				}
				currentLimitsForID = item
				found = true
				break
			}
			if !found || !currentLimitsForID.IsChangeable {
				continue
			}

			limitValue := model.NewScaledNumberType(limit)
			for _, elLimit := range elLimits {
				if elLimit.ConnectionID != uint(*elConnectionid) {
					continue
				}
				if elLimit.Scope != model.ScopeTypeTypeACCurrent {
					continue
				}
				if limit < elLimit.Min {
					limitValue = model.NewScaledNumberType(elLimit.Default)
				}
				if limit > elLimit.Max {
					limitValue = model.NewScaledNumberType(elLimit.Max)
				}
			}

			active := true
			newLimit := model.LoadControlLimitDataType{
				LimitId:       limitId,
				IsLimitActive: &active,
				Value:         limitValue,
			}
			limitData = append(limitData, newLimit)
		}
	}

	_, err = e.evLoadControl.WriteLimitValues(limitData)

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
	support, err := e.evDeviceConfiguration.GetDescriptionKeyNameSupport(model.DeviceConfigurationKeyNameTypeCommunicationsStandard)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}
	if !support {
		return EVCommunicationStandardTypeUnknown, features.ErrNotSupported
	}

	data, err := e.evDeviceConfiguration.GetEVCommunicationStandard()
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}

	return EVCommunicationStandardType(*data), err
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
		if identification.Identifier != "" {
			return identification.Identifier, nil
		}
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
	support, err := e.evLoadControl.GetLimitDescriptionCategorySupport(model.LoadControlCategoryTypeRecommendation)
	if err != nil {
		return false, err
	}
	return support, nil
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
	desc, err := e.evMeasurement.GetDescriptionForScope(model.ScopeTypeTypeStateOfCharge)
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
		return 0.0, ErrEVDisconnected
	}

	if e.evMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	// check if the SoC is supported
	support, err := e.EVSoCSupported()
	if err != nil {
		return 0.0, err
	}
	if !support {
		return 0.0, features.ErrNotSupported
	}

	return e.evMeasurement.GetSoC()
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
