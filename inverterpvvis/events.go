package inverterpvvis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Internal EventHandler Interface for the CEM
func (i *InverterPVVis) HandleEvent(payload api.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != i.ski {
		return
	}

	// we care only about events for this remote device
	if payload.Device != nil && payload.Device.Ski() != i.ski {
		return
	}

	switch payload.EventType {
	case api.EventTypeDeviceChange:
		switch payload.ChangeType {
		case api.ElementChangeRemove:
			i.inverterDisconnected()
		}

	case api.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()
		if entityType != model.EntityTypeTypeBatterySystem {
			return
		}

		switch payload.ChangeType {
		case api.ElementChangeAdd:
			i.inverterConnected(payload.Ski, payload.Entity)

		case api.ElementChangeRemove:
			i.inverterDisconnected()
		}

	case api.EventTypeDataChange:
		if payload.ChangeType != api.ElementChangeUpdate {
			return
		}

		entityType := payload.Entity.EntityType()
		if entityType != model.EntityTypeTypeBatterySystem {
			return
		}

		switch payload.Data.(type) {
		case *model.DeviceConfigurationKeyValueDescriptionListDataType:
			if i.inverterDeviceConfiguration == nil {
				break
			}

			// key value descriptions received, now get the data
			if _, err := i.inverterDeviceConfiguration.RequestKeyValues(); err != nil {
				logging.Log().Error("Error getting configuration key values:", err)
			}

		case *model.ElectricalConnectionParameterDescriptionListDataType:
			if i.inverterElectricalConnection == nil {
				break
			}
			if _, err := i.inverterElectricalConnection.RequestPermittedValueSets(); err != nil {
				logging.Log().Error("Error getting electrical permitted values:", err)
			}

		case *model.ElectricalConnectionDescriptionListDataType:
			if i.inverterElectricalConnection == nil {
				break
			}
			if err := i.inverterElectricalConnection.RequestDescriptions(); err != nil {
				logging.Log().Error("Error getting electrical permitted values:", err)
			}

		case *model.MeasurementDescriptionListDataType:
			if i.inverterMeasurement == nil {
				break
			}
			if _, err := i.inverterMeasurement.RequestValues(); err != nil {
				logging.Log().Error("Error getting measurement list values:", err)
			}
		}
	}
}

// process required steps when a pv device entity is connected
func (e *InverterPVVis) inverterConnected(ski string, entity api.EntityRemoteInterface) {
	e.inverterEntity = entity
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f1, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.inverterElectricalConnection = f1

	f2, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.inverterMeasurement = f2

	f3, err := features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.inverterDeviceConfiguration = f3

	// subscribe
	if err := e.inverterDeviceConfiguration.Subscribe(); err != nil {
		logging.Log().Error(err)
	}
	if err := e.inverterElectricalConnection.Subscribe(); err != nil {
		logging.Log().Error(err)
	}
	if err := e.inverterMeasurement.Subscribe(); err != nil {
		logging.Log().Error(err)
	}

	// get device configuration data
	if err := e.inverterDeviceConfiguration.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	// get electrical connection parameter
	if err := e.inverterElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	if err := e.inverterElectricalConnection.RequestParameterDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	// get measurement parameters
	if err := e.inverterMeasurement.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	if err := e.inverterMeasurement.RequestConstraints(); err != nil {
		logging.Log().Error(err)
	}
}

// a pv device entity was disconnected
func (e *InverterPVVis) inverterDisconnected() {
	e.inverterMeasurement = nil

	e.inverterElectricalConnection = nil
	e.inverterMeasurement = nil
	e.inverterDeviceConfiguration = nil
}
