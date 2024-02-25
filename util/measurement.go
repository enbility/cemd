package util

import (
	"slices"

	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the measurement value for a scope or an error
func MeasurementValueForScope(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	scope model.ScopeTypeType) (float64, error) {
	measurementF, err := Measurement(service, entity)
	if err != nil {
		return 0, features.ErrFunctionNotSupported
	}

	if data, err := measurementF.GetDescriptionsForScope(scope); err == nil {
		for _, item := range data {
			if item.MeasurementId == nil {
				continue
			}

			if value, err := measurementF.GetValueForMeasurementId(*item.MeasurementId); err == nil {
				return value, nil
			}
		}
	}

	return 0, features.ErrDataNotAvailable
}

// return the phase specific voltage details
func MeasurementValuesForTypeCommodityScope(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	measurementType model.MeasurementTypeType,
	commodityType model.CommodityTypeType,
	scopeType model.ScopeTypeType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
) ([]float64, error) {

	measurement := measurementType
	commodity := commodityType
	scope := scopeType
	data, err := GetValuesForTypeCommodityScope(service, entity, measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	electricalConnection, err := ElectricalConnection(service, entity)
	if err != nil || electricalConnection == nil {
		return nil, err
	}

	var result []float64

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		if validPhaseNameTypes != nil {
			param, err := electricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil ||
				param.AcMeasuredPhases == nil ||
				!slices.Contains(validPhaseNameTypes, *param.AcMeasuredPhases) {
				continue
			}
		}

		if energyDirection != "" {
			desc, err := electricalConnection.GetDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil {
				continue
			}

			// if energy direction is not consume
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != energyDirection {
				return nil, err
			}
		}

		value := item.Value.GetValue()

		result = append(result, value)
	}

	return result, nil
}

func GetValuesForTypeCommodityScope(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	measurement model.MeasurementTypeType,
	commodity model.CommodityTypeType,
	scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {

	measurementFeature, err := Measurement(service, entity)
	if err != nil || measurementFeature == nil {
		return nil, err
	}

	return measurementFeature.GetValuesForTypeCommodityScope(measurement, commodity, scope)
}
