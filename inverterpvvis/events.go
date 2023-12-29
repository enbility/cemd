package inverterpvvis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (i *InverterPVVisImpl) HandleEvent(payload spine.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != i.ski {
		return
	}

	// we care only about events for this remote device
	if payload.Device != nil && payload.Device.Ski() != i.ski {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeRemove:
			i.inverterDisconnected()
		}

	case spine.EventTypeEntityChange:
		entityType := payload.Entity.EntityType()
		if entityType != model.EntityTypeTypeBatterySystem {
			return
		}

		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			i.inverterConnected(payload.Ski, payload.Entity)

		case spine.ElementChangeRemove:
			i.inverterDisconnected()
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType != spine.ElementChangeUpdate {
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
				logging.Log.Error("Error getting configuration key values:", err)
			}

		case *model.ElectricalConnectionParameterDescriptionListDataType:
			if i.inverterElectricalConnection == nil {
				break
			}
			if _, err := i.inverterElectricalConnection.RequestPermittedValueSets(); err != nil {
				logging.Log.Error("Error getting electrical permitted values:", err)
			}

		case *model.ElectricalConnectionDescriptionListDataType:
			if i.inverterElectricalConnection == nil {
				break
			}
			if err := i.inverterElectricalConnection.RequestDescriptions(); err != nil {
				logging.Log.Error("Error getting electrical permitted values:", err)
			}

		case *model.MeasurementDescriptionListDataType:
			if i.inverterMeasurement == nil {
				break
			}
			if _, err := i.inverterMeasurement.RequestValues(); err != nil {
				logging.Log.Error("Error getting measurement list values:", err)
			}
		}
	}
}

// process required steps when a pv device entity is connected
func (e *InverterPVVisImpl) inverterConnected(ski string, entity *spine.EntityRemoteImpl) {
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
	if err := e.inverterDeviceConfiguration.SubscribeForEntity(); err != nil {
		logging.Log.Error(err)
	}
	if err := e.inverterElectricalConnection.SubscribeForEntity(); err != nil {
		logging.Log.Error(err)
	}
	if err := e.inverterMeasurement.SubscribeForEntity(); err != nil {
		logging.Log.Error(err)
	}

	// get device configuration data
	if err := e.inverterDeviceConfiguration.RequestDescriptions(); err != nil {
		logging.Log.Error(err)
	}

	// get electrical connection parameter
	if err := e.inverterElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log.Error(err)
	}

	if err := e.inverterElectricalConnection.RequestParameterDescriptions(); err != nil {
		logging.Log.Error(err)
	}

	// get measurement parameters
	if err := e.inverterMeasurement.RequestDescriptions(); err != nil {
		logging.Log.Error(err)
	}

	if err := e.inverterMeasurement.RequestConstraints(); err != nil {
		logging.Log.Error(err)
	}
}

// a pv device entity was disconnected
func (e *InverterPVVisImpl) inverterDisconnected() {
	e.inverterMeasurement = nil

	e.inverterElectricalConnection = nil
	e.inverterMeasurement = nil
	e.inverterDeviceConfiguration = nil
}
