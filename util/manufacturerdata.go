package util

import (
	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the current manufacturer data for a entity
//
// possible errors:
//   - ErrNoCompatibleEntity if entity is not compatible
//   - and others
func ManufacturerData(service eebusapi.ServiceInterface, entity spineapi.EntityRemoteInterface, entityTypes []model.EntityTypeType) (*api.ManufacturerData, error) {
	if entity == nil || !IsCompatibleEntity(entity, entityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	deviceClassification, err := DeviceClassification(service, entity)
	if err != nil {
		return nil, err
	}

	data, err := deviceClassification.GetManufacturerDetails()
	if err != nil {
		return nil, err
	}

	ret := &api.ManufacturerData{
		DeviceName:                     Deref((*string)(data.DeviceName)),
		DeviceCode:                     Deref((*string)(data.DeviceCode)),
		SerialNumber:                   Deref((*string)(data.SerialNumber)),
		SoftwareRevision:               Deref((*string)(data.SoftwareRevision)),
		HardwareRevision:               Deref((*string)(data.HardwareRevision)),
		VendorName:                     Deref((*string)(data.VendorName)),
		VendorCode:                     Deref((*string)(data.VendorCode)),
		BrandName:                      Deref((*string)(data.BrandName)),
		PowerSource:                    Deref((*string)(data.PowerSource)),
		ManufacturerNodeIdentification: Deref((*string)(data.ManufacturerNodeIdentification)),
		ManufacturerLabel:              Deref((*string)(data.ManufacturerLabel)),
		ManufacturerDescription:        Deref((*string)(data.ManufacturerDescription)),
	}

	return ret, nil
}
