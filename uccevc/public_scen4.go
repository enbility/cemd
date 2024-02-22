package uccevc

import (
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *UCCEVC) ChargePlanConstraints(entity spineapi.EntityRemoteInterface) ([]api.DurationSlotValue, error) {
	constraints := []api.DurationSlotValue{}

	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return constraints, api.ErrEVDisconnected
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return constraints, features.ErrDataNotAvailable
	}

	data, err := evTimeSeries.GetValueForType(model.TimeSeriesTypeTypeConstraints)
	if err != nil {
		return constraints, features.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data.TimeSeriesSlot == nil {
		return constraints, features.ErrDataNotAvailable
	}

	// get the values for all slots
	for _, slot := range data.TimeSeriesSlot {
		newSlot := api.DurationSlotValue{}

		if slot.Duration != nil {
			if duration, err := slot.Duration.GetTimeDuration(); err == nil {
				newSlot.Duration = duration
			}
		} else if slot.TimePeriod != nil {
			var slotStart, slotEnd time.Time
			if slot.TimePeriod.StartTime != nil {
				if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
					slotStart = time
				}
			}
			if slot.TimePeriod.EndTime != nil {
				if time, err := slot.TimePeriod.EndTime.GetTime(); err == nil {
					slotEnd = time
				}
			}
			newSlot.Duration = slotEnd.Sub(slotStart)
		}

		if slot.MaxValue != nil {
			newSlot.Value = slot.MaxValue.GetValue()
		}

		constraints = append(constraints, newSlot)
	}

	return constraints, nil
}

func (e *UCCEVC) ChargePlan(entity spineapi.EntityRemoteInterface) (api.ChargePlan, error) {
	plan := api.ChargePlan{}

	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return plan, api.ErrEVDisconnected
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return plan, features.ErrDataNotAvailable
	}

	data, err := evTimeSeries.GetValueForType(model.TimeSeriesTypeTypePlan)
	if err != nil {
		return plan, features.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data.TimeSeriesSlot == nil {
		return plan, features.ErrDataNotAvailable
	}

	startAvailable := false
	// check the start time relative to now of the plan, default is now
	currentStart := time.Now()
	currentEnd := currentStart
	if data.TimePeriod != nil && data.TimePeriod.StartTime != nil {

		if start, err := data.TimePeriod.StartTime.GetTimeDuration(); err == nil {
			currentStart = currentStart.Add(start)
			startAvailable = true
		}
	}

	// get the values for all slots
	for index, slot := range data.TimeSeriesSlot {
		newSlot := api.ChargePlanSlotValue{}

		slotStartDefined := false
		if index == 0 && startAvailable && (slot.TimePeriod == nil || slot.TimePeriod.StartTime == nil) {
			newSlot.Start = currentStart
			slotStartDefined = true
		}
		if slot.TimePeriod != nil && slot.TimePeriod.StartTime != nil {
			if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
				newSlot.Start = time
				slotStartDefined = true
			}
		}
		if !slotStartDefined {
			newSlot.Start = currentEnd
		}

		if slot.Duration != nil {
			if duration, err := slot.Duration.GetTimeDuration(); err == nil {
				newSlot.End = newSlot.Start.Add(duration)
				currentEnd = newSlot.End
			}
		} else if slot.TimePeriod != nil && slot.TimePeriod.EndTime != nil {
			if time, err := slot.TimePeriod.StartTime.GetTime(); err == nil {
				newSlot.End = time
				currentEnd = newSlot.End
			}
		}

		if slot.Value != nil {
			newSlot.Value = slot.Value.GetValue()
		}
		if slot.MinValue != nil {
			newSlot.MinValue = slot.MinValue.GetValue()
		}
		if slot.MaxValue != nil {
			newSlot.MaxValue = slot.MaxValue.GetValue()
		}

		plan.Slots = append(plan.Slots, newSlot)
	}

	return plan, nil
}
