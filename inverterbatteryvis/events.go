package inverterbatteryvis

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (i *InverterBatteryVisImpl) HandleEvent(payload spine.EventPayload) {
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

// process required steps when a battery device entity is connected
func (i *InverterBatteryVisImpl) inverterConnected(ski string, entity spine.EntityRemote) {
	i.inverterEntity = entity
	localDevice := i.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f1, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	i.inverterElectricalConnection = f1

	f2, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	i.inverterMeasurement = f2

	// subscribe
	if err := i.inverterElectricalConnection.SubscribeForEntity(); err != nil {
		logging.Log().Error(err)
	}
	if err := i.inverterMeasurement.SubscribeForEntity(); err != nil {
		logging.Log().Error(err)
	}

	// get electrical connection parameter
	if err := i.inverterElectricalConnection.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	if err := i.inverterElectricalConnection.RequestParameterDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	// get measurement parameters
	if err := i.inverterMeasurement.RequestDescriptions(); err != nil {
		logging.Log().Error(err)
	}

	if err := i.inverterMeasurement.RequestConstraints(); err != nil {
		logging.Log().Error(err)
	}
}

// a battery device entity was disconnected
func (i *InverterBatteryVisImpl) inverterDisconnected() {
	i.inverterEntity = nil

	i.inverterElectricalConnection = nil
	i.inverterMeasurement = nil
}
