package ucevcem

import (
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the number of ac connected phases of the EV or 0 if it is unknown
func (e *UCEVCEM) PhasesConnected(entity spineapi.EntityRemoteInterface) (uint, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	evElectricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return 0, features.ErrDataNotAvailable
	}

	data, err := evElectricalConnection.GetDescriptions()
	if err != nil {
		return 0, features.ErrDataNotAvailable
	}

	for _, item := range data {
		if item.ElectricalConnectionId != nil && item.AcConnectedPhases != nil {
			return *item.AcConnectedPhases, nil
		}
	}

	// default to 0 if the value is not available
	return 0, nil
}

// return the last current measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCEVCEM) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := util.Measurement(e.service, entity)
	evElectricalConnection, err2 := util.ElectricalConnection(e.service, entity)
	if err != nil || err2 != nil {
		return nil, err
	}

	measurement := model.MeasurementTypeTypeCurrent
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACCurrent
	data, err := evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return nil, err
	}

	var result []float64
	refetch := true
	compare := time.Now().Add(-1 * time.Minute)

	for _, phase := range util.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			elParam, err := evElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil || elParam.AcMeasuredPhases == nil || *elParam.AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			result = append(result, phaseValue)

			if item.Timestamp != nil {
				if timestamp, err := item.Timestamp.GetTime(); err == nil {
					refetch = timestamp.Before(compare)
				}
			}
		}
	}

	// if there was no timestamp provided or the time for the last value
	// is older than 1 minute, send a read request
	if refetch {
		_, _ = evMeasurement.RequestValues()
	}

	return result, nil
}

// return the last power measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCEVCEM) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := util.Measurement(e.service, entity)
	evElectricalConnection, err2 := util.ElectricalConnection(e.service, entity)
	if err != nil || err2 != nil {
		return nil, err
	}

	var data []model.MeasurementDataType

	powerAvailable := true
	measurement := model.MeasurementTypeTypePower
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeACPower
	data, err = evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil || len(data) == 0 {
		powerAvailable = false

		// If power is not provided, fall back to power calculations via currents
		measurement = model.MeasurementTypeTypeCurrent
		scope = model.ScopeTypeTypeACCurrent
		data, err = evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
		if err != nil {
			return nil, err
		}
	}

	var result []float64

	for _, phase := range util.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			elParam, err := evElectricalConnection.GetParameterDescriptionForMeasurementId(*item.MeasurementId)
			if err != nil || elParam.AcMeasuredPhases == nil || *elParam.AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			if !powerAvailable {
				phaseValue *= e.service.Configuration().Voltage()
			}

			result = append(result, phaseValue)
		}
	}

	return result, nil
}

// return the charged energy measurement in Wh of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCEVCEM) EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return 0, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := util.Measurement(e.service, entity)
	if err != nil {
		return 0, err
	}

	measurement := model.MeasurementTypeTypeEnergy
	commodity := model.CommodityTypeTypeElectricity
	scope := model.ScopeTypeTypeCharge
	data, err := evMeasurement.GetValuesForTypeCommodityScope(measurement, commodity, scope)
	if err != nil {
		return 0, err
	}

	// we assume there is only one result
	value := data[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), err
}
