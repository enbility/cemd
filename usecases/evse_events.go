package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"golang.org/x/exp/slices"
)

// Internal EventHandler Interface for the CEM
func (e *EVSECommissioningAndConfiguration) HandleEvent(payload spine.EventPayload) {
	// if this is not an event for any connected SKIs, ignore it
	if !slices.Contains(e.connectedSKIs, payload.Ski) {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			e.evseConnected(payload.Ski)
		}
	case spine.EventTypeSubscriptionChange:
		switch payload.Data.(type) {
		case model.SubscriptionManagementRequestCallType:
			data := payload.Data.(model.SubscriptionManagementRequestCallType)
			if *data.ServerFeatureType == model.FeatureTypeTypeDeviceDiagnosis {
				remoteDevice := e.service.RemoteDeviceForSki(payload.Ski)
				if remoteDevice == nil {
					fmt.Println("No remote device found for SKI:", payload.Ski)
					return
				}
				switch payload.ChangeType {
				case spine.ElementChangeAdd:
					// start sending heartbeats
					senderAddr := e.entity.Device().FeatureByTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer).Address()
					destinationAddr := payload.Feature.Address()
					if senderAddr == nil || destinationAddr == nil {
						fmt.Println("No sender or destination address found for SKI:", payload.Ski)
						return
					}
					remoteDevice := e.service.RemoteDeviceForSki(payload.Ski)
					remoteDevice.StartHeartbeatSend(senderAddr, destinationAddr)
				}
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceClassificationManufacturerDataType:
				// ignore the dataset if it is not for the registered one for the SKI
				if e.remoteEntity[payload.Ski] != payload.Entity {
					return
				}

				// don't with the received possible changeset, but the full dataset
				_, err := features.GetManufacturerDetails(e.service, e.remoteEntity[payload.Ski])
				if err != nil {
					return
				}

				// TODO: provide the current data to the CEM

			case *model.DeviceDiagnosisStateDataType:
				// if e.Delegate == nil {
				// 	return
				// }

				// deviceDiagnosisStateData := payload.Data.(model.DeviceDiagnosisStateDataType)
				// failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				// e.Delegate.HandleEVSEDeviceState(payload.Ski, failure)
			}
		}
	}
}

// process required steps when an evse is connected
func (e *EVSECommissioningAndConfiguration) evseConnected(ski string) {
	remoteDevice := e.service.RemoteDeviceForSki(ski)

	_, _ = features.RequestManufacturerDetailsForDevice(e.service, remoteDevice)
	_, _ = features.RequestDiagnosisStateForDevice(e.service, remoteDevice)
}
