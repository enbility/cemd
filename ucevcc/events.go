package ucevcc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVCC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if util.IsDeviceDisconnected(payload) {
		e.evDisconnected(payload.Ski, payload.Entity)
		return
	}

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evConnected(payload.Ski, payload.Entity)
		return
	} else if util.IsEntityDisconnected(payload) {
		e.evDisconnected(payload.Ski, payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.evConfigurationDescriptionDataUpdate(payload.Entity)
	case *model.DeviceConfigurationKeyValueListDataType:
		e.evConfigurationDataUpdate(payload.Ski, payload.Entity)
	case *model.DeviceClassificationManufacturerDataType:
		e.evManufacturerDataUpdate(payload.Ski, payload.Entity)
	case *model.ElectricalConnectionParameterDescriptionListDataType:
		e.evElectricalParamerDescriptionUpdate(payload.Ski, payload.Entity)
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload.Ski, payload.Entity)
	case *model.IdentificationListDataType:
		e.evIdentificationDataUpdate(payload.Ski, payload.Entity)
	}
}

// an EV was connected
func (e *UCEVCC) evConnected(ski string, entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evDeviceClassification, err := util.DeviceClassification(e.service, entity); err == nil {
		if _, err := evDeviceClassification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get manufacturer details
		if _, err := evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		if _, err := evDeviceConfiguration.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
		// get ev configuration data
		if _, err := evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity); err == nil {
		if _, err := evDeviceDiagnosis.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get device diagnosis state
		if _, err := evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := evElectricalConnection.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection parameter descriptions
		if _, err := evElectricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical permitted values descriptions
		if _, err := evElectricalConnection.RequestPermittedValueSets(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evIdentification, err := util.Identification(e.service, entity); err == nil {
		if _, err := evIdentification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get identification
		if _, err := evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCEventConnected)
}

// an EV was disconnected
func (e *UCEVCC) evDisconnected(ski string, entity spineapi.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCEventDisconnected)
}

// the configuration key description data of an EV was updated
func (e *UCEVCC) evConfigurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := evDeviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data of an EV was updated
func (e *UCEVCC) evConfigurationDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil {
		return
	}

	// Scenario 2
	if _, err := evDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard, model.DeviceConfigurationKeyValueTypeTypeString); err == nil {
		e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCCommunicationStandardConfigurationDataUpdate)
	}

	// Scenario 3
	if _, err := evDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported, model.DeviceConfigurationKeyValueTypeTypeString); err == nil {
		e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCAsymmetricChargingConfigurationDataUpdate)
	}
}

// the identification data of an EV was updated
func (e *UCEVCC) evIdentificationDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evIdentification, err := util.Identification(e.service, entity)
	if err != nil {
		return
	}

	// Scenario 4
	if values, err := evIdentification.GetValues(); err == nil {
		for _, item := range values {
			if item.IdentificationId == nil || item.IdentificationValue == nil {
				continue
			}

			e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCIdentificationDataUpdate)
			return
		}
	}
}

// the manufacturer data of an EV was updated
func (e *UCEVCC) evManufacturerDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evDeviceClassification, err := util.DeviceClassification(e.service, entity)
	if err != nil {
		return
	}

	// Scenario 5
	if _, err := evDeviceClassification.GetManufacturerDetails(); err == nil {
		e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCManufacturerDataUpdate)
	}

}

// the electrical connection parameter description data of an EV was updated
func (e *UCEVCC) evElectricalParamerDescriptionUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := evElectricalConnection.RequestPermittedValueSets(); err != nil {
			logging.Log().Error("Error getting electrical permitted values:", err)
		}
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *UCEVCC) evElectricalPermittedValuesUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evElectricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return
	}

	data, err := evElectricalConnection.GetParameterDescriptionForScopeType(model.ScopeTypeTypeACPower)
	if err != nil || data.ParameterId == nil {
		return
	}

	values, err := evElectricalConnection.GetPermittedValueSetForParameterId(*data.ParameterId)
	if err != nil || values == nil {
		return
	}

	// Scenario 6
	e.reader.SpineEvent(ski, entity.Device(), entity, api.UCEVCCChargingPowerLimitsDataUpdate)
}
