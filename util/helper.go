package util

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC}

// check if the given usecase, actor is supported by the remote device
func IsUsecaseSupported(
	usecase model.UseCaseNameType,
	actor model.UseCaseActorType,
	remoteDevice spineapi.DeviceRemoteInterface) bool {
	uci := remoteDevice.UseCases()

	for _, element := range uci {
		if *element.Actor != actor {
			continue
		}
		for _, uc := range element.UseCaseSupport {
			if uc.UseCaseName != nil && *uc.UseCaseName == usecase {
				return true
			}
		}
	}

	return false
}

// return the remote entity of a given type and device ski
func EntityOfTypeForSki(
	service api.ServiceInterface,
	entityType model.EntityTypeType,
	ski string) (spineapi.EntityRemoteInterface, error) {
	rDevice := service.LocalDevice().RemoteDeviceForSki(ski)

	if rDevice == nil {
		return nil, features.ErrEntityNotFound
	}

	entities := rDevice.Entities()
	for _, entity := range entities {
		if entity.EntityType() == entityType {
			return entity, nil
		}
	}

	return nil, features.ErrEntityNotFound
}

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
