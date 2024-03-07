package ucevsecc

import (
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVSECC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EVSE entity or device changes for this remote device

	if util.IsDeviceDisconnected(payload) {
		e.evseDisconnected(payload)
		return
	}

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evseConnected(payload)
		return
	} else if util.IsEntityDisconnected(payload) {
		e.evseDisconnected(payload)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceClassificationManufacturerDataType:
		e.evseManufacturerDataUpdate(payload)
	case *model.DeviceDiagnosisStateDataType:
		e.evseStateUpdate(payload)
	}
}

// an EVSE was connected
func (e *UCEVSECC) evseConnected(payload spineapi.EventPayload) {
	if evseDeviceClassification, err := util.DeviceClassification(e.service, payload.Entity); err == nil {
		_, _ = evseDeviceClassification.RequestManufacturerDetails()
	}

	if evseDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, payload.Entity); err == nil {
		_, _ = evseDeviceDiagnosis.RequestState()
	}

	e.eventCB(payload.Ski, payload.Device, payload.Entity, EvseConnected)
}

// an EVSE was disconnected
func (e *UCEVSECC) evseDisconnected(payload spineapi.EventPayload) {
	e.eventCB(payload.Ski, payload.Device, payload.Entity, EvseDisconnected)
}

// the manufacturer Data of an EVSE was updated
func (e *UCEVSECC) evseManufacturerDataUpdate(payload spineapi.EventPayload) {
	evDeviceClassification, err := util.DeviceClassification(e.service, payload.Entity)
	if err != nil {
		return
	}

	if _, err := evDeviceClassification.GetManufacturerDetails(); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateManufacturerData)
	}
}

// the operating State of an EVSE was updated
func (e *UCEVSECC) evseStateUpdate(payload spineapi.EventPayload) {
	evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, payload.Entity)
	if err != nil {
		return
	}

	if _, err := evDeviceDiagnosis.GetState(); err == nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateOperatingState)
	}
}
