package ucevsecc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCEvseCCInterface interface {
	api.UseCaseInterface

	// the manufacturer data of an EVSE
	// returns deviceName, serialNumber, error
	EVSEManufacturerData(ski string, entity spineapi.EntityRemoteInterface) (string, string, error)

	// the operating state data of an EVSE
	// returns operatingState, lastErrorCode, error
	EVSEOperatingState(ski string, entity spineapi.EntityRemoteInterface) (model.DeviceDiagnosisOperatingStateType, string, error)
}

const (
	// An EVSE was connected
	UCEvseCCEventConnected api.UseCaseEventType = "ucEvseConnected"

	// An EVSE was disconnected
	UCEvseCCEventDisconnected api.UseCaseEventType = "ucEvseDisonnected"

	// EVSE manufacturer data was updated
	UCEvseCCEventManufacturerUpdate api.UseCaseEventType = "ucEvseManufacturerUpdate"

	// EVSE operation state was updated
	UCEvseCCEventOperationStateUpdate api.UseCaseEventType = "ucEvseOperationStateUpdate"
)
