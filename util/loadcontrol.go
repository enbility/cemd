package util

import (
	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the current loadcontrol limits for a categoriy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func LoadControlLimits(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	entityType model.EntityTypeType,
	category model.LoadControlCategoryType) ([]float64, error) {
	if entity == nil || entity.EntityType() != entityType {
		return nil, api.ErrNoCompatibleEntity
	}

	evLoadControl, err := LoadControl(service, entity)
	evElectricalConnection, err2 := ElectricalConnection(service, entity)
	if err != nil || err2 != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	// find out the appropriate limitId for each phase value
	// limitDescription contains the measurementId for each limitId
	limitDescriptions, err := evLoadControl.GetLimitDescriptionsForCategory(category)
	if err != nil {
		return nil, features.ErrDataNotAvailable
	}

	var result []float64

	for i := 0; i < 3; i++ {
		phaseName := PhaseNameMapping[i]

		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.MeasurementId == nil {
			// there is no data for this phase, the phase may not exit
			result = append(result, 0)
			continue
		}

		var limitDesc *model.LoadControlLimitDescriptionDataType
		for _, desc := range limitDescriptions {
			if desc.MeasurementId != nil &&
				elParamDesc.MeasurementId != nil &&
				*desc.MeasurementId == *elParamDesc.MeasurementId {
				safeDesc := desc
				limitDesc = &safeDesc
				break
			}
		}

		if limitDesc == nil || limitDesc.LimitId == nil {
			return nil, features.ErrDataNotAvailable
		}

		limitIdData, err := evLoadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
		if err != nil {
			return nil, features.ErrDataNotAvailable
		}

		var limitValue float64
		if limitIdData.Value == nil || (limitIdData.IsLimitActive != nil && !*limitIdData.IsLimitActive) {
			// report maximum possible if no limit is available or the limit is not active
			_, dataMax, _, err := evElectricalConnection.GetLimitsForParameterId(*elParamDesc.ParameterId)
			if err != nil {
				return nil, features.ErrDataNotAvailable
			}

			limitValue = dataMax
		} else {
			limitValue = limitIdData.Value.GetValue()
		}

		result = append(result, limitValue)
	}

	return result, nil
}

// generic helper to be used in UCOPEV & UCOSCEV
// send new LoadControlLimits to the remote EV
//
// parameters:
//   - limits: a set of limits for a  given limit category containing phase specific limit data
//
// category obligations:
// Sets a maximum A limit for each phase that the EV may not exceed.
// Mainly used for implementing overload protection of the site or limiting the
// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
// and with ISO15118 if the EV does not support the Optimization of Self Consumption
// usecase.
//
// category recommendations:
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
func WriteLoadControlLimits(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	entityType model.EntityTypeType,
	category model.LoadControlCategoryType,
	limits []api.LoadLimitsPhase) (*model.MsgCounterType, error) {
	if entity == nil || entity.EntityType() != entityType {
		return nil, api.ErrNoCompatibleEntity
	}

	evLoadControl, err := LoadControl(service, entity)
	evElectricalConnection, err2 := ElectricalConnection(service, entity)
	if err != nil || err2 != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	for _, phaseLimit := range limits {
		// find out the appropriate limitId for each phase value
		// limitDescription contains the measurementId for each limitId
		limitDescriptions, err := evLoadControl.GetLimitDescriptionsForCategory(category)
		if err != nil {
			continue
		}

		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseLimit.Phase)
		if err != nil || elParamDesc.MeasurementId == nil {
			continue
		}

		var limitDesc *model.LoadControlLimitDescriptionDataType
		for _, desc := range limitDescriptions {
			if desc.MeasurementId != nil &&
				elParamDesc.MeasurementId != nil &&
				*desc.MeasurementId == *elParamDesc.MeasurementId {
				safeDesc := desc
				limitDesc = &safeDesc
				break
			}
		}

		if limitDesc == nil || limitDesc.LimitId == nil {
			continue
		}

		limitIdData, err := evLoadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
		if err != nil {
			continue
		}

		// EEBus_UC_TS_OverloadProtectionByEvChargingCurrentCurtailment V1.01b 3.2.1.2.2.2
		// If omitted or set to "true", the timePeriod, value and isLimitActive element SHALL be writeable by a client.
		if limitIdData.IsLimitChangeable != nil && !*limitIdData.IsLimitChangeable {
			continue
		}

		// electricalPermittedValueSet contains the allowed min, max and the default values per phase
		limit := evElectricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(phaseLimit.Value, *elParamDesc.ParameterId)

		newLimit := model.LoadControlLimitDataType{
			LimitId:       limitDesc.LimitId,
			IsLimitActive: eebusutil.Ptr(phaseLimit.IsActive),
			Value:         model.NewScaledNumberType(limit),
		}
		limitData = append(limitData, newLimit)
	}

	msgCounter, err := evLoadControl.WriteLimitValues(limitData)

	return msgCounter, err
}
