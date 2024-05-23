package uclppserver

import (
	"errors"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the current production limit data
//
// return values:
//   - limit: load limit data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCLPPServer) ProductionLimit() (limit api.LoadLimit, resultErr error) {
	limit = api.LoadLimit{
		Value:        0.0,
		IsChangeable: false,
		IsActive:     false,
		Duration:     0,
	}
	resultErr = eebusapi.ErrDataNotAvailable

	limidId, err := e.loadControlLimitId()
	if err != nil {
		return
	}

	value := util.GetLocalLimitValueForLimitId(e.service, limidId)
	if value.LimitId == nil || value.Value == nil {
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

	return limit, nil
}

// set the current production limit data
func (e *UCLPPServer) SetProductionLimit(limit api.LoadLimit) (resultErr error) {
	resultErr = eebusapi.ErrDataNotAvailable

	limidId, err := e.loadControlLimitId()
	if err != nil {
		return
	}

	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	loadControl := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	if loadControl == nil {
		return
	}

	limitData := model.LoadControlLimitDataType{
		LimitId:           eebusutil.Ptr(limidId),
		IsLimitChangeable: eebusutil.Ptr(limit.IsChangeable),
		IsLimitActive:     eebusutil.Ptr(limit.IsActive),
		Value:             model.NewScaledNumberType(limit.Value),
	}
	if limit.Duration > 0 {
		limitData.TimePeriod = &model.TimePeriodType{
			EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(limit.Duration),
		}
	}
	// TODO: this overwrites LPC data as well
	limits := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{limitData},
	}

	loadControl.SetData(model.FunctionTypeLoadControlLimitListData, limits)

	return nil
}

// return the currently pending incoming consumption write limits
func (e *UCLPPServer) PendingProductionLimits() map[model.MsgCounterType]api.LoadLimit {
	result := make(map[model.MsgCounterType]api.LoadLimit)

	limitId, err := e.loadControlLimitId()
	if err != nil {
		return result
	}

	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	for key, msg := range e.pendingLimits {
		data := msg.Cmd.LoadControlLimitListData

		// elements are only added to the map if all required fields exist
		// therefor not check for these are needed here

		// find the item which contains the limit for this usecase
		for _, item := range data.LoadControlLimitData {
			if item.LimitId == nil ||
				limitId != *item.LimitId {
				continue
			}

			limit := api.LoadLimit{}

			if item.TimePeriod != nil {
				if duration, err := item.TimePeriod.GetDuration(); err == nil {
					limit.Duration = duration
				}
			}

			if item.IsLimitActive != nil {
				limit.IsActive = *item.IsLimitActive
			}

			if item.Value != nil {
				limit.Value = item.Value.GetValue()
			}

			result[key] = limit
		}
	}

	return result
}

// accept or deny an incoming consumption write limit
//
// use PendingProductionLimits to get the list of currently pending requests
func (e *UCLPPServer) ApproveOrDenyProductionLimit(msgCounter model.MsgCounterType, approve bool, reason string) {
	e.pendingMux.Lock()
	defer e.pendingMux.Unlock()

	msg, ok := e.pendingLimits[msgCounter]
	if !ok {
		return
	}

	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	f := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	result := model.ErrorType{
		ErrorNumber: model.ErrorNumberType(0),
	}
	if !approve {
		result.ErrorNumber = model.ErrorNumberType(7)
		result.Description = eebusutil.Ptr(model.DescriptionType(reason))
	}
	f.ApproveOrDenyWrite(msg, result)
}

// Scenario 2

// return Failsafe limit for the produced active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *UCLPPServer) FailsafeProductionActivePowerLimit() (limit float64, isChangeable bool, resultErr error) {
	limit = 0
	isChangeable = false
	resultErr = eebusapi.ErrDataNotAvailable

	keyName := model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit
	keyData := util.GetLocalDeviceConfigurationKeyValueForKeyName(e.service, keyName)
	if keyData.KeyId == nil || keyData.Value == nil || keyData.Value.ScaledNumber == nil {
		return
	}

	limit = keyData.Value.ScaledNumber.GetValue()
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set Failsafe limit for the produced active (real) power of the
// Controllable System. This limit becomes activated in "init" state or "failsafe state".
func (e *UCLPPServer) SetFailsafeProductionActivePowerLimit(value float64, changeable bool) error {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit
	keyValue := model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(value),
	}

	return util.SetLocalDeviceConfigurationKeyValue(e.service, keyName, changeable, keyValue)
}

// return minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
func (e *UCLPPServer) FailsafeDurationMinimum() (duration time.Duration, isChangeable bool, resultErr error) {
	duration = 0
	isChangeable = false
	resultErr = eebusapi.ErrDataNotAvailable

	keyName := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum
	keyData := util.GetLocalDeviceConfigurationKeyValueForKeyName(e.service, keyName)
	if keyData.KeyId == nil || keyData.Value == nil || keyData.Value.Duration == nil {
		return
	}

	durationValue, err := keyData.Value.Duration.GetTimeDuration()
	if err != nil {
		return
	}

	duration = durationValue
	isChangeable = (keyData.IsValueChangeable != nil && *keyData.IsValueChangeable)
	resultErr = nil
	return
}

// set minimum time the Controllable System remains in "failsafe state" unless conditions
// specified in this Use Case permit leaving the "failsafe state"
//
// parameters:
//   - duration: has to be >= 2h and <= 24h
//   - changeable: boolean if the client service can change this value
func (e *UCLPPServer) SetFailsafeDurationMinimum(duration time.Duration, changeable bool) error {
	if duration < time.Duration(time.Hour*2) || duration > time.Duration(time.Hour*24) {
		return errors.New("duration outside of allowed range")
	}
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum
	keyValue := model.DeviceConfigurationKeyValueValueType{
		Duration: model.NewDurationType(duration),
	}

	return util.SetLocalDeviceConfigurationKeyValue(e.service, keyName, changeable, keyValue)
}

// Scenario 3

func (e *UCLPPServer) IsHeartbeatWithinDuration() bool {
	if e.heartbeatDiag == nil {
		return false
	}

	return e.heartbeatDiag.IsHeartbeatWithinDuration(2 * time.Minute)
}

// Scenario 4

// return nominal maximum active (real) power the Controllable System is
// allowed to produce due to the customer's contract.
func (e *UCLPPServer) ContractualProductionNominalMax() (value float64, resultErr error) {
	value = 0
	resultErr = eebusapi.ErrDataNotAvailable

	charData := util.GetLocalElectricalConnectionCharacteristicForContextType(
		e.service,
		model.ElectricalConnectionCharacteristicContextTypeEntity,
		model.ElectricalConnectionCharacteristicTypeTypeContractualProductionNominalMax,
	)
	if charData.CharacteristicId == nil || charData.Value == nil {
		return
	}

	return charData.Value.GetValue(), nil
}

// set nominal maximum active (real) power the Controllable System is
// allowed to produce due to the customer's contract.
func (e *UCLPPServer) SetContractualProductionNominalMax(value float64) error {
	return util.SetLocalElectricalConnectionCharacteristicForContextType(
		e.service,
		model.ElectricalConnectionCharacteristicContextTypeEntity,
		model.ElectricalConnectionCharacteristicTypeTypeContractualProductionNominalMax,
		value,
	)
}
