package emobility

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/util"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EMobilityImpl) HandleEvent(payload spine.EventPayload) {
	// we only care about events from an EVSE or EV entity
	if payload.Entity == nil {
		return
	}
	entityType := payload.Entity.EntityType()
	if entityType != model.EntityTypeTypeEVSE && entityType != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseConnected(payload.Ski)
			case model.EntityTypeTypeEV:
				e.evConnected(payload.Entity)
			}
		case spine.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeEV:
				e.evDisconnected(payload.Entity)
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceClassificationManufacturerDataType:
				_, err := util.GetManufacturerDetails(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting manufacturer data:", err)
					return
				}

				// TODO: provide the current data to the CEM
			case *model.DeviceConfigurationKeyValueDescriptionListDataType:
				// key value descriptions received, now get the data
				_, err := util.RequestDeviceConfigurationKeyValueList(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting configuration key values:", err)
				}

			case *model.DeviceConfigurationKeyValueListDataType:
				data, err := util.GetDeviceConfigurationValues(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting device configuration values:", err)
					return
				}

				// TODO: provide the device configuration data
				fmt.Printf("Device Configuration Values: %#v\n", data)

			case *model.DeviceDiagnosisStateDataType:
				_, err := util.GetDeviceDiagnosisState(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting device diagnosis state:", err)
				}

			case *model.ElectricalConnectionDescriptionListDataType:
				data, err := util.GetElectricalDescription(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting electrical description:", err)
					return
				}

				// TODO: provide the electrical description data
				fmt.Printf("Electrical Description: %#v\n", data)
			case *model.ElectricalConnectionParameterDescriptionListDataType:
				_, err := util.RequestElectricalPermittedValueSet(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting electrical permitted values:", err)
				}

			case *model.ElectricalConnectionPermittedValueSetListDataType:
				data, err := util.GetElectricalLimitValues(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting electrical limit values:", err)
					return
				}

				// TODO: provide the electrical limit data
				fmt.Printf("Electrical Permitted Values: %#v\n", data)

			case *model.LoadControlLimitDescriptionListDataType:
				_, err := util.RequestLoadControlLimitList(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
				}

			case *model.LoadControlLimitListDataType:
				data, err := util.GetLoadControlLimitValues(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
					return
				}

				// TODO: provide the loadcontrol limit data
				fmt.Printf("Loadcontrol Limits: %#v\n", data)

			case *model.MeasurementDescriptionListDataType:
				_, err := util.RequestMeasurementList(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting measurement values:", err)
				}

			case *model.MeasurementListDataType:
				data, err := util.GetMeasurementValues(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting measurement values:", err)
					return
				}

				// TODO: provide the measurement data
				fmt.Printf("Measurements: %#v\n", data)

			case *model.IdentificationListDataType:
				data, err := util.GetIdentificationValues(e.service, payload.Entity)
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
func (e *EMobilityImpl) evseConnected(ski string) {
	remoteDevice := e.service.RemoteDeviceForSki(ski)

	_, _ = util.RequestManufacturerDetailsForDevice(e.service, remoteDevice)
	_, _ = util.RequestDiagnosisStateForDevice(e.service, remoteDevice)
}

// an EV was disconnected, trigger required cleanup
func (e *EMobilityImpl) evDisconnected(entity *spine.EntityRemoteImpl) {
	fmt.Println("EV DISCONNECTED")

	// TODO: add error handling

}

// an EV was connected, trigger required communication
func (e *EMobilityImpl) evConnected(entity *spine.EntityRemoteImpl) {
	fmt.Println("EV CONNECTED")

	// TODO: add error handling

	// get ev configuration data
	if err := util.RequestDeviceConfiguration(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get manufacturer details
	if _, err := util.RequestManufacturerDetailsForEntity(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get device diagnosis state
	if _, err := util.RequestDiagnosisStateForEntity(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get electrical connection parameter
	if err := util.RequestElectricalConnectionDescription(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	if err := util.RequestElectricalConnectionParameterDescription(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get measurement parameters
	if err := util.RequestMeasurementDescription(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get identification
	if _, err := util.RequestIdentification(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

	// get loadlimit parameter
	if err := util.RequestLoadControlLimitDescription(e.service, entity); err != nil {
		fmt.Println(err)
		return
	}

}
