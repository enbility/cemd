package ucmpc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Scenario 1

// return the momentary active power consumption or production
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCMPC) MomentaryTotalPower(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPowerTotal
	data, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume there is only one value
	mId := data[0].MeasurementId
	value := data[0].Value
	if mId == nil || value == nil {
		return 0, features.ErrDataNotAvailable
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil || electricalConnection == nil {
		return 0, err
	}

	desc, err := electricalConnection.GetDescriptionForMeasurementId(*mId)
	if err != nil {
		return 0, err
	}

	// if energy direction is not consume, report an error
	if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
		return 0, features.ErrMissingData
	}

	return value.GetValue(), nil
}

// return the momentary active phase specific power consumption or production per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCMPC) MomentaryPhasePower(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPower
	values, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil || electricalConnection == nil {
		return nil, err
	}

	var phaseA, phaseB, phaseC float64

	for _, item := range values {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		param, err := electricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
		if err != nil || param.AcMeasuredPhases == nil {
			continue
		}

		value := item.Value.GetValue()

		if desc, err := electricalConnection.GetDescriptionForMeasurementId(*item.MeasurementId); err == nil {
			// if energy direction is not consume, invert it
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
				return nil, err
			}
		}

		switch *param.AcMeasuredPhases {
		case model.ElectricalConnectionPhaseNameTypeA:
			phaseA = value
		case model.ElectricalConnectionPhaseNameTypeB:
			phaseB = value
		case model.ElectricalConnectionPhaseNameTypeC:
			phaseC = value
		}
	}

	return []float64{phaseA, phaseB, phaseC}, nil
}

// Scenario 2

// return the total consumption energy
//
//   - positive values are used for consumption
func (e *UCMPC) TotalConsumedEnergy(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACEnergyConsumed
	data, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume thre is only one result
	value := data[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// return the total feed in energy
//
//   - negative values are used for production
func (e *UCMPC) TotalProducedEnergy(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACEnergyProduced
	data, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume thre is only one result
	value := data[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// Scenario 3

// return the momentary phase specific current consumption or production
//
//   - positive values are used for consumption
//   - negative values are used for production
func (e *UCMPC) MomentaryCurrents(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	values, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil || electricalConnection == nil {
		return nil, err
	}

	var phaseA, phaseB, phaseC float64

	for _, item := range values {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		param, err := electricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
		if err != nil || param.AcMeasuredPhases == nil {
			continue
		}

		value := item.Value.GetValue()

		if desc, err := electricalConnection.GetDescriptionForMeasurementId(*item.MeasurementId); err == nil {
			// if energy direction is not consume, invert it
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
				return nil, err
			}
		}

		switch *param.AcMeasuredPhases {
		case model.ElectricalConnectionPhaseNameTypeA:
			phaseA = value
		case model.ElectricalConnectionPhaseNameTypeB:
			phaseB = value
		case model.ElectricalConnectionPhaseNameTypeC:
			phaseC = value
		}
	}

	return []float64{phaseA, phaseB, phaseC}, nil
}

// Scenario 4

// return the phase specific voltage details
func (e *UCMPC) Voltages(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeVoltage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACVoltage
	data, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	electricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil || electricalConnection == nil {
		return nil, err
	}

	var phaseA, phaseB, phaseC float64

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		param, err := electricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
		if err != nil || param.AcMeasuredPhases == nil {
			continue
		}

		value := item.Value.GetValue()

		switch *param.AcMeasuredPhases {
		case model.ElectricalConnectionPhaseNameTypeA:
			phaseA = value
		case model.ElectricalConnectionPhaseNameTypeB:
			phaseB = value
		case model.ElectricalConnectionPhaseNameTypeC:
			phaseC = value
		}
	}

	return []float64{phaseA, phaseB, phaseC}, nil
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
	item, err := e.getValuesForTypeCommodityScope(entity, measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// take the first item
	value := item[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// helper

func (e *UCMPC) getValuesForTypeCommodityScope(
	entity spineapi.EntityRemoteInterface,
	measurement model.MeasurementTypeType,
	commodity model.CommodityTypeType,
	scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	measurementFeature, err := util.Measurement(e.service, entity)
	if err != nil || measurementFeature == nil {
		return nil, err
	}

	return measurementFeature.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
