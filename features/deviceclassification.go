package features

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type ManufacturerDetails struct {
	BrandName                      string
	VendorName                     string
	VendorCode                     string
	DeviceName                     string
	DeviceCode                     string
	SerialNumber                   string
	SoftwareRevision               string
	HardwareRevision               string
	PowerSource                    string
	ManufacturerNodeIdentification string
	ManufacturerLabel              string
	ManufacturerDescription        string
}

// request DeviceClassificationManufacturerData from a remote device entity 1
func RequestManufacturerDetailsForDevice(service *service.EEBUSService, device *spine.DeviceRemoteImpl) (*model.MsgCounterType, error) {
	return RequestManufacturerDetailsForEntity(service, device.Entity([]model.AddressEntityType{1}))
}

// request DeviceClassificationManufacturerData from a remote device entity
func RequestManufacturerDetailsForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// get the current manufacturer details for a remote device entity
func GetManufacturerDetails(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*ManufacturerDetails, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := featureRemote.Data(model.FunctionTypeDeviceClassificationManufacturerData).(*model.DeviceClassificationManufacturerDataType)

	details := &ManufacturerDetails{}

	if data.BrandName != nil {
		details.BrandName = string(*data.BrandName)
	}
	if data.VendorName != nil {
		details.VendorName = string(*data.VendorName)
	}
	if data.VendorCode != nil {
		details.VendorCode = string(*data.VendorCode)
	}
	if data.DeviceName != nil {
		details.DeviceName = string(*data.DeviceName)
	}
	if data.DeviceCode != nil {
		details.DeviceCode = string(*data.DeviceCode)
	}
	if data.SerialNumber != nil {
		details.SerialNumber = string(*data.SerialNumber)
	}
	if data.SoftwareRevision != nil {
		details.SoftwareRevision = string(*data.SoftwareRevision)
	}
	if data.HardwareRevision != nil {
		details.HardwareRevision = string(*data.HardwareRevision)
	}
	if data.PowerSource != nil {
		details.PowerSource = string(*data.PowerSource)
	}
	if data.ManufacturerNodeIdentification != nil {
		details.ManufacturerNodeIdentification = string(*data.ManufacturerNodeIdentification)
	}
	if data.ManufacturerLabel != nil {
		details.ManufacturerLabel = string(*data.ManufacturerLabel)
	}
	if data.ManufacturerDescription != nil {
		details.ManufacturerDescription = string(*data.ManufacturerDescription)
	}

	return details, nil
}
