package ucevsecc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVSECC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EVSE entity or device changes for this remote device

	if util.IsDeviceDisconnected(payload) {
		e.evseDisconnected(payload.Ski, payload.Entity)
		return
	}

	entityType := model.EntityTypeTypeEVSE
	if !util.IsPayloadForEntityType(payload, entityType) {
		return
	}

	if util.IsEntityTypeConnected(payload, entityType) {
		e.evseConnected(payload.Ski, payload.Entity)
		return
	} else if util.IsEntityTypeDisconnected(payload, entityType) {
		e.evseDisconnected(payload.Ski, payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceClassificationManufacturerDataType:
		e.evseManufacturerDataUpdate(payload.Ski, payload.Entity)
	case *model.DeviceDiagnosisStateDataType:
		e.evseStateUpdate(payload.Ski, payload.Entity)
	}
}

// an EVSE was connected
func (e *UCEVSECC) evseConnected(ski string, entity spineapi.EntityRemoteInterface) {
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	if evseDeviceClassification, err := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceClassification.RequestManufacturerDetails()
	}

	if evseDeviceDiagnosis, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceDiagnosis.RequestState()
	}

	e.reader.SpineEvent(ski, entity, api.UCEVSECCEventConnected)
}

// an EVSE was disconnected
func (e *UCEVSECC) evseDisconnected(ski string, entity spineapi.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, api.UCEVSECCEventDisconnected)
}

// the manufacturer Data of an EVSE was updated
func (e *UCEVSECC) evseManufacturerDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evDeviceClassification, err := util.DeviceClassification(e.service, entity)
	if err != nil {
		return
	}

	if _, err := evDeviceClassification.GetManufacturerDetails(); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCEVSECCEventManufacturerUpdate)
	}
}

// the operating State of an EVSE was updated
func (e *UCEVSECC) evseStateUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity)
	if err != nil {
		return
	}

	if _, err := evDeviceDiagnosis.GetState(); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCEVSECCEventOperationStateUpdate)
	}
}