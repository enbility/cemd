package util

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceDiagnosisType struct {
	OperatingState       model.DeviceDiagnosisOperatingStateType
	PowerSupplyCondition model.PowerSupplyConditionType
}

// request DeviceDiagnosisStateData from a remote device entity 1
func RequestDiagnosisStateForDevice(service *service.EEBUSService, device *spine.DeviceRemoteImpl) (*model.MsgCounterType, error) {
	return RequestDiagnosisStateForEntity(service, device.Entity([]model.AddressEntityType{1}))
}

// request DeviceDiagnosisStateData from a remote entity
func RequestDiagnosisStateForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// request FunctionTypeDeviceDiagnosisStateData from a remote entity
	msgCounter, err := requestData(featureLocal, featureRemote, model.FunctionTypeDeviceDiagnosisStateData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return msgCounter, nil
}

// get the current diagnosis state for an device entity
func GetDeviceDiagnosisState(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*DeviceDiagnosisType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := featureRemote.Data(model.FunctionTypeDeviceDiagnosisStateData).(*model.DeviceDiagnosisStateDataType)

	details := &DeviceDiagnosisType{}
	if data.OperatingState != nil {
		details.OperatingState = *data.OperatingState
	}
	if data.PowerSupplyCondition != nil {
		details.PowerSupplyCondition = *data.PowerSupplyCondition
	}

	return details, nil
}

func SendDeviceDiagnosisState(service *service.EEBUSService, entity *spine.EntityRemoteImpl, operatingState *model.DeviceDiagnosisStateDataType) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)

	_, _ = featureLocal.NotifyData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote)
}
