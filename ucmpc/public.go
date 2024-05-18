package ucmpc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the momentary active power consumption or production
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCMPC) Power(entity spineapi.EntityRemoteInterface) (float64, error) {
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
		return 0, eebusapi.ErrDataNotAvailable
	}
	return values[0], nil
}

// return the momentary active phase specific power consumption or production per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCMPC) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	return util.MeasurementValuesForTypeCommodityScope(
		e.service,
		entity,
		model.MeasurementTypeTypePower,
		model.CommodityTypeTypeElectricity,
		model.ScopeTypeTypeACPower,
		model.EnergyDirectionTypeConsume,
		util.PhaseNameMapping,
	)
}

// Scenario 2

// return the total consumption energy
//
//   - positive values are used for consumption
func (e *UCMPC) EnergyConsumed(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACEnergyConsumed
	values, err := util.GetValuesForTypeCommodityScope(e.service, entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, eebusapi.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// return the total feed in energy
//
//   - negative values are used for production
func (e *UCMPC) EnergyProduced(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACEnergyProduced
	values, err := util.GetValuesForTypeCommodityScope(e.service, entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}
	if len(values) == 0 {
		return 0, eebusapi.ErrDataNotAvailable
	}

	// we assume thre is only one result
	value := values[0].Value
	if value == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// Scenario 3

// return the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *UCMPC) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
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

// Scenario 4

// return the phase specific voltage details
func (e *UCMPC) VoltagePerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
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

// Scenario 5

// return frequency
func (e *UCMPC) Frequency(entity spineapi.EntityRemoteInterface) (float64, error) {
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
		return 0, eebusapi.ErrDataNotAvailable
	}

	// take the first item
	value := values[0].Value

	if value == nil {
		return 0, eebusapi.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}
