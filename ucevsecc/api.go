package ucevsecc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCEVSECCInterface interface {
	api.UseCaseInterface

	// the manufacturer data of an EVSE
	// returns deviceName, serialNumber, error
	ManufacturerData(entity spineapi.EntityRemoteInterface) (string, string, error)

	// the operating state data of an EVSE
	// returns operatingState, lastErrorCode, error
	OperatingState(entity spineapi.EntityRemoteInterface) (model.DeviceDiagnosisOperatingStateType, string, error)
}
