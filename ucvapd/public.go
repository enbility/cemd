package ucvapd

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the current photovoltaic production power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVAPD) CurrentProductionPower(entity spineapi.EntityRemoteInterface) (float64, error) {
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

// return the nominal photovoltaic peak power (W)
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCVAPD) NominalPeakPower(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	deviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil {
		return 0, features.ErrFunctionNotSupported
	}

	keyName := model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem
	if _, err = deviceConfiguration.GetDescriptionForKeyName(keyName); err != nil {
		return 0, err
	}

	data, err := deviceConfiguration.GetKeyValueForKeyName(keyName, model.DeviceConfigurationKeyValueTypeTypeScaledNumber)
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
func (e *UCVAPD) TotalPVYield(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.isCompatibleEntity(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACYieldTotal
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

func (e *UCVAPD) isCompatibleEntity(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil ||
		(entity.EntityType() != model.EntityTypeTypePVSystem) {
		return false
	}

	return true
}

func (e *UCVAPD) getValuesForTypeCommodityScope(
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
