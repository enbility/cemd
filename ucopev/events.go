package ucopev

import (
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCOPEV) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !util.IsCompatibleEntity(payload.Entity, e.validEntityTypes) {
		return
	}

	if util.IsEntityConnected(payload) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)
	case *model.LoadControlLimitDescriptionListDataType:
		e.evLoadControlLimitDescriptionDataUpdate(payload.Entity)
	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload)
	}
}

// an EV was connected
func (e *UCOPEV) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evLoadControl, err := util.LoadControl(e.service, entity); err == nil {
		if _, err := evLoadControl.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := evLoadControl.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get descriptions
		if _, err := evLoadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get constraints
		if _, err := evLoadControl.RequestLimitConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit description data of an EV was updated
func (e *UCOPEV) evLoadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evLoadControl, err := util.LoadControl(e.service, entity); err == nil {
		// get values
		if _, err := evLoadControl.RequestLimitValues(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data of an EV was updated
func (e *UCOPEV) evLoadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if util.LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false,
		e.service, payload, model.LoadControlLimitTypeTypeMaxValueLimit,
		model.LoadControlCategoryTypeObligation, "", model.ScopeTypeTypeOverloadProtection) {
		e.eventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *UCOPEV) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
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
