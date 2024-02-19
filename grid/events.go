package grid

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Internal EventHandler Interface for the CEM
func (e *Grid) HandleEvent(payload api.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != e.ski {
		return
	}

	// we care only about events for this remote device
	if payload.Device != nil && payload.Device.Ski() != e.ski {
		return
	}

	switch payload.EventType {
	case api.EventTypeDeviceChange:
		switch payload.ChangeType {
		case api.ElementChangeRemove:
			e.gridDisconnected()
		}

	case api.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()

		switch payload.ChangeType {
		case api.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.gridConnected(payload.Ski, payload.Entity)
			}
		case api.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeGridConnectionPointOfPremises:
				e.gridDisconnected()
			}
		}

	case api.EventTypeDataChange:
		if payload.ChangeType == api.ElementChangeUpdate {
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
func (e *Grid) gridConnected(ski string, entity api.EntityRemoteInterface) {
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
	if _, err := e.gridDeviceConfiguration.Subscribe(); err != nil {
		logging.Log().Error(err)
		return
	}
	if _, err := e.gridElectricalConnection.Subscribe(); err != nil {
		logging.Log().Error(err)
		return
	}
	if _, err := e.gridMeasurement.Subscribe(); err != nil {
		logging.Log().Error(err)
		return
	}

	// get configuration data
	if _, err := e.gridDeviceConfiguration.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	// get electrical connection parameter
	if _, err := e.gridElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
		return
	}

	if _, err := e.gridElectricalConnection.RequestParameterDescriptions(); err != nil {
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
func (e *Grid) gridDisconnected() {
	e.gridEntity = nil

	e.gridDeviceConfiguration = nil
	e.gridElectricalConnection = nil
	e.gridMeasurement = nil
}
