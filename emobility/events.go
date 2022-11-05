package emobility

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EMobilityImpl) HandleEvent(payload spine.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != e.ski {
		return
	}

	// we only care about events from an EVSE or EV entity
	if payload.Entity == nil {
		return
	}
	entityType := payload.Entity.EntityType()
	if entityType != model.EntityTypeTypeEVSE && entityType != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case spine.EventTypeEntityChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseConnected(payload.Ski, payload.Entity)
			case model.EntityTypeTypeEV:
				e.evConnected(payload.Entity)
			}
		case spine.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseDisconnected(payload.Entity)
			case model.EntityTypeTypeEV:
				e.evDisconnected(payload.Entity)
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceClassificationManufacturerDataType:
				entity, exists := e.deviceClassification[payload.Entity]
				if !exists {
					return
				}
				_, err := entity.GetManufacturerDetails()
				if err != nil {
					fmt.Println("Error getting manufacturer data:", err)
					return
				}

				// TODO: provide the current data to the CEM
			case *model.DeviceConfigurationKeyValueDescriptionListDataType:
				// key value descriptions received, now get the data
				_, err := e.evDeviceConfiguration.RequestKeyValueList()
				if err != nil {
					fmt.Println("Error getting configuration key values:", err)
				}

			case *model.DeviceConfigurationKeyValueListDataType:
				data, err := e.evDeviceConfiguration.GetValues()
				if err != nil {
					fmt.Println("Error getting device configuration values:", err)
					return
				}

				// TODO: provide the device configuration data
				fmt.Printf("Device Configuration Values: %#v\n", data)

			case *model.DeviceDiagnosisStateDataType:
				entity, exists := e.deviceDiagnosis[payload.Entity]
				if !exists {
					return
				}
				_, err := entity.GetState()
				if err != nil {
					fmt.Println("Error getting device diagnosis state:", err)
				}

			case *model.ElectricalConnectionDescriptionListDataType:
				data, err := e.evElectricalConnection.GetDescription()
				if err != nil {
					fmt.Println("Error getting electrical description:", err)
					return
				}

				// TODO: provide the electrical description data
				fmt.Printf("Electrical Description: %#v\n", data)
			case *model.ElectricalConnectionParameterDescriptionListDataType:
				_, err := e.evElectricalConnection.RequestPermittedValueSet()
				if err != nil {
					fmt.Println("Error getting electrical permitted values:", err)
				}

			case *model.ElectricalConnectionPermittedValueSetListDataType:
				data, err := e.evElectricalConnection.GetEVLimitValues()
				if err != nil {
					fmt.Println("Error getting electrical limit values:", err)
					return
				}

				// TODO: provide the electrical limit data
				fmt.Printf("Electrical Permitted Values: %#v\n", data)

			case *model.LoadControlLimitDescriptionListDataType:
				_, err := e.evLoadControl.RequestLimits()
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
				}

			case *model.LoadControlLimitListDataType:
				data, err := e.evLoadControl.GetLimitValues()
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
					return
				}

				// TODO: provide the loadcontrol limit data
				fmt.Printf("Loadcontrol Limits: %#v\n", data)

			case *model.MeasurementDescriptionListDataType:
				_, err := e.evMeasurement.Request()
				if err != nil {
					fmt.Println("Error getting measurement list values:", err)
				}

			case *model.MeasurementListDataType:
				data, err := e.evMeasurement.GetValues()
				if err != nil {
					fmt.Println("Error getting measurement values:", err)
					return
				}

				// TODO: provide the measurement data
				fmt.Printf("Measurements: %#v\n", data)

			case *model.IdentificationListDataType:
				data, err := e.evIdentification.GetValues()
				if err != nil {
					fmt.Println("Error getting identification values:", err)
					return
				}

				// TODO: provide the device configuration data
				fmt.Printf("Identification Values: %#v\n", data)
			}
		}
	}
}

// process required steps when an evse is connected
func (e *EMobilityImpl) evseConnected(ski string, entity *spine.EntityRemoteImpl) {
	e.evseEntity = entity
	localDevice := e.service.LocalDevice()

	f1, err := features.NewDeviceClassification(localDevice, entity)
	if err != nil {
		return
	}
	e.deviceClassification[entity] = f1

	f2, err := features.NewDeviceDiagnosis(localDevice, entity)
	if err != nil {
		return
	}
	e.deviceDiagnosis[entity] = f2

	_, _ = e.deviceClassification[entity].RequestManufacturerDetailsForEntity()
	_, _ = e.deviceDiagnosis[entity].RequestStateForEntity()
}

// an EV was disconnected
func (e *EMobilityImpl) evseDisconnected(entity *spine.EntityRemoteImpl) {
	e.evseEntity = nil

	delete(e.deviceClassification, entity)
	delete(e.deviceDiagnosis, entity)
}

// an EV was disconnected, trigger required cleanup
func (e *EMobilityImpl) evDisconnected(entity *spine.EntityRemoteImpl) {
	e.evEntity = nil

	delete(e.deviceClassification, entity)
	delete(e.deviceDiagnosis, entity)
	e.evDeviceConfiguration = nil
	e.evElectricalConnection = nil
	e.evMeasurement = nil
	e.evIdentification = nil
	e.evLoadControl = nil

	fmt.Println("EV DISCONNECTED")

	// TODO: add error handling

}

// an EV was connected, trigger required communication
func (e *EMobilityImpl) evConnected(entity *spine.EntityRemoteImpl) {
	e.evEntity = entity
	localDevice := e.service.LocalDevice()

	fmt.Println("EV CONNECTED")

	// TODO: add error handling

	// setup features
	e.deviceClassification[entity], _ = features.NewDeviceClassification(localDevice, entity)
	e.deviceDiagnosis[entity], _ = features.NewDeviceDiagnosis(localDevice, entity)
	e.evDeviceConfiguration, _ = features.NewDeviceConfiguration(localDevice, entity)
	e.evElectricalConnection, _ = features.NewElectricalConnection(localDevice, entity)
	e.evMeasurement, _ = features.NewMeasurement(localDevice, entity)
	e.evIdentification, _ = features.NewIdentification(localDevice, entity)
	e.evLoadControl, _ = features.NewLoadControl(localDevice, entity)

	// subscribe
	if err := e.deviceClassification[entity].SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.evDeviceConfiguration.SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.deviceDiagnosis[entity].SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.evElectricalConnection.SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.evMeasurement.SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.evIdentification.SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	if err := e.evLoadControl.SubscribeForEntity(); err != nil {
		fmt.Println(err)
		return
	}
	// if err := util.SubscribeTimeSeriesForEntity(e.service, entity); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// if err := util.SubscribeIncentiveTableForEntity(e.service, entity); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// bindings
	if err := e.evLoadControl.Bind(); err != nil {
		fmt.Println(err)
		return
	}

	// get ev configuration data
	if err := e.evDeviceConfiguration.Request(); err != nil {
		fmt.Println(err)
		return
	}

	// get manufacturer details
	if _, err := e.deviceClassification[entity].RequestManufacturerDetailsForEntity(); err != nil {
		fmt.Println(err)
		return
	}

	// get device diagnosis state
	if _, err := e.deviceDiagnosis[entity].RequestStateForEntity(); err != nil {
		fmt.Println(err)
		return
	}

	// get electrical connection parameter
	if err := e.evElectricalConnection.RequestDescription(); err != nil {
		fmt.Println(err)
		return
	}

	if err := e.evElectricalConnection.RequestParameterDescription(); err != nil {
		fmt.Println(err)
		return
	}

	// get measurement parameters
	if err := e.evMeasurement.RequestDescription(); err != nil {
		fmt.Println(err)
		return
	}

	// get identification
	if _, err := e.evIdentification.Request(); err != nil {
		fmt.Println(err)
		return
	}

	// get loadlimit parameter
	if err := e.evLoadControl.RequestLimitDescription(); err != nil {
		fmt.Println(err)
		return
	}

}
