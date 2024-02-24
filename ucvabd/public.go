package ucvabd

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the current battery (dis-)charge power (W)
//
//   - positive values charge power
//   - negative values discharge power
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVABD) CurrentChargePower(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
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

	return value.GetValue(), nil
}

// return the total charge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVABD) TotalChargeEnergy(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeCharge
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

// return the total discharge energy (Wh)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVABD) TotalDischargeEnergy(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeDischarge
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

// return the current state of charge in %
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVABD) CurrentStateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypePercentage
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeStateOfCharge
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

// helper

func (e *UCVABD) isCompatibleEntity(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil ||
		(entity.EntityType() != model.EntityTypeTypeElectricityStorageSystem) {
		return false
	}

	return true
}

func (e *UCVABD) getValuesForTypeCommodityScope(
	entity spineapi.EntityRemoteInterface,
	measurement model.MeasurementTypeType,
	commodity model.CommodityTypeType,
	scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	if entity == nil {
		return nil, util.ErrDeviceDisconnected
	}

	measurementF, err := util.Measurement(e.service, entity)
	if err != nil {
		return nil, features.ErrFunctionNotSupported
	}

	return measurementF.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
