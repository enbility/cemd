package uclpcserver

import (
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCLPCServer) HandleEvent(payload spineapi.EventPayload) {
	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	localEntity := e.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

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

// the load control limit data was updated
func (e *UCLPCServer) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if _, err := e.LoadControlLimit(); err != nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the configuration key data of an SMGW was updated
func (e *UCLPCServer) configurationDataUpdate(payload spineapi.EventPayload) {
	if _, _, err := e.FailsafeConsumptionActivePowerLimit(); err != nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeConsumptionActivePowerLimit)
	}
	if _, _, err := e.FailsafeDurationMinimum(); err != nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
	}
}
