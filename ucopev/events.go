package ucopev

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCOPEV) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	entityType := model.EntityTypeTypeEV
	if !util.IsPayloadForEntityType(payload, entityType) {
		return
	}

	if util.IsEntityTypeConnected(payload, entityType) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.LoadControlLimitDescriptionListDataType:
		e.evLoadControlLimitDescriptionDataUpdate(payload.Entity)
	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload.Ski, payload.Entity)
	}
}

// an EV was connected
func (e *UCOPEV) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evLoadControl, err := util.LoadControl(e.service, entity); err == nil {
		if _, err := evLoadControl.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement descriptions
		if _, err := evLoadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement constraints
		if _, err := evLoadControl.RequestLimitConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit description data of an EV was updated
func (e *UCOPEV) evLoadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evLoadControl, err := util.LoadControl(e.service, entity); err == nil {
		// get measurement values
		if _, err := evLoadControl.RequestLimitValues(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data of an EV was updated
func (e *UCOPEV) evLoadControlLimitDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evLoadControl, err := util.LoadControl(e.service, entity)
	if err != nil {
		return
	}

	data, err := evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeObligation)
	if err != nil {
		return
	}

	for _, item := range data {
		if item.LimitId == nil {
			continue
		}

		_, err := evLoadControl.GetLimitValueForLimitId(*item.LimitId)
		if err != nil {
			continue
		}

		e.reader.SpineEvent(ski, entity, api.UCOPEVLoadControlLimitDataUpdate)
		return
	}

}
