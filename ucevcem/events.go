package ucevcem

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEvCEM) HandleEvent(payload api.EventPayload) {
	// only about events from an EVSE entity or device changes for this remote device

	if payload.Entity == nil {
		return
	}

	entityType := payload.Entity.EntityType()
	if entityType != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case api.EventTypeEntityChange:
		switch payload.ChangeType {
		case api.ElementChangeAdd:
			e.evConnected(payload.Ski, payload.Entity)
		}

	case api.EventTypeDataChange:
		if payload.ChangeType != api.ElementChangeUpdate {
			return
		}

		switch payload.Data.(type) {
		case *model.MeasurementDescriptionListDataType:
			e.evMeasurementDescriptionDataUpdate(payload.Ski, payload.Entity)
		case *model.MeasurementListDataType:
			e.evMeasurementDataUpdate(payload.Ski, payload.Entity)
		}
	}
}

// an EV was connected
func (e *UCEvCEM) evConnected(ski string, entity api.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evMeasurement, err := util.Measurement(e.service, entity); err == nil {
		if _, err := evMeasurement.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement descriptions
		if err := evMeasurement.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement constraints
		if err := evMeasurement.RequestConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the measurement description data of an EV was updated
func (e *UCEvCEM) evMeasurementDescriptionDataUpdate(ski string, entity api.EntityRemoteInterface) {
	if evMeasurement, err := util.Measurement(e.service, entity); err == nil {
		// get measurement values
		if _, err := evMeasurement.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the measurement data of an EV was updated
func (e *UCEvCEM) evMeasurementDataUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvCEMMeasurementDataUpdate)
}
