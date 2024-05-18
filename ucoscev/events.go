package ucoscev

import (
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
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)
	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload)
	}
}

// the load control limit data of an EV was updated
func (e *UCOSCEV) evLoadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if util.LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false,
		e.service, payload, model.LoadControlLimitTypeTypeMaxValueLimit,
		model.LoadControlCategoryTypeRecommendation, "",
		model.ScopeTypeTypeSelfConsumption) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *UCOSCEV) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
	evElectricalConnection, err := util.ElectricalConnection(e.service, payload.Entity)
	if err != nil {
		return
	}

	data, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(model.ElectricalConnectionPhaseNameTypeA)
	if err != nil || data.ParameterId == nil {
		return
	}

	values, err := evElectricalConnection.GetPermittedValueSetForParameterId(*data.ParameterId)
	if err != nil || values == nil {
		return
	}

	// Scenario 6
	e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentLimits)
}
