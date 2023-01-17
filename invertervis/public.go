package invertervis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine/model"
)

// return the current battery system (dis)charge power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *InverterVisImpl) CurrentDisChargePower() (float64, error) {
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

	return value.GetValue(), nil
}

// return the total charge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *InverterVisImpl) TotalChargeEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeCharge
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

// return the total discharge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *InverterVisImpl) TotalDischargeEnergy() (float64, error) {
	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeDischarge
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

// return the current state of charge in %
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (g *InverterVisImpl) CurrentStateOfCharge() (float64, error) {
	measurement := model.MeasurementTypeTypePercentage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeStateOfCharge
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

// helper

func (g *InverterVisImpl) getValuesForTypeCommodityScope(measurement model.MeasurementTypeType, commodity model.CommodityTypeType, scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if g.inverterEntity == nil {
		return nil, ErrDeviceDisconnected
	}

	if g.inverterMeasurement == nil {
		return nil, features.ErrDataNotAvailable
	}

	return g.inverterMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
