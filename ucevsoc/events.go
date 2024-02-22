package ucevsoc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (e *UCEVSOC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !util.IsPayloadForEntityType(payload, model.EntityTypeTypeEV) {
		return
	}

	if util.IsEvConnected(payload) {
		e.evConnected(payload.Entity)
		return
	}

	switch payload.EventType {
	case spineapi.EventTypeDataChange:
		if payload.ChangeType != spineapi.ElementChangeUpdate {
			return
		}

		switch payload.Data.(type) {
		case *model.MeasurementListDataType:
			e.evMeasurementDataUpdate(payload.Ski, payload.Entity)
		case *model.ElectricalConnectionCharacteristicListDataType:
			e.evElectricalConnectionCharacteristicsDataUpdate(payload.Ski, payload.Entity)
		}
	}
}

// an EV was connected
func (e *UCEVSOC) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evMeasurement, err := util.Measurement(e.service, entity); err == nil {
		if _, err := evMeasurement.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement descriptions
		if _, err := evMeasurement.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement constraints
		if _, err := evMeasurement.RequestConstraints(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := util.ElectricalConnection(e.service, entity); err == nil {
		if _, err := evElectricalConnection.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := evElectricalConnection.RequestCharacteristics(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the measurement data of an EV was updated
func (e *UCEVSOC) evMeasurementDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	// Scenario 1
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeStateOfCharge); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCEVSOCStateOfChargeDataUpdate)
	}

	// Scenario 3
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeStateOfHealth); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCEVSOCStateOfHealthDataUpdate)
	}

	// Scenario 4
	if _, err := util.MeasurementValueForScope(e.service, entity, model.ScopeTypeTypeTravelRange); err == nil {
		e.reader.SpineEvent(ski, entity, api.UCEVSOCActualRangeDataUpdate)
	}
}

// the elecrical connection characteristic data of an EV was updated
func (e *UCEVSOC) evElectricalConnectionCharacteristicsDataUpdate(ski string, entity spineapi.EntityRemoteInterface) {
	evElectricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return
	}

	data, err := evElectricalConnection.GetCharacteristics()
	if err != nil {
		return
	}

	for _, item := range data {
		if item.CharacteristicType == nil || item.Value == nil {
			continue
		}

		if *item.CharacteristicType == model.ElectricalConnectionCharacteristicTypeTypeEnergyCapacityNominalMax {
			e.reader.SpineEvent(ski, entity, api.UCEVSOCNominalCapacityDataUpdate)
			return
		}
	}

}
