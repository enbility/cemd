package ucmgcp

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the current power limitation factor
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCMGCP) PowerLimitationFactor(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement, err := util.Measurement(e.service, entity)
	if err != nil || measurement == nil {
		return 0, err
	}

	keyname := model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil || deviceConfiguration == nil {
		return 0, err
	}

	// check if device configuration description has curtailment limit factor key name
	_, err = deviceConfiguration.GetDescriptionForKeyName(keyname)
	if err != nil {
		return 0, err
	}

	data, err := deviceConfiguration.GetKeyValueForKeyName(keyname, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, features.ErrDataNotAvailable
	}

	value, ok := data.(*model.ScaledNumberType)
	if !ok || value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// Scenario 2

// return the momentary power consumption or production at the grid connection point
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *UCMGCP) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	values, err := util.MeasurementValuesForTypeCommodityScope(
		e.service,
		entity,
		model.MeasurementTypeTypePower,
		model.CommodityTypeTypeElectricity,
		model.ScopeTypeTypeACPowerTotal,
		model.EnergyDirectionTypeConsume,
		nil,
	)
	if err != nil {
		return 0, err
	}
	if len(values) != 1 {
		return 0, features.ErrDataNotAvailable
	}

	return values[0], nil
}

// Scenario 3

// return the total feed in energy at the grid connection point
//
//   - negative values are used for production
func (e *UCMGCP) EnergyFeedIn(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridFeedIn
	values, err := util.GetValuesForTypeCommodityScope(e.service, entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, features.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// Scenario 4

// return the total consumption energy at the grid connection point
//
//   - positive values are used for consumption
func (e *UCMGCP) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridConsumption
	values, err := util.GetValuesForTypeCommodityScope(e.service, entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, features.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// Scenario 5

// return the momentary current consumption or production at the grid connection point
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *UCMGCP) CurrentsPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	return util.MeasurementValuesForTypeCommodityScope(
		e.service,
		entity,
		model.MeasurementTypeTypeCurrent,
		model.CommodityTypeTypeElectricity,
		model.ScopeTypeTypeACCurrent,
		model.EnergyDirectionTypeConsume,
		util.PhaseNameMapping,
	)
}

// Scenario 6

// return the voltage phase details at the grid connection point
func (e *UCMGCP) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	return util.MeasurementValuesForTypeCommodityScope(
		e.service,
		entity,
		model.MeasurementTypeTypeVoltage,
		model.CommodityTypeTypeElectricity,
		model.ScopeTypeTypeACVoltage,
		"",
		util.PhaseNameMapping,
	)
}

// Scenario 7

// return frequency at the grid connection point
func (e *UCMGCP) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeFrequency
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACFrequency
	values, err := util.GetValuesForTypeCommodityScope(e.service, entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, features.ErrDataNotAvailable
	}

	// take the first item
	value := values[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}
