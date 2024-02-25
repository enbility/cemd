package ucevcem

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVCEM) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}
	switch payload.Data.(type) {
	case *model.ElectricalConnectionDescriptionListDataType:
		e.evElectricalConnectionDescriptionDataUpdate(payload.Ski, payload.Entity)
	case *model.MeasurementDescriptionListDataType:
		e.evMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.evMeasurementDataUpdate(payload.Ski, payload.Entity)
	}
}

// an EV was connected
func (e *UCEVCEM) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evMeasurement, err := util.Measurement(e.service, entity); err == nil {
		if _, err := evMeasurement.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement descriptions
		if _, err := evMeasurement.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement constraints
		if _, err := evMeasurement.RequestConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the electrical connection description data of an EV was updated
func (e *UCEVCEM) evElectricalConnectionDescriptionDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	if _, err := e.PhasesConnected(entity); err != nil {
		return
	}

	e.eventCB(ski, entity.Device(), entity, api.UCEVCEMNumberOfConnectedPhasesDataUpdate)
}

// the measurement description data of an EV was updated
func (e *UCEVCEM) evMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evMeasurement, err := util.Measurement(e.service, entity); err == nil {
		// get measurement values
		if _, err := evMeasurement.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the measurement data of an EV was updated
func (e *UCEVCEM) evMeasurementDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 1
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACCurrent); err == nil {
		e.eventCB(ski, entity.Device(), entity, api.UCEVCEMCurrentMeasurementDataUpdate)
	}

	// Scenario 2
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACPower); err == nil {
		e.eventCB(ski, entity.Device(), entity, api.UCEVCEMPowerMeasurementDataUpdate)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeCharge); err == nil {
		e.eventCB(ski, entity.Device(), entity, api.UCEVCEMChargingEnergyMeasurementDataUpdate)
	}
}
