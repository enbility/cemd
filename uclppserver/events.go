package uclppserver

import (
	"slices"

	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCLPPServer) HandleEvent(payload spineapi.EventPayload) {
	if util.IsDeviceConnected(payload) {
		e.deviceConnected(payload)
		return
	}

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// did we receive a binding to the loadControl server and the
	// heartbeatWorkaround is required?
	if payload.EventType == spineapi.EventTypeBindingChange &&
		payload.ChangeType == spineapi.ElementChangeAdd &&
		payload.LocalFeature != nil &&
		payload.LocalFeature.Type() == model.FeatureTypeTypeLoadControl &&
		payload.LocalFeature.Role() == model.RoleTypeServer {
		e.subscribeHeartbeatWorkaround(payload)
		return
	}

	if util.IsHeartbeat(localEntity, payload) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateHeartbeat)
		return
	}

	if localEntity == nil ||
		payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate ||
		payload.CmdClassifier == nil ||
		*payload.CmdClassifier != model.CmdClassifierTypeWrite {
		return
	}

	// the codefactor warning is invalid, as .(type) check can not be replaced with if then
	//revive:disable-next-line
	switch payload.Data.(type) {
	case *model.LoadControlLimitListDataType:
		serverF := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

		if payload.Function != model.FunctionTypeLoadControlLimitListData ||
			payload.LocalFeature != serverF {
			return
		}

		e.loadControlLimitDataUpdate(payload)
	case *model.DeviceConfigurationKeyValueListDataType:
		serverF := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

		if payload.Function != model.FunctionTypeDeviceConfigurationKeyValueListData ||
			payload.LocalFeature != serverF {
			return
		}

		e.configurationDataUpdate(payload)
	}
}

// a remote device was connected and we know its entities
func (e *UCLPPServer) deviceConnected(payload spineapi.EventPayload) {
	if payload.Device == nil {
		return
	}

	// check if there is a DeviceDiagnosis server on one or more entities
	remoteDevice := payload.Device

	var deviceDiagEntites []spineapi.EntityRemoteInterface

	entites := remoteDevice.Entities()
	for _, entity := range entites {
		if !slices.Contains(e.validEntityTypes, entity.EntityType()) {
			continue
		}

		deviceDiagF := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
		if deviceDiagF == nil {
			continue
		}

		deviceDiagEntites = append(deviceDiagEntites, entity)
	}

	// the remote device does not have a DeviceDiagnosis Server, which it should
	if len(deviceDiagEntites) == 0 {
		return
	}

	// we only found one matching entity, as it should be, subscribe
	if len(deviceDiagEntites) == 1 {
		if localDeviceDiag, err := util.DeviceDiagnosis(e.service, deviceDiagEntites[0]); err == nil {
			e.heartbeatDiag = localDeviceDiag
			if _, err := localDeviceDiag.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}

			if _, err := localDeviceDiag.RequestHeartbeat(); err != nil {
				logging.Log().Debug(err)
			}
		}

		return
	}

	// we found more than one matching entity, this is not good
	// according to KEO the subscription should be done on the entity that requests a binding to
	// the local loadControlLimit server feature
	e.heartbeatKeoWorkaround = true
}

// subscribe to the DeviceDiagnosis Server of the entity that created a binding
func (e *UCLPPServer) subscribeHeartbeatWorkaround(payload spineapi.EventPayload) {
	// the workaround is not needed, exit
	if !e.heartbeatKeoWorkaround {
		return
	}

	if localDeviceDiag, err := util.DeviceDiagnosis(e.service, payload.Entity); err == nil {
		e.heartbeatDiag = localDeviceDiag
		if _, err := localDeviceDiag.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := localDeviceDiag.RequestHeartbeat(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *UCLPPServer) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if util.LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(
		true, e.service, payload,
		model.LoadControlLimitTypeTypeSignDependentAbsValueLimit,
		model.LoadControlCategoryTypeObligation,
		model.EnergyDirectionTypeProduce,
		model.ScopeTypeTypeActivePowerLimit) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the configuration key data of an SMGW was updated
func (e *UCLPPServer) configurationDataUpdate(payload spineapi.EventPayload) {
	if util.DeviceConfigurationCheckDataPayloadForKeyName(true, e.service, payload, model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeProductionActivePowerLimit)
	}
	if util.DeviceConfigurationCheckDataPayloadForKeyName(true, e.service, payload, model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
	}
}
