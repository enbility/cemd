package grid

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
)

// return the power limitation factor
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - ErrNotSupported if getting the communication standard is not supported
//   - and others
func (g *GridImpl) PowerLimitationFactor() (float64, error) {
	if g.gridEntity == nil {
		return 0.0, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	keyname := model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor

	// check if device configuration description has curtailment limit factor key name
	support, err := g.gridDeviceConfiguration.GetDescriptionKeyNameSupport(keyname)
	if err != nil {
		return 0.0, nil
	}
	if !support {
		return 0.0, features.ErrNotSupported
	}

	data, err := g.gridDeviceConfiguration.GetValueForKeyName(keyname, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	if err != nil {
		return 0.0, err
	}

	if data == nil {
		return 0.0, features.ErrDataNotAvailable
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
	if g.gridEntity == nil {
		return 0.0, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPowerTotal
	value, mId, err := g.gridMeasurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0.0, err
	}

	// according to UC_TS_MonitoringOfGridConnectionPoint 3.2.2.2.4.1
	// positive values are with description "PositiveEnergyDirection" value "consume"
	// but we verify this
	desc, err := g.gridElectricalConnection.GetDescriptionForMeasurementId(mId)
	if err != nil {
		return 0.0, err
	}

	// if energy direction is not consume, invert it
	if desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
		return -1 * value, nil
	}

	return value, nil
}

// return the total feed-in energy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) TotalFeedInEnergy() (float64, error) {
	if g.gridEntity == nil {
		return 0.0, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridFeedIn
	value, _, err := g.gridMeasurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}

// return the total consumed energy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) TotalConsumedEnergy() (float64, error) {
	if g.gridEntity == nil {
		return 0.0, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeGridConsumption
	value, _, err := g.gridMeasurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}

// return the momentary current consumption (positive) or production (negative) per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) MomentaryCurrentConsumptionOrProduction() ([]float64, error) {
	if g.gridEntity == nil {
		return nil, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	data, measuremtIdData, err := g.gridMeasurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, g.gridElectricalConnection)
	if err != nil {
		return nil, err
	}

	var result []float64

	for _, phase := range util.PhaseMapping {
		value := 0.0
		if theValue, exists := data[phase]; exists {
			value = theValue
		}

		// according to UC_TS_MonitoringOfGridConnectionPoint 3.2.1.3.2.4
		// positive values are with description "PositiveEnergyDirection" value "consume"
		// but we verify this
		mId, exists := measuremtIdData[phase]
		if exists {
			if desc, err := g.gridElectricalConnection.GetDescriptionForMeasurementId(mId); err == nil {
				// if energy direction is not consume, invert it
				if desc.PositiveEnergyDirection != model.EnergyDirectionTypeConsume {
					value = -1 * value
				}
			}
		}

		result = append(result, value)
	}

	return result, nil
}

// return the voltage per phase
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) Voltage() ([]float64, error) {
	if g.gridEntity == nil {
		return nil, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypeVoltage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACVoltage
	data, _, err := g.gridMeasurement.GetValuesPerPhaseForTypeCommodityScope(measurement, commodity, scope, g.gridElectricalConnection)
	if err != nil {
		return nil, err
	}

	var result []float64

	for _, phase := range util.PhaseMapping {
		value := 0.0
		if theValue, exists := data[phase]; exists {
			value = theValue
		}

		result = append(result, value)
	}

	return result, nil
}

// return the frequence
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *GridImpl) Frequency() (float64, error) {
	if g.gridEntity == nil {
		return 0.0, ErrDeviceDisconnected
	}

	if g.gridMeasurement == nil {
		return 0.0, features.ErrDataNotAvailable
	}

	measurement := model.MeasurementTypeTypeFrequency
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACFrequency
	value, _, err := g.gridMeasurement.GetValueForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0.0, err
	}

	return value, nil
}
