package usecases

import (
	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EVSECommissioningAndConfiguration) HandleEvent(payload spine.EventPayload) {
	// if this is not an event for any connected SKIs, ignore it
	// if !slices.Contains(e.connectedSKIs, payload.Ski) {
	// 	return
	// }

	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			e.evseConnected(payload.Ski)
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceClassificationManufacturerDataType:
				// ignore the dataset if it is not for the registered one for the SKI
				if e.remoteEntity[payload.Ski] != payload.Entity {
					return
				}

				// don't proceed with the received possibly changeset, but the full dataset
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
