package util

import (
	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

// Check the payload data if it contains measurementId values for a given scope
func LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(
	localServer bool,
	service eebusapi.ServiceInterface,
	payload spineapi.EventPayload,
	limitType model.LoadControlLimitTypeType,
	limitCategory model.LoadControlCategoryType,
	direction model.EnergyDirectionType,
	scope model.ScopeTypeType) bool {
	var descs []model.LoadControlLimitDescriptionDataType

	if payload.Data == nil {
		return false
	}

	limits := payload.Data.(*model.LoadControlLimitListDataType)

	if localServer {
		descs = GetLocalLimitDescriptionsForTypeCategoryDirectionScope(service, limitType, limitCategory, direction, scope)
	} else {
		loadcontrolF, err := LoadControl(service, payload.Entity)
		if err != nil {
			return false
		}

		descs, err = loadcontrolF.GetLimitDescriptionsForTypeCategoryDirectionScope(limitType, limitCategory, direction, scope)
		if err != nil {
			return false
		}
	}

	for _, item := range descs {
		if item.LimitId == nil {
			continue
		}

		for _, limit := range limits.LoadControlLimitData {
			if limit.LimitId != nil &&
				*limit.LimitId == *item.LimitId &&
				limit.Value != nil {
				return true
			}
		}
	}

	return false
}

// return the current loadcontrol limits for a categoriy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func LoadControlLimits(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	entityTypes []model.EntityTypeType,
	limitType model.LoadControlLimitTypeType,
	limitCategory model.LoadControlCategoryType,
	scopeType model.ScopeTypeType) (limits []api.LoadLimitsPhase, resultErr error) {
	limits = nil
	resultErr = api.ErrNoCompatibleEntity
	if entity == nil || !IsCompatibleEntity(entity, entityTypes) {
		return
	}

	evLoadControl, err := LoadControl(service, entity)
	evElectricalConnection, err2 := ElectricalConnection(service, entity)
	if err != nil || err2 != nil {
		return
	}

	resultErr = eebusapi.ErrDataNotAvailable
	// find out the appropriate limitId for each phase value
	// limitDescription contains the measurementId for each limitId
	limitDescriptions, err := evLoadControl.GetLimitDescriptionsForTypeCategoryDirectionScope(
		limitType, limitCategory, "", scopeType)
	if err != nil {
		return
	}

	var result []api.LoadLimitsPhase

	for i := 0; i < len(PhaseNameMapping); i++ {
		phaseName := PhaseNameMapping[i]

		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.MeasurementId == nil {
			// there is no data for this phase, the phase may not exist
			result = append(result, api.LoadLimitsPhase{Phase: phaseName})
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
			return
		}

		limitIdData, err := evLoadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
		if err != nil {
			return
		}

		var limitValue float64
		if limitIdData.Value == nil || (limitIdData.IsLimitActive != nil && !*limitIdData.IsLimitActive) {
			// report maximum possible if no limit is available or the limit is not active
			_, dataMax, _, err := evElectricalConnection.GetLimitsForParameterId(*elParamDesc.ParameterId)
			if err != nil {
				return
			}

			limitValue = dataMax
		} else {
			limitValue = limitIdData.Value.GetValue()
		}

		newLimit := api.LoadLimitsPhase{
			Phase:        phaseName,
			IsChangeable: (limitIdData.IsLimitChangeable != nil && *limitIdData.IsLimitChangeable),
			IsActive:     (limitIdData.IsLimitActive != nil && *limitIdData.IsLimitActive),
			Value:        limitValue,
		}

		result = append(result, newLimit)
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
	entityTypes []model.EntityTypeType,
	category model.LoadControlCategoryType,
	limits []api.LoadLimitsPhase) (*model.MsgCounterType, error) {
	if entity == nil || !IsCompatibleEntity(entity, entityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	loadControl, err := LoadControl(service, entity)
	electricalConnection, err2 := ElectricalConnection(service, entity)
	if err != nil || err2 != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	for _, phaseLimit := range limits {
		// find out the appropriate limitId for each phase value
		// limitDescription contains the measurementId for each limitId
		limitDescriptions, err := loadControl.GetLimitDescriptionsForCategory(category)
		if err != nil {
			continue
		}

		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := electricalConnection.GetParameterDescriptionForMeasuredPhase(phaseLimit.Phase)
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

		limitIdData, err := loadControl.GetLimitValueForLimitId(*limitDesc.LimitId)
		if err != nil {
			continue
		}

		// EEBus_UC_TS_OverloadProtectionByEvChargingCurrentCurtailment V1.01b 3.2.1.2.2.2
		// If omitted or set to "true", the timePeriod, value and isLimitActive element SHALL be writeable by a client.
		if limitIdData.IsLimitChangeable != nil && !*limitIdData.IsLimitChangeable {
			continue
		}

		// electricalPermittedValueSet contains the allowed min, max and the default values per phase
		limit := electricalConnection.AdjustValueToBeWithinPermittedValuesForParameter(phaseLimit.Value, *elParamDesc.ParameterId)

		newLimit := model.LoadControlLimitDataType{
			LimitId:       limitDesc.LimitId,
			IsLimitActive: eebusutil.Ptr(phaseLimit.IsActive),
			Value:         model.NewScaledNumberType(limit),
		}
		limitData = append(limitData, newLimit)
	}

	msgCounter, err := loadControl.WriteLimitValues(limitData)

	return msgCounter, err
}

func GetLocalLimitDescriptionsForTypeCategoryDirectionScope(
	service eebusapi.ServiceInterface,
	limitType model.LoadControlLimitTypeType,
	limitCategory model.LoadControlCategoryType,
	limitDirection model.EnergyDirectionType,
	scopeType model.ScopeTypeType,
) (descriptions []model.LoadControlLimitDescriptionDataType) {
	descriptions = []model.LoadControlLimitDescriptionDataType{}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	loadControl := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	if loadControl == nil {
		return
	}

	data, err := spine.LocalFeatureDataCopyOfType[*model.LoadControlLimitDescriptionListDataType](
		loadControl, model.FunctionTypeLoadControlLimitDescriptionListData)
	if err != nil || data == nil || data.LoadControlLimitDescriptionData == nil {
		return
	}

	for _, desc := range data.LoadControlLimitDescriptionData {
		if desc.LimitId != nil &&
			(limitType == "" || (desc.LimitType != nil && *desc.LimitType == limitType)) &&
			(limitCategory == "" || (desc.LimitCategory != nil && *desc.LimitCategory == limitCategory)) &&
			(limitDirection == "" || (desc.LimitDirection != nil && *desc.LimitDirection == limitDirection)) &&
			(scopeType == "" || (desc.ScopeType != nil && *desc.ScopeType == scopeType)) {
			descriptions = append(descriptions, desc)
		}
	}

	return descriptions
}

func GetLocalLimitValueForLimitId(
	service eebusapi.ServiceInterface,
	limitId model.LoadControlLimitIdType,
) (value model.LoadControlLimitDataType) {
	value = model.LoadControlLimitDataType{}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	loadControl := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	if loadControl == nil {
		return
	}

	values, err := spine.LocalFeatureDataCopyOfType[*model.LoadControlLimitListDataType](
		loadControl, model.FunctionTypeLoadControlLimitListData)
	if err != nil || values == nil || values.LoadControlLimitData == nil {
		return
	}

	for _, item := range values.LoadControlLimitData {
		if item.LimitId != nil && *item.LimitId == limitId {
			value = item
			break
		}
	}

	return
}
