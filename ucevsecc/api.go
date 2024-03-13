package ucevsecc

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery
type ManufacturerData struct {
	DeviceName                     string `json:"deviceName,omitempty"`
	DeviceCode                     string `json:"deviceCode,omitempty"`
	SerialNumber                   string `json:"serialNumber,omitempty"`
	SoftwareRevision               string `json:"softwareRevision,omitempty"`
	HardwareRevision               string `json:"hardwareRevision,omitempty"`
	VendorName                     string `json:"vendorName,omitempty"`
	VendorCode                     string `json:"vendorCode,omitempty"`
	BrandName                      string `json:"brandName,omitempty"`
	PowerSource                    string `json:"powerSource,omitempty"`
	ManufacturerNodeIdentification string `json:"manufacturerNodeIdentification,omitempty"`
	ManufacturerLabel              string `json:"manufacturerLabel,omitempty"`
	ManufacturerDescription        string `json:"manufacturerDescription,omitempty"`
}

// interface for the EVSE Commissioning and Configuration UseCase
type UCEVSECCInterface interface {
	api.UseCaseInterface

	// the manufacturer data of an EVSE
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// returns deviceName, serialNumber, error
	ManufacturerData(entity spineapi.EntityRemoteInterface) (*ManufacturerData, error)

	// the operating state data of an EVSE
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// returns operatingState, lastErrorCode, error
	OperatingState(entity spineapi.EntityRemoteInterface) (model.DeviceDiagnosisOperatingStateType, string, error)
}
