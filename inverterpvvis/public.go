package inverterpvvis

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
)

// return the current photovoltaic production power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterPVVisImpl) CurrentProductionPower() (float64, error) {
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

// return the nominal photovoltaic peak power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterPVVisImpl) NominalPeakPower() (float64, error) {
	if i.inverterEntity == nil {
		return 0, util.ErrDeviceDisconnected
	}

	if i.inverterDeviceConfiguration == nil {
		return 0, features.ErrDataNotAvailable
	}

	_, err := i.inverterDeviceConfiguration.GetDescriptionForKeyName(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem)
	if err != nil {
		return 0, err
	}

	data, err := i.inverterDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
	if err != nil {
		return 0, err
	}

	if data == nil {
		return 0, features.ErrDataNotAvailable
	}

	value := data.(*model.ScaledNumberType)

	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}

// return the total photovoltaic yield (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (i *InverterPVVisImpl) TotalPVYield() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACYieldTotal
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

func (i *InverterPVVisImpl) getValuesForTypeCommodityScope(measurement model.MeasurementTypeType, commodity model.CommodityTypeType, scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if i.inverterEntity == nil {
		return nil, util.ErrDeviceDisconnected
	}

	if i.inverterMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	return i.inverterMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
