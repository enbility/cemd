package ucmpc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCMPC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.deviceConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.MeasurementDescriptionListDataType:
		e.deviceMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.deviceMeasurementDataUpdate(payload.Ski, payload.Entity)
	}
}

// process required steps when a device is connected
func (e *UCMPC) deviceConnected(entity spineapi.EntityRemoteInterface) {
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

// the measurement descriptiondata of a device was updated
func (e *UCMPC) deviceMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := util.Measurement(e.service, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestValues(); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of a device was updated
func (e *UCMPC) deviceMeasurementDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 1
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACPowerTotal); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCPowerTotalMeasurementDataUpdate)
	}

	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACPower); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCPowerPerPhaseMeasurementDataUpdate)
	}

	// Scenario 2
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACEnergyConsumed); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCEnergyConsumedMeasurementDataUpdate)
	}

	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACEnergyProduced); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCEnergyProcudedMeasurementDataUpdate)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACCurrent); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCCurrentsMeasurementDataUpdate)
	}

	// Scenario 4
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACVoltage); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCVoltagesMeasurementDataUpdate)
	}

	// Scenario 5
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACFrequency); err == nil {
		e.reader.Event(ski, entity.Device(), entity, api.UCMPCFrequencyMeasurementDataUpdate)
	}

}
