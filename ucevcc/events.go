package ucevcc

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVCC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evConnected(payload)
		return
	} else if util.IsEntityDisconnected(payload) {
		e.evDisconnected(payload)
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
		e.evConfigurationDataUpdate(payload)
	case *model.DeviceDiagnosisOperatingStateType:
		e.evOperatingStateDataUpdate(payload)
	case *model.DeviceClassificationManufacturerDataType:
		e.evManufacturerDataUpdate(payload)
	case *model.ElectricalConnectionParameterDescriptionListDataType:
		e.evElectricalParamerDescriptionUpdate(payload.Entity)
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)
	case *model.IdentificationListDataType:
		e.evIdentificationDataUpdate(payload)
	}
}

// an EV was connected
func (e *UCEVCC) evConnected(payload spineapi.EventPayload) {
	// initialise features, e.g. subscriptions, descriptions
	if evDeviceClassification, err := util.DeviceClassification(e.service, payload.Entity); err == nil {
		if _, err := evDeviceClassification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get manufacturer details
		if _, err := evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceConfiguration, err := util.DeviceConfiguration(e.service, payload.Entity); err == nil {
		if _, err := evDeviceConfiguration.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
		// get ev configuration data
		if _, err := evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, payload.Entity); err == nil {
		if _, err := evDeviceDiagnosis.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get device diagnosis state
		if _, err := evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := util.ElectricalConnection(e.service, payload.Entity); err == nil {
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

	if evIdentification, err := util.Identification(e.service, payload.Entity); err == nil {
		if _, err := evIdentification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get identification
		if _, err := evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	e.eventCB(payload.Ski, payload.Device, payload.Entity, EvConnected)
}

// an EV was disconnected
func (e *UCEVCC) evDisconnected(payload spineapi.EventPayload) {
	e.eventCB(payload.Ski, payload.Device, payload.Entity, EvDisconnected)
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
func (e *UCEVCC) evConfigurationDataUpdate(payload spineapi.EventPayload) {
	evDeviceConfiguration, err := util.DeviceConfiguration(e.service, payload.Entity)
	if err != nil {
		return
	}

	// Scenario 2
	if _, err := evDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeCommunicationsStandard, model.DeviceConfigurationKeyValueTypeTypeString); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCommunicationStandard)
	}

	// Scenario 3
	if _, err := evDeviceConfiguration.GetKeyValueForKeyName(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported, model.DeviceConfigurationKeyValueTypeTypeString); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateAsymmetricChargingSupport)
	}
}

// the operating state of an EV was updated
func (e *UCEVCC) evOperatingStateDataUpdate(payload spineapi.EventPayload) {
	deviceDiagnosis, err := util.DeviceDiagnosis(e.service, payload.Entity)
	if err != nil {
		return
	}

	if _, err := deviceDiagnosis.GetState(); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIdentifications)
	}
}

// the identification data of an EV was updated
func (e *UCEVCC) evIdentificationDataUpdate(payload spineapi.EventPayload) {
	evIdentification, err := util.Identification(e.service, payload.Entity)
	if err != nil {
		return
	}

	// Scenario 4
	if values, err := evIdentification.GetValues(); err == nil {
		for _, item := range values {
			if item.IdentificationId == nil || item.IdentificationValue == nil {
				continue
			}

			e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIdentifications)
			return
		}
	}
}

// the manufacturer data of an EV was updated
func (e *UCEVCC) evManufacturerDataUpdate(payload spineapi.EventPayload) {
	evDeviceClassification, err := util.DeviceClassification(e.service, payload.Entity)
	if err != nil {
		return
	}

	// Scenario 5
	if _, err := evDeviceClassification.GetManufacturerDetails(); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateManufacturerData)
	}
}

// the electrical connection parameter description data of an EV was updated
func (e *UCEVCC) evElectricalParamerDescriptionUpdate(entity spineapi.EntityRemoteInterface) {
	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := evElectricalConnection.RequestPermittedValueSets(); err != nil {
			logging.Log().Error("Error getting electrical permitted values:", err)
		}
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *UCEVCC) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
	evElectricalConnection, err := util.ElectricalConnection(e.service, payload.Entity)
	if err != nil {
		return
	}

	data, err := evElectricalConnection.GetParameterDescriptionForScopeType(model.ScopeTypeTypeACPowerTotal)
	if err != nil || data.ParameterId == nil {
		return
	}

	values, err := evElectricalConnection.GetPermittedValueSetForParameterId(*data.ParameterId)
	if err != nil || values == nil {
		return
	}

	// Scenario 6
	e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentLimits)
}
