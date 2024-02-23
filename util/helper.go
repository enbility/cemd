package util

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC}

func IsPayloadForEntityType(payload spineapi.EventPayload, entityType model.EntityTypeType) bool {
	if payload.Entity == nil {
		return false
	}

	theEntityType := payload.Entity.EntityType()
	return theEntityType == entityType
}

func IsDeviceDisconnected(payload spineapi.EventPayload) bool {
	return (payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeRemove)
}

func IsEvseConnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeAdd {
		return IsPayloadForEntityType(payload, model.EntityTypeTypeEVSE)
	}

	return false
}

func IsEvseDisconnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeRemove {
		return IsPayloadForEntityType(payload, model.EntityTypeTypeEVSE)
	}

	return false
}

func IsEvConnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeAdd {
		return IsPayloadForEntityType(payload, model.EntityTypeTypeEV)
	}

	return false
}

func IsEvDisconnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeRemove {
		return IsPayloadForEntityType(payload, model.EntityTypeTypeEV)
	}

	return false
}
