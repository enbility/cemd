package usecases

import "github.com/DerAndereAndi/eebus-go/spine/model"

type EVSEData struct {
	OperatingState model.DeviceDiagnosisOperatingStateType
}

// // Delegate Interface for the EVSE
// type EVSEDelegate interface {
// 	// handle device state updates from the remote EVSE device
// 	HandleEVSEDeviceState(ski string, failure bool)

// 	// handle device manufacturer data updates from the remote EVSE device
// 	HandleEVSEDeviceManufacturerData(ski string, details ManufacturerDetails)
// }
