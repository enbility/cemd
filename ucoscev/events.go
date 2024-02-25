package ucoscev

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCOSCEV) HandleEvent(payload spineapi.EventPayload) {
	// most of the events are identical to OPEV, and OPEV is required to be used,
	// we don't handle the same events in here

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	// the codefactor warning is invalid, as .(type) check can not be replaced with if then
	//revive:disable-next-line
	switch payload.Data.(type) {
	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload.Ski, payload.Entity)
	}
}

// the load control limit data of an EV was updated
func (e *UCOSCEV) evLoadControlLimitDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evLoadControl, err := util.LoadControl(e.service, entity)
	if err != nil {
		return
	}

	data, err := evLoadControl.GetLimitDescriptionsForCategory(model.LoadControlCategoryTypeRecommendation)
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

		e.reader.Event(ski, entity.Device(), entity, api.UCOPEVLoadControlLimitDataUpdate)
		return
	}

}
