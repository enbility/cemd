package cem

import "github.com/DerAndereAndi/eebus-go/spine"

type EVCommunicationStandardType string

const (
	EVCommunicationStandardTypeUnknown      EVCommunicationStandardType = "unknown"
	EVCommunicationStandardTypeISO151182ED1 EVCommunicationStandardType = "iso15118-2ed1"
	EVCommunicationStandardTypeISO151182ED2 EVCommunicationStandardType = "iso15118-2ed2"
	EVCommunicationStandardTypeIEC61851     EVCommunicationStandardType = "iec61851"
)

type EVIdentificationType string

const (
	EVIdentificationTypeEUI48 EVIdentificationType = "eui48" // eui48 MAC address
	EVIdentificationTypeEUI64 EVIdentificationType = "eui64" // eui64 MAC address
)

type EVData struct {
	CommunicationStandard       EVCommunicationStandardType
	AsymmetricChargingSupported bool
	IdentificationType          EVIdentificationType
	Identification              string
	ManufacturerDetails         ManufacturerDetails
}

// get the remote device specific data element
func (e *EV) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *EVData {
	if evdata, ok := e.data[remoteDevice.Ski()]; ok {
		return evdata
	}

	evData := &EVData{
		CommunicationStandard:       EVCommunicationStandardTypeIEC61851,
		AsymmetricChargingSupported: false,
	}
	e.data[remoteDevice.Ski()] = evData

	return evData
}
