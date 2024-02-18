package ucevcc

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEvCC) HandleEvent(payload api.EventPayload) {
	// only about events from an EVSE entity or device changes for this remote device

	if payload.Entity == nil {
		return
	}

	entityType := payload.Entity.EntityType()
	if entityType != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case api.EventTypeDeviceChange:
		if payload.ChangeType == api.ElementChangeRemove {
			e.evDisconnected(payload.Ski, payload.Entity)
		}

	case api.EventTypeEntityChange:
		switch payload.ChangeType {
		case api.ElementChangeAdd:
			e.evConnected(payload.Ski, payload.Entity)
		case api.ElementChangeRemove:
			e.evDisconnected(payload.Ski, payload.Entity)
		}

	case api.EventTypeDataChange:
		if payload.ChangeType != api.ElementChangeUpdate {
			return
		}

		switch payload.Data.(type) {
		case *model.DeviceConfigurationKeyValueDescriptionListDataType:
			e.evConfigurationDataUpdate(payload.Ski, payload.Entity)
		case *model.DeviceClassificationManufacturerDataType:
			e.evManufacturerDataUpdate(payload.Ski, payload.Entity)
		case *model.ElectricalConnectionParameterDescriptionListDataType:
			e.evElectricalParamerDescriptionUpdate(payload.Ski, payload.Entity)
		case *model.ElectricalConnectionPermittedValueSetListDataType:
			e.evElectricalParamerDescriptionUpdate(payload.Ski, payload.Entity)
		}
	}
}

// an EVSE was connected
func (e *UCEvCC) evConnected(ski string, entity api.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, bindings
	if evDeviceClassification, err := util.DeviceClassification(e.service, entity); err == nil {
		if err := evDeviceClassification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get manufacturer details
		if _, err := evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		if err := evDeviceConfiguration.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
		// get ev configuration data
		if err := evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity); err == nil {
		if err := evDeviceDiagnosis.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get device diagnosis state
		if _, err := evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if err := evElectricalConnection.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection parameter
		if err := evElectricalConnection.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		if err := evElectricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

	}

	if evIdentification, err := util.Identification(e.service, entity); err == nil {
		if err := evIdentification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get identification
		if _, err := evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	e.reader.SpineEvent(ski, entity, UCEvCCEventConnected)
}

// an EV was disconnected
func (e *UCEvCC) evDisconnected(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvCCEventDisconnected)
}

// the configuration key Data of an EV was updated
func (e *UCEvCC) evConfigurationDataUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvCCEventConfigurationUdpate)
}

// the manufacturer Data of an EV was updated
func (e *UCEvCC) evManufacturerDataUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvCCEventManufacturerUpdate)
}

// the manufacturer Data of an EV was updated
func (e *UCEvCC) evElectricalParamerDescriptionUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvCCEventChargingPowerLimitsUpdate)
}
