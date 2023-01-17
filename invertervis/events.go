package invertervis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *InverterVisImpl) HandleEvent(payload spine.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != e.ski {
		return
	}

	// we care only about events for this remote device
	if payload.Device != nil && payload.Device.Ski() != e.ski {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeRemove:
			e.inverterDisconnected()
		}

	case spine.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()

		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.inverterConnected(payload.Ski, payload.Entity)
			}
		case spine.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.inverterDisconnected()
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {

			case *model.MeasurementDescriptionListDataType:
				if _, err := e.inverterMeasurement.RequestValues(); err != nil {
					logging.Log.Error("Error getting measurement list values:", err)
				}
			}

		}

	}
}

// process required steps when a inverter device is connected
func (e *InverterVisImpl) inverterConnected(ski string, entity *spine.EntityRemoteImpl) {
	e.inverterEntity = entity
	localDevice := e.service.LocalDevice()

	f1, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	if err != nil {
		return
	}
	e.inverterElectricalConnection = f1

	f2, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	if err != nil {
		return
	}
	e.inverterMeasurement = f2

	// subscribe
	if err := e.inverterElectricalConnection.SubscribeForEntity(); err != nil {
		logging.Log.Error(err)
		return
	}
	if err := e.inverterMeasurement.SubscribeForEntity(); err != nil {
		logging.Log.Error(err)
		return
	}

	// get electrical connection parameter
	if err := e.inverterElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log.Error(err)
		return
	}

	if err := e.inverterElectricalConnection.RequestParameterDescriptions(); err != nil {
		logging.Log.Error(err)
		return
	}

	// get measurement parameters
	if err := e.inverterMeasurement.RequestDescriptions(); err != nil {
		logging.Log.Error(err)
		return
	}

	if err := e.inverterMeasurement.RequestConstraints(); err != nil {
		logging.Log.Error(err)
		return
	}
}

// a inverter device was disconnected
func (e *InverterVisImpl) inverterDisconnected() {
	e.inverterEntity = nil

	e.inverterElectricalConnection = nil
	e.inverterMeasurement = nil
}
