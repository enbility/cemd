package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EVSECommissioningAndConfiguration) HandleEvent(payload spine.EventPayload) {
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
				_, err := features.GetManufacturerDetails(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting manufacturer data:", err)
					return
				}

				// TODO: provide the current data to the CEM
			case *model.DeviceDiagnosisStateDataType:
				_, err := features.GetDeviceDiagnosisState(e.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting device diagnosis state:", err)
				}
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
