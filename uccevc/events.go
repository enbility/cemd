package uccevc

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCCEVC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange {
		return
	}

	if payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.TimeSeriesDescriptionListDataType:
		e.evTimeSeriesDescriptionDataUpdate(payload.Ski, payload.Entity)

	case *model.TimeSeriesListDataType:
		e.evTimeSeriesDataUpdate(payload.Ski, payload.Entity)

	case *model.IncentiveTableDescriptionDataType:
		e.evIncentiveTableDescriptionDataUpdate(payload.Ski, payload.Entity)

	case *model.IncentiveDataType:
		e.evIncentiveTableDataUpdate(payload.Ski, payload.Entity)
	}
}

// an EV was connected
func (e *UCCEVC) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		if _, err := evDeviceConfiguration.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get device configuration descriptions
		if _, err := evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evTimeSeries, err := util.TimeSeries(e.service, entity); err == nil {
		if _, err := evTimeSeries.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := evTimeSeries.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get time series descriptions
		if _, err := evTimeSeries.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get time series constraints
		if _, err := evTimeSeries.RequestConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evIncentiveTable, err := util.IncentiveTable(e.service, entity); err == nil {
		if _, err := evIncentiveTable.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := evIncentiveTable.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get incentivetable descriptions
		if _, err := evIncentiveTable.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the time series description data of an EV was updated
func (e *UCCEVC) evTimeSeriesDescriptionDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	if evTimeSeries, err := util.TimeSeries(e.service, entity); err == nil {
		// get time series values
		if _, err := evTimeSeries.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// check if we are required to update the plan
	if !e.evCheckTimeSeriesDescriptionConstraintsUpdateRequired(entity) {
		return
	}

	_, err := e.EnergyDemand(entity)
	if err != nil {
		return
	}

	e.eventCB(ski, entity.Device(), entity, DataUpdateEnergyDemand)

	_, err = e.TimeSlotConstraints(entity)
	if err != nil {
		logging.Log().Error("Error getting timeseries constraints:", err)
		return
	}

	_, err = e.IncentiveConstraints(entity)
	if err != nil {
		logging.Log().Error("Error getting incentive constraints:", err)
		return
	}

	e.eventCB(ski, entity.Device(), entity, DataRequestedPowerLimitsAndIncentives)
}

// the load control limit data of an EV was updated
func (e *UCCEVC) evTimeSeriesDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	if _, err := e.ChargePlan(entity); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdateChargePlan)
	}

	if _, err := e.ChargePlanConstraints(entity); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdateTimeSlotConstraints)
	}
}

// the incentive table description data of an EV was updated
func (e *UCCEVC) evIncentiveTableDescriptionDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	if evIncentiveTable, err := util.IncentiveTable(e.service, entity); err == nil {
		// get time series values
		if _, err := evIncentiveTable.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// check if we are required to update the plan
	if e.evCheckIncentiveTableDescriptionUpdateRequired(entity) {
		e.eventCB(ski, entity.Device(), entity, DataRequestedIncentiveTableDescription)
	}

}

// the load control limit data of an EV was updated
func (e *UCCEVC) evIncentiveTableDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	e.eventCB(ski, entity.Device(), entity, DataUpdateIncentiveTable)
}

// check timeSeries descriptions if constraints element has updateRequired set to true
// as this triggers the CEM to send power tables within 20s
func (e *UCCEVC) evCheckTimeSeriesDescriptionConstraintsUpdateRequired(entity spineapi.EntityRemoteInterface) bool {
	evTimeSeries, err := util.TimeSeries(e.service, entity)
	if err != nil {
		logging.Log().Error("timeseries feature not found")
		return false
	}

	data, err := evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints)
	if err != nil {
		return false
	}

	if data.UpdateRequired != nil {
		return *data.UpdateRequired
	}

	return false
}

// check incentibeTable descriptions if the tariff description has updateRequired set to true
// as this triggers the CEM to send incentive tables within 20s
func (e *UCCEVC) evCheckIncentiveTableDescriptionUpdateRequired(entity spineapi.EntityRemoteInterface) bool {
	evIncentiveTable, err := util.IncentiveTable(e.service, entity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return false
	}

	data, err := evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
	if err != nil || len(data) == 0 {
		return false
	}

	// only use the first description and therein the first tariff
	item := data[0].TariffDescription
	if item != nil && item.UpdateRequired != nil {
		return *item.UpdateRequired
	}

	return false
}
