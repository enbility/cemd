package ucevsecc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *UCEvseCC) EVSEManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	string,
	string,
	error,
) {
	deviceName := ""
	serialNumber := ""

	if entity == nil || entity.EntityType() != model.EntityTypeTypeEVSE {
		return deviceName, serialNumber, api.ErrNoEvseEntity
	}

	evseDeviceClassification, err := util.DeviceClassification(e.service, entity)
	if err != nil {
		return deviceName, serialNumber, err
	}

	data, err := evseDeviceClassification.GetManufacturerDetails()
	if err != nil {
		return deviceName, serialNumber, err
	}

	if data.DeviceName != nil {
		deviceName = string(*data.DeviceName)
	}

	if data.SerialNumber != nil {
		serialNumber = string(*data.SerialNumber)
	}

	return deviceName, serialNumber, nil
}

// the operating state data of an EVSE
// returns operatingState, lastErrorCode, error
func (e *UCEvseCC) EVSEOperatingState(
	entity spineapi.EntityRemoteInterface,
) (
	model.DeviceDiagnosisOperatingStateType, string, error,
) {
	operatingState := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	lastErrorCode := ""

	if entity == nil || entity.EntityType() != model.EntityTypeTypeEVSE {
		return operatingState, lastErrorCode, api.ErrNoEvseEntity
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
