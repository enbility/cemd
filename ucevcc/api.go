package ucevcc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

// interface for the EV Commissioning and Configuration UseCase
type UCEVCCInterface interface {
	api.UseCaseInterface

	// return the current charge state of the EV
	CurrentChargeState(entity spineapi.EntityRemoteInterface) (api.EVChargeStateType, error)

	// Scenario 1 & 8

	// return if the EV is connected
	EVConnected(entity spineapi.EntityRemoteInterface) bool

	// Scenario 2

	// return the current communication standard type used to communicate between EVSE and EV
	CommunicationStandard(entity spineapi.EntityRemoteInterface) (string, error)

	// Scenario 3

	// return if the EV supports asymmetric charging
	AsymmetricChargingSupported(entity spineapi.EntityRemoteInterface) (bool, error)

	// Scenario 4

	// return the identifications of the currently connected EV or nil if not available
	// these can be multiple, e.g. PCID, Mac Address, RFID
	Identifications(entity spineapi.EntityRemoteInterface) ([]IdentificationItem, error)

	// Scenario 5

	// the manufacturer data of an EVSE
	// returns deviceName, serialNumber, error
	ManufacturerData(entity spineapi.EntityRemoteInterface) (string, string, error)

	// Scenario 6

	// return the min, max, default limits for each phase of the connected EV
	CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error)

	// Scenario 7

	// is the EV in sleep mode
	EVInSleepMode(entity spineapi.EntityRemoteInterface) (bool, error)
}

// EV identification
type IdentificationItem struct {
	// the identification value
	Value string

	// the type of the identification value, e.g.
	ValueType model.IdentificationTypeType
}
