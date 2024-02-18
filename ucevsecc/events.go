package ucevsecc

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEvseCC) HandleEvent(payload api.EventPayload) {
	// only about events from an EVSE entity or device changes for this remote device

	if payload.Entity == nil {
		return
	}

	entityType := payload.Entity.EntityType()
	if entityType != model.EntityTypeTypeEVSE {
		return
	}

	switch payload.EventType {
	case api.EventTypeDeviceChange:
		if payload.ChangeType == api.ElementChangeRemove {
			e.evseDisconnected(payload.Ski, payload.Entity)
		}

	case api.EventTypeEntityChange:
		switch payload.ChangeType {
		case api.ElementChangeAdd:
			e.evseConnected(payload.Ski, payload.Entity)
		case api.ElementChangeRemove:
			e.evseDisconnected(payload.Ski, payload.Entity)
		}

	case api.EventTypeDataChange:
		if payload.ChangeType != api.ElementChangeUpdate {
			return
		}

		switch payload.Data.(type) {
		case *model.DeviceClassificationManufacturerDataType:
			e.evseManufacturerDataUpdate(payload.Ski, payload.Entity)
		case *model.DeviceDiagnosisStateDataType:
			e.evseStateUpdate(payload.Ski, payload.Entity)
		}
	}
}

// an EVSE was connected
func (e *UCEvseCC) evseConnected(ski string, entity api.EntityRemoteInterface) {
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	if evseDeviceClassification, err := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceClassification.RequestManufacturerDetails()
	}

	if evseDeviceDiagnosis, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceDiagnosis.RequestState()
	}

	e.reader.SpineEvent(ski, entity, UCEvseCCEventConnected)
}

// an EVSE was disconnected
func (e *UCEvseCC) evseDisconnected(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvseCCEventDisconnected)
}

// the manufacturer Data of an EVSE was updated
func (e *UCEvseCC) evseManufacturerDataUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvseCCEventManufacturerUpdate)
}

// the operating State of an EVSE was updated
func (e *UCEvseCC) evseStateUpdate(ski string, entity api.EntityRemoteInterface) {
	e.reader.SpineEvent(ski, entity, UCEvseCCEventOperationStateUpdate)
}
