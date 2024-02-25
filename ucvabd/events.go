package ucvabd

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCVABD) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.inverterConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.inverterMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.inverterMeasurementDataUpdate(payload.Ski, payload.Entity)
	}
}

// process required steps when a grid device is connected
func (e *UCVABD) inverterConnected(entity spineapi.EntityRemoteInterface) {
	if electricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := electricalConnection.Subscribe(); err != nil {
			logging.Log().Error(err)
		}

		// get electrical connection parameter
		if _, err := electricalConnection.RequestDescriptions(); err != nil {
			logging.Log().Error(err)
		}

		if _, err := electricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Error(err)
		}
	}

	if measurement, err := util.Measurement(e.service, entity); err == nil {
		if _, err := measurement.Subscribe(); err != nil {
			logging.Log().Error(err)
		}

		// get measurement parameters
		if _, err := measurement.RequestDescriptions(); err != nil {
			logging.Log().Error(err)
		}

		if _, err := measurement.RequestConstraints(); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the measurement descriptiondata of an SMGW was updated
func (e *UCVABD) inverterMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := util.Measurement(e.service, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestValues(); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *UCVABD) inverterMeasurementDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 1
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACPowerTotal); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdatePower)
	}

	// Scenario 2
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeCharge); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdateEnergyCharged)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeDischarge); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdateEnergyDischarged)
	}

	// Scenario 4
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeStateOfCharge); err == nil {
		e.eventCB(ski, entity.Device(), entity, DataUpdateStateOfCharge)
	}
}
