package ucevcem

import (
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
		e.evElectricalConnectionDescriptionDataUpdate(payload)
	case *model.MeasurementDescriptionListDataType:
		e.evMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.evMeasurementDataUpdate(payload)
	}
}

// an EV was connected
func (e *UCEVCEM) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions

	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := evElectricalConnection.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection descriptions
		if _, err := evElectricalConnection.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection parameter descriptions
		if _, err := evElectricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

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
func (e *UCEVCEM) evElectricalConnectionDescriptionDataUpdate(payload spineapi.EventPayload) {
	if payload.Data == nil {
		return
	}

	data := payload.Data.(*model.ElectricalConnectionDescriptionListDataType)

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId != nil && item.AcConnectedPhases != nil {
			e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePhasesConnected)
			return
		}
	}
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
func (e *UCEVCEM) evMeasurementDataUpdate(payload spineapi.EventPayload) {
	// Scenario 1
	if util.MeasurementCheckPayloadDataForScope(e.service, payload, model.ScopeTypeTypeACCurrent) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentPerPhase)
	}

	// Scenario 2
	if util.MeasurementCheckPayloadDataForScope(e.service, payload, model.ScopeTypeTypeACPower) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerPerPhase)
	}

	// Scenario 3
	if util.MeasurementCheckPayloadDataForScope(e.service, payload, model.ScopeTypeTypeCharge) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyCharged)
	}
}
