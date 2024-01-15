package grid

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/spine-go/model"
)

// return the power limitation factor
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - ErrNotSupported if getting the communication standard is not supported
//   - and others
func (g *GridImpl) PowerLimitationFactor() (float64, error) {
	if g.gridEntity == nil {
		return 0, util.ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0, features.ErrDataNotAvailable
	}

	keyname := model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor

	// check if device configuration description has curtailment limit factor key name
	_, err := g.gridDeviceConfiguration.GetDescriptionForKeyName(keyname)
	if err != nil {
		return 0, err
	}

	data, err := g.gridDeviceConfiguration.GetKeyValueForKeyName(keyname, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, features.ErrDataNotAvailable
	}

	value := data.(*model.ScaledNumberType)
	return value.GetValue(), nil
}

// return the momentary power consumption (positive) or production (negative)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) MomentaryPowerConsumptionOrProduction() (float64, error) {
	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPowerTotal
	data, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume there is only one value
	mId := data[0].MeasurementId
	value := data[0].Value
	if mId == nil || value == nil {
		return 0, features.ErrDataNotAvailable
	}

	// according to UC_TS_MonitoringOfGridConnectionPoint 3.2.2.2.4.1
	// positive values are with description "PositiveEnergyDirection" value "consume"
	// but we verify this
	desc, err := g.gridElectricalConnection.GetDescriptionForMeasurementId(*mId)
	if err != nil {
		return 0, err
	}

	// if energy direction is not consume, invert it
	if desc.PositiveEnergyDirection != nil && *desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
		return -1 * value.GetValue(), nil
	}

	return value.GetValue(), nil
}

// return the total feed-in energy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) TotalFeedInEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridFeedIn
	data, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

// return the total consumed energy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) TotalConsumedEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridConsumption
	data, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

// return the momentary current consumption (positive) or production (negative) per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) MomentaryCurrentConsumptionOrProduction() ([]float64, error) {
	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	values, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	var phaseA, phaseB, phaseC float64

	for _, item := range values {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		param, err := g.gridElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
		if err != nil || param.AcMeasuredPhases == nil {
			continue
		}

		value := item.Value.GetValue()

		// according to UC_TS_MonitoringOfGridConnectionPoint 3.2.1.3.2.4
		// positive values are with description "PositiveEnergyDirection" value "consume"
		// but we should verify this
		if desc, err := g.gridElectricalConnection.GetDescriptionForMeasurementId(*item.MeasurementId); err == nil {
			// if energy direction is not consume, invert it
			if desc.PositiveEnergyDirection != nil && *desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
				value = -1 * value
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

// return the voltage per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) Voltage() ([]float64, error) {
	measurement := model.MeasurementTypeTypeVoltage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACVoltage
	data, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	var phaseA, phaseB, phaseC float64

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		param, err := g.gridElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
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

// return the frequence
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) Frequency() (float64, error) {
	measurement := model.MeasurementTypeTypeFrequency
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACFrequency
	item, err := g.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

func (g *GridImpl) getValuesForTypeCommodityScope(measurement model.MeasurementTypeType, commodity model.CommodityTypeType, scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if g.gridEntity == nil {
		return nil, util.ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	return g.gridMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
