package uccevc

import (
	"errors"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// returns the constraints for the time slots
func (e *UCCEVC) TimeSlotConstraints(entity spineapi.EntityRemoteInterface) (api.TimeSlotConstraints, error) {
	result := api.TimeSlotConstraints{}

	if !e.isCompatibleEntity(entity) {
		return result, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return result, features.ErrDataNotAvailable
	}

	constraints, err := evTimeSeries.GetConstraints()
	if err != nil {
		return result, err
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.SlotCountMin)
	}
	if constraint.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.SlotCountMax)
	}
	if constraint.SlotDurationMin != nil {
		if duration, err := constraint.SlotDurationMin.GetTimeDuration(); err == nil {
			result.MinSlotDuration = duration
		}
	}
	if constraint.SlotDurationMax != nil {
		if duration, err := constraint.SlotDurationMax.GetTimeDuration(); err == nil {
			result.MaxSlotDuration = duration
		}
	}
	if constraint.SlotDurationStepSize != nil {
		if duration, err := constraint.SlotDurationStepSize.GetTimeDuration(); err == nil {
			result.SlotDurationStepSize = duration
		}
	}

	return result, nil
}

// send power limits to the EV
// if no data is provided, default power limits with the max possible value for 7 days will be sent
func (e *UCCEVC) WritePowerLimits(entity spineapi.EntityRemoteInterface, data []api.DurationSlotValue) error {
	if !e.isCompatibleEntity(entity) {
		return api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return features.ErrDataNotAvailable
	}

	if len(data) == 0 {
		data, err = e.defaultPowerLimits(entity)
		if err != nil {
			return err
		}
	}

	constraints, err := e.TimeSlotConstraints(entity)
	if err != nil {
		return err
	}

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	desc, err := evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints)
	if err != nil {
		return features.ErrDataNotAvailable
	}

	timeSeriesSlots := []model.TimeSeriesSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		timeSeriesSlot := model.TimeSeriesSlotType{
			TimeSeriesSlotId: eebusutil.Ptr(model.TimeSeriesSlotIdType(index)),
			TimePeriod: &model.TimePeriodType{
				StartTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeStart),
			},
			MaxValue: model.NewScaledNumberType(slot.Value),
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeSeriesSlot.TimePeriod.EndTime = model.NewAbsoluteOrRelativeTimeTypeFromDuration(relativeEndTime)
		}
		timeSeriesSlots = append(timeSeriesSlots, timeSeriesSlot)

		totalDuration += slot.Duration
	}

	timeSeriesData := model.TimeSeriesDataType{
		TimeSeriesId: desc.TimeSeriesId,
		TimePeriod: &model.TimePeriodType{
			StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
			EndTime:   model.NewAbsoluteOrRelativeTimeTypeFromDuration(totalDuration),
		},
		TimeSeriesSlot: timeSeriesSlots,
	}

	_, err = evTimeSeries.WriteValues([]model.TimeSeriesDataType{timeSeriesData})

	return err
}

func (e *UCCEVC) defaultPowerLimits(entity spineapi.EntityRemoteInterface) ([]api.DurationSlotValue, error) {
	// send default power limits for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default power limits")

	evElectricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		logging.Log().Error("electrical connection feature not found")
		return nil, err
	}

	paramDesc, err := evElectricalConnection.GetParameterDescriptionForScopeType(model.ScopeTypeTypeACPower)
	if err != nil {
		logging.Log().Error("Error getting parameter descriptions:", err)
		return nil, err
	}

	permitted, err := evElectricalConnection.GetPermittedValueSetForParameterId(*paramDesc.ParameterId)
	if err != nil {
		logging.Log().Error("Error getting permitted values:", err)
		return nil, err
	}

	if len(permitted.PermittedValueSet) < 1 || len(permitted.PermittedValueSet[0].Range) < 1 {
		text := "No permitted value set available"
		logging.Log().Error(text)
		return nil, errors.New(text)
	}

	data := []api.DurationSlotValue{
		{
			Duration: 7 * time.Hour * 24,
			Value:    permitted.PermittedValueSet[0].Range[0].Max.GetValue(),
		},
	}
	return data, nil
}
