package util

import (
	"slices"

	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Check the payload data if it contains measurementId values for a given scope
func MeasurementCheckPayloadDataForScope(service eebusapi.ServiceInterface, payload spineapi.EventPayload, scope model.ScopeTypeType) bool {
	measurementF, err := Measurement(service, payload.Entity)
	if err != nil || payload.Data == nil {
		return false
	}

	if data, err := measurementF.GetDescriptionsForScope(scope); err == nil {
		measurements := payload.Data.(*model.MeasurementListDataType)

		for _, item := range data {
			if item.MeasurementId == nil {
				continue
			}

			for _, measurement := range measurements.MeasurementData {
				if measurement.MeasurementId != nil &&
					*measurement.MeasurementId == *item.MeasurementId &&
					measurement.Value != nil {
					return true
				}
			}
		}
	}

	return false
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
