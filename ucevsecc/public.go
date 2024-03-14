package ucevsecc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *UCEVSECC) ManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	*api.ManufacturerData,
	error,
) {
	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evseDeviceClassification, err := util.DeviceClassification(e.service, entity)
	if err != nil {
		return nil, err
	}

	data, err := evseDeviceClassification.GetManufacturerDetails()
	if err != nil {
		return nil, err
	}

	ret := &api.ManufacturerData{
		DeviceName:                     util.Deref((*string)(data.DeviceName)),
		DeviceCode:                     util.Deref((*string)(data.DeviceCode)),
		SerialNumber:                   util.Deref((*string)(data.SerialNumber)),
		SoftwareRevision:               util.Deref((*string)(data.SoftwareRevision)),
		HardwareRevision:               util.Deref((*string)(data.HardwareRevision)),
		VendorName:                     util.Deref((*string)(data.VendorName)),
		VendorCode:                     util.Deref((*string)(data.VendorCode)),
		BrandName:                      util.Deref((*string)(data.BrandName)),
		PowerSource:                    util.Deref((*string)(data.PowerSource)),
		ManufacturerNodeIdentification: util.Deref((*string)(data.ManufacturerNodeIdentification)),
		ManufacturerLabel:              util.Deref((*string)(data.ManufacturerLabel)),
		ManufacturerDescription:        util.Deref((*string)(data.ManufacturerDescription)),
	}

	return ret, nil
}

// the operating state data of an EVSE
// returns operatingState, lastErrorCode, error
func (e *UCEVSECC) OperatingState(
	entity spineapi.EntityRemoteInterface,
) (
	model.DeviceDiagnosisOperatingStateType, string, error,
) {
	operatingState := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	lastErrorCode := ""

	if entity == nil || !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return operatingState, lastErrorCode, api.ErrNoCompatibleEntity
	}

	evseDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity)
	if err != nil {
		return operatingState, lastErrorCode, err
	}

	data, err := evseDeviceDiagnosis.GetState()
	if err != nil {
		return operatingState, lastErrorCode, err
	}

	if data.OperatingState != nil {
		operatingState = *data.OperatingState
	}
	if data.LastErrorCode != nil {
		lastErrorCode = string(*data.LastErrorCode)
	}

	return operatingState, lastErrorCode, nil
}
