package uclpp

import (
	"errors"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the current loadcontrol limit data
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCLPP) ProductionLimit(entity spineapi.EntityRemoteInterface) (
	limit api.LoadLimit, resultErr error) {
	limit = api.LoadLimit{
		Value:        0.0,
		IsChangeable: false,
		IsActive:     false,
	}

	resultErr = api.ErrNoCompatibleEntity
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return
	}

	resultErr = eebusapi.ErrDataNotAvailable
	loadControl, err := util.LoadControl(e.service, entity)
	if err != nil || loadControl == nil {
		return
	}

	limitDescriptions, err := loadControl.GetLimitDescriptionsForTypeCategoryDirectionScope(
		model.LoadControlLimitTypeTypeSignDependentAbsValueLimit,
		model.LoadControlCategoryTypeObligation,
		model.EnergyDirectionTypeProduce,
		model.ScopeTypeTypeActivePowerLimit)
	if err != nil || len(limitDescriptions) != 1 {
		return
	}

	value, err := loadControl.GetLimitValueForLimitId(*limitDescriptions[0].LimitId)
	if err != nil || value.Value == nil {
		return
	}

	limit.Value = value.Value.GetValue()
	limit.IsChangeable = (value.IsLimitChangeable != nil && *value.IsLimitChangeable)
	limit.IsActive = (value.IsLimitActive != nil && *value.IsLimitActive)
	if value.TimePeriod != nil && value.TimePeriod.EndTime != nil {
		if duration, err := value.TimePeriod.EndTime.GetTimeDuration(); err == nil {
			limit.Duration = duration
		}
	}

	resultErr = nil

	return
}

// send new LoadControlLimits
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - limit: load limit data
func (e *UCLPP) WriteProductionLimit(
	entity spineapi.EntityRemoteInterface,
	limit api.LoadLimit) (*model.MsgCounterType, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	loadControl, err := util.LoadControl(e.service, entity)
	if err != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	limitDescriptions, err := loadControl.GetLimitDescriptionsForTypeCategoryDirectionScope(
		model.LoadControlLimitTypeTypeSignDependentAbsValueLimit,
		model.LoadControlCategoryTypeObligation,
		model.EnergyDirectionTypeProduce,
		model.ScopeTypeTypeActivePowerLimit,
	)
	if err != nil ||
		len(limitDescriptions) != 1 ||
		limitDescriptions[0].LimitId == nil {
		return nil, eebusapi.ErrMetadataNotAvailable
	}

	limitDesc := limitDescriptions[0]

	if _, err := loadControl.GetLimitValueForLimitId(*limitDesc.LimitId); err != nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	currentLimits, err := loadControl.GetLimitValues()
	if err != nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	for index, item := range currentLimits {
		if item.LimitId == nil ||
			*item.LimitId != *limitDesc.LimitId {
			continue
		}

		// EEBus_UC_TS_LimitationOfPowerProduction V1.0.0 3.2.2.2.2.2
		// If set to "true", the timePeriod, value and isLimitActive Elements SHALL be writeable by a client.
		if item.IsLimitChangeable != nil && !*item.IsLimitChangeable {
			return nil, eebusapi.ErrNotSupported
		}

		newLimit := model.LoadControlLimitDataType{
			LimitId:       limitDesc.LimitId,
			IsLimitActive: eebusutil.Ptr(limit.IsActive),
			Value:         model.NewScaledNumberType(limit.Value),
		}
		if limit.Duration > 0 {
			newLimit.TimePeriod = &model.TimePeriodType{
				EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(limit.Duration),
			}
		}

		currentLimits[index] = newLimit
		break
	}

	msgCounter, err := loadControl.WriteLimitValues(limitData)

	return msgCounter, err
}

// Scenario 2

// return Failsafe limit for the produced active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *UCLPP) FailsafeProductionActivePowerLimit(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	data, err := deviceConfiguration.GetKeyValueForKeyName(keyname, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	if err != nil || data == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	value, ok := data.(*model.ScaledNumberType)
	if !ok || value == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// send new Failsafe Production Active Power Limit
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - value: the new limit in W
func (e *UCLPP) WriteFailsafeProductionActivePowerLimit(entity spineapi.EntityRemoteInterface, value float64) (*model.MsgCounterType, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil || deviceConfiguration == nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	data, err := deviceConfiguration.GetDescriptionForKeyName(keyname)
	if err != nil || data == nil || data.KeyId == nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	keyData := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: data.KeyId,
			Value: &model.DeviceConfigurationKeyValueValueType{
				ScaledNumber: model.NewScaledNumberType(value),
			},
		},
	}

	msgCounter, err := deviceConfiguration.WriteKeyValues(keyData)

	return msgCounter, err
}

// return minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
func (e *UCLPP) FailsafeDurationMinimum(entity spineapi.EntityRemoteInterface) (time.Duration, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	data, err := deviceConfiguration.GetKeyValueForKeyName(keyname, model.DeviceConfigurationKeyValueTypeTypeDuration)
	if err != nil || data == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	value, ok := data.(*model.DurationType)
	if !ok || value == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	return value.GetTimeDuration()
}

// send new Failsafe Duration Minimum
//
// parameters:
//   - entity: the entity of the e.g. EVSE
//   - duration: the duration, between 2h and 24h
func (e *UCLPP) WriteFailsafeDurationMinimum(entity spineapi.EntityRemoteInterface, duration time.Duration) (*model.MsgCounterType, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	if duration < time.Duration(time.Hour*2) || duration > time.Duration(time.Hour*24) {
		return nil, errors.New("duration outside of allowed range")
	}

	keyname := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil || deviceConfiguration == nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	data, err := deviceConfiguration.GetDescriptionForKeyName(keyname)
	if err != nil || data == nil || data.KeyId == nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	keyData := []model.DeviceConfigurationKeyValueDataType{
		{
			KeyId: data.KeyId,
			Value: &model.DeviceConfigurationKeyValueValueType{
				Duration: model.NewDurationType(duration),
			},
		},
	}

	msgCounter, err := deviceConfiguration.WriteKeyValues(keyData)

	return msgCounter, err
}

// Scenario 4

// return nominal maximum active (real) power the Controllable System is
// able to produce according to the device label or data sheet.
func (e *UCLPP) PowerProductionNominalMax(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil || electricalConnection == nil {
		return 0, err
	}

	data, err := electricalConnection.GetCharacteristicForContextType(
		model.ElectricalConnectionCharacteristicContextTypeEntity,
		model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax,
	)
	if err != nil || data.Value == nil {
		return 0, err
	}

	return data.Value.GetValue(), nil
}
