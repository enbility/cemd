package ucmgcp

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCMGCP) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.gridConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.gridConfigurationDescriptionDataUpdate(payload.Entity)
	case *model.DeviceConfigurationKeyValueListDataType:
		e.gridConfigurationDataUpdate(payload)
	case *model.MeasurementDescriptionListDataType:
		e.gridMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.gridMeasurementDataUpdate(payload)
	}
}

// process required steps when a grid device is connected
func (e *UCMGCP) gridConnected(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		if _, err := deviceConfiguration.Subscribe(); err != nil {
			logging.Log().Error(err)
		}

		// get configuration data
		if _, err := deviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Error(err)
		}
	}

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

// the configuration key description data of an SMGW was updated
func (e *UCMGCP) gridConfigurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data of an SMGW was updated
func (e *UCMGCP) gridConfigurationDataUpdate(payload spineapi.EventPayload) {
	if _, err := e.PowerLimitationFactor(payload.Entity); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerLimitationFactor)
	}
}

// the measurement descriptiondata of an SMGW was updated
func (e *UCMGCP) gridMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := util.Measurement(e.service, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestValues(); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *UCMGCP) gridMeasurementDataUpdate(payload spineapi.EventPayload) {
	// Scenario 2
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeACPowerTotal); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePower)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeGridFeedIn); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyFeedIn)
	}

	// Scenario 4
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeGridConsumption); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateEnergyConsumed)
	}

	// Scenario 5
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeACCurrent); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentPerPhase)
	}

	// Scenario 6
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeACVoltage); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateVoltagePerPhase)
	}

	// Scenario 7
	if _, err := util.MeasurementValueForScope(e.service, payload.Entity, model.ScopeTypeTypeACFrequency); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFrequency)
	}
}
