package uccevc

import (
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// returns the current charging strategy
func (e *UCCEVC) ChargeStrategy(entity spineapi.EntityRemoteInterface) api.EVChargeStrategyType {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return api.EVChargeStrategyTypeUnknown
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return api.EVChargeStrategyTypeUnknown
	}

	// only the time series data for singledemand is relevant for detecting the charging strategy
	data, err := evTimeSeries.GetValueForType(model.TimeSeriesTypeTypeSingleDemand)
	if err != nil {
		return api.EVChargeStrategyTypeUnknown
	}

	// without time series slots, there is no known strategy
	if data.TimeSeriesSlot == nil || len(data.TimeSeriesSlot) == 0 {
		return api.EVChargeStrategyTypeUnknown
	}

	// get the value for the first slot
	firstSlot := data.TimeSeriesSlot[0]

	switch {
	case firstSlot.Duration == nil:
		// if value is > 0 and duration does not exist, the EV is direct charging
		if firstSlot.Value != nil && firstSlot.Value.GetValue() > 0 {
			return api.EVChargeStrategyTypeDirectCharging
		}

		// maxValue will show the maximum amount the battery could take
		return api.EVChargeStrategyTypeNoDemand

	case firstSlot.Duration != nil:
		if _, err := firstSlot.Duration.GetTimeDuration(); err != nil {
			// we got an invalid duration
			return api.EVChargeStrategyTypeUnknown
		}

		if firstSlot.MinValue != nil && firstSlot.MinValue.GetValue() > 0 {
			return api.EVChargeStrategyTypeMinSoC
		}

		if firstSlot.Value != nil {
			if firstSlot.Value.GetValue() > 0 {
				// there is demand and a duration
				return api.EVChargeStrategyTypeTimedCharging
			}

			return api.EVChargeStrategyTypeNoDemand
		}
	}

	return api.EVChargeStrategyTypeUnknown
}

// returns the current energy demand in Wh and the duration
func (e *UCCEVC) EnergyDemand(entity spineapi.EntityRemoteInterface) (api.Demand, error) {
	demand := api.Demand{}

	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return demand, api.ErrNoCompatibleEntity
	}

	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		return demand, eebusapi.ErrDataNotAvailable
	}

	data, err := evTimeSeries.GetValueForType(model.TimeSeriesTypeTypeSingleDemand)
	if err != nil {
		return demand, eebusapi.ErrDataNotAvailable
	}

	// we need at least a time series slot
	if data.TimeSeriesSlot == nil {
		return demand, eebusapi.ErrDataNotAvailable
	}

	// get the value for the first slot, ignore all others, which
	// in the tests so far always have min/max/value 0
	firstSlot := data.TimeSeriesSlot[0]
	if firstSlot.MinValue != nil {
		demand.MinDemand = firstSlot.MinValue.GetValue()
	}
	if firstSlot.Value != nil {
		demand.OptDemand = firstSlot.Value.GetValue()
	}
	if firstSlot.MaxValue != nil {
		demand.MaxDemand = firstSlot.MaxValue.GetValue()
	}
	if firstSlot.Duration != nil {
		if tempDuration, err := firstSlot.Duration.GetTimeDuration(); err == nil {
			demand.DurationUntilEnd = tempDuration.Seconds()
		}
	}

	// start time has to be defined either in TimePeriod or the first slot
	relStartTime := time.Duration(0)

	startTimeSet := false
	if data.TimePeriod != nil && data.TimePeriod.StartTime != nil {
		if temp, err := data.TimePeriod.StartTime.GetTimeDuration(); err == nil {
			relStartTime = temp
			startTimeSet = true
		}
	}

	if !startTimeSet {
		if firstSlot.TimePeriod != nil && firstSlot.TimePeriod.StartTime != nil {
			if temp, err := firstSlot.TimePeriod.StartTime.GetTimeDuration(); err == nil {
				relStartTime = temp
			}
		}
	}

	demand.DurationUntilStart = relStartTime.Seconds()

	return demand, nil
}
