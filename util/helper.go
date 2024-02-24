package util

import (
	"slices"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC}

func IsCompatibleEntity(entity spineapi.EntityRemoteInterface, entityTypes []model.EntityTypeType) bool {
	if entity == nil {
		return false
	}

	return slices.Contains(entityTypes, entity.EntityType())
}

func IsDeviceDisconnected(payload spineapi.EventPayload) bool {
	return (payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeRemove)
}

func IsEntityConnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeAdd {
		return true
	}

	return false
}

func IsEntityDisconnected(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeEntityChange &&
		payload.ChangeType == spineapi.ElementChangeRemove {
		return true
	}

	return false
}
