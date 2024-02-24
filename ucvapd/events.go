package ucvapd

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCVAPD) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	entityType := model.EntityTypeTypePVSystem
	if !util.IsPayloadForEntityType(payload, entityType) {
		return
	}

	if util.IsEntityTypeConnected(payload, entityType) {
		e.inverterConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.inverterConfigurationDescriptionDataUpdate(payload.Entity)
	case *model.DeviceConfigurationKeyValueListDataType:
		e.inverterConfigurationDataUpdate(payload.Ski, payload.Entity)
	case *model.MeasurementDescriptionListDataType:
		e.inverterMeasurementDescriptionDataUpdate(payload.Entity)
	case *model.MeasurementListDataType:
		e.inverterMeasurementDataUpdate(payload.Ski, payload.Entity)
	}
}

// process required steps when a grid device is connected
func (e *UCVAPD) inverterConnected(entity spineapi.EntityRemoteInterface) {
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
func (e *UCVAPD) inverterConfigurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *UCVAPD) inverterConfigurationDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 1
	if deviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		if _, err := deviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem, model.DeviceConfigurationKeyValueTypeTypeScaledNumber); err == nil {
			e.reader.SpineEvent(ski, entity, api.UCVAPDPeakPowerDataUpdate)
		}
	}
}

// the measurement descriptiondata of an SMGW was updated
func (e *UCVAPD) inverterMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := util.Measurement(e.service, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestValues(); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *UCVAPD) inverterMeasurementDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 2
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACPowerTotal); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCVAPDPowerTotalMeasurementDataUpdate)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeACYieldTotal); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCVAPDYieldTotalMeasurementDataUpdate)
	}
}
