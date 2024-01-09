package grid

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *GridImpl) HandleEvent(payload spine.EventPayload) {
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
			e.gridDisconnected()
		}

	case spine.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()

		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.gridConnected(payload.Ski, payload.Entity)
			}
		case spine.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.gridDisconnected()
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {

			case *model.DeviceConfigurationKeyValueDescriptionListDataType:
				// key value descriptions received, now get the data
				if _, err := e.gridDeviceConfiguration.RequestKeyValues(); err != nil {
					logging.Log().Error("Error getting configuration key values:", err)
				}

			case *model.MeasurementDescriptionListDataType:
				if _, err := e.gridMeasurement.RequestValues(); err != nil {
					logging.Log().Error("Error getting measurement list values:", err)
				}
			}

		}

	}
}

// process required steps when a grid device is connected
func (e *GridImpl) gridConnected(ski string, entity spine.EntityRemote) {
	e.gridEntity = entity
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f1, err := features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.gridDeviceConfiguration = f1

	f2, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.gridElectricalConnection = f2

	f3, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.gridMeasurement = f3

	// subscribe
	if err := e.gridDeviceConfiguration.SubscribeForEntity(); err != nil {
		logging.Log().Error(err)
		return
	}
	if err := e.gridElectricalConnection.SubscribeForEntity(); err != nil {
		logging.Log().Error(err)
		return
	}
	if err := e.gridMeasurement.SubscribeForEntity(); err != nil {
		logging.Log().Error(err)
		return
	}

	// get configuration data
	if err := e.gridDeviceConfiguration.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	// get electrical connection parameter
	if err := e.gridElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	if err := e.gridElectricalConnection.RequestParameterDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	// get measurement parameters
	if err := e.gridMeasurement.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	if err := e.gridMeasurement.RequestConstraints(); err != nil {
		logging.Log().Error(err)
		return
	}
}

// a grid device was disconnected
func (e *GridImpl) gridDisconnected() {
	e.gridEntity = nil

	e.gridDeviceConfiguration = nil
	e.gridElectricalConnection = nil
	e.gridMeasurement = nil
}
