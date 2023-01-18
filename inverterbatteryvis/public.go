package inverterbatteryvis

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
)

// return the current battery (dis-)charge power (W)
//
//   - positive values charge power
//   - negative values discharge power
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterBatteryVisImpl) CurrentDisChargePower() (float64, error) {
	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPowerTotal

	data, err := i.getValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume there is only one value
	mId := data[0].MeasurementId
	value := data[0].Value
	if mId == nil || value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// return the total charge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterBatteryVisImpl) TotalChargeEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeCharge
	data, err := i.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

// return the total discharge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterBatteryVisImpl) TotalDischargeEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeDischarge
	data, err := i.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

// return the current state of charge in %
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterBatteryVisImpl) CurrentStateOfCharge() (float64, error) {
	measurement := model.MeasurementTypeTypePercentage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeStateOfCharge
	data, err := i.getValuesForTypeCommodityScope(measurement, commodity, scope)
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

// helper

func (i *InverterBatteryVisImpl) getValuesForTypeCommodityScope(measurement model.MeasurementTypeType, commodity model.CommodityTypeType, scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if i.inverterEntity == nil {
		return nil, util.ErrDeviceDisconnected
	}

	if i.inverterMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	return i.inverterMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
