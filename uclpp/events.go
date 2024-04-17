package uclpp

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCLPP) HandleEvent(payload spineapi.EventPayload) {
	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.connected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.LoadControlLimitDescriptionListDataType:
		e.loadControlLimitDescriptionDataUpdate(payload.Entity)
	case *model.LoadControlLimitListDataType:
		e.loadControlLimitDataUpdate(payload)
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.configurationDescriptionDataUpdate(payload.Entity)
	case *model.DeviceConfigurationKeyValueListDataType:
		e.configurationDataUpdate(payload)
	}
}

// the remote entity was connected
func (e *UCLPP) connected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if loadControl, err := util.LoadControl(e.service, entity); err == nil {
		if _, err := loadControl.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get descriptions
		if _, err := loadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if localDeviceDiag, err := util.DeviceDiagnosis(e.service, entity); err == nil {
		if _, err := localDeviceDiag.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit description data was updated
func (e *UCLPP) loadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if loadControl, err := util.LoadControl(e.service, entity); err == nil {
		// get values
		if _, err := loadControl.RequestLimitValues(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *UCLPP) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	loadControl, err := util.LoadControl(e.service, payload.Entity)
	if err != nil {
		return
	}

	data, err := loadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation)
	if err != nil {
		return
	}

	for _, item := range data {
		if item.LimitId == nil {
			continue
		}

		_, err := loadControl.GetLimitValueForLimitId(*item.LimitId)
		if err != nil {
			continue
		}

		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
		return
	}
}

// the configuration key description data was updated
func (e *UCLPP) configurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := util.DeviceConfiguration(e.service, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data was updated
func (e *UCLPP) configurationDataUpdate(payload spineapi.EventPayload) {
	if _, err := e.FailsafeProductionActivePowerLimit(payload.Entity); err != nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeProductionActivePowerLimit)
	}
	if _, err := e.FailsafeDurationMinimum(payload.Entity); err != nil {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
	}
}
