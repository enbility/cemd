package util

import (
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func localCemEntity(service eebusapi.ServiceInterface) api.EntityLocalInterface {
	localDevice := service.LocalDevice()

	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	return localEntity
}

func DeviceClassification(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceClassification, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceClassification(localEntity, remoteEntity)
}

func DeviceConfiguration(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceConfiguration, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceConfiguration(localEntity, remoteEntity)
}

func DeviceDiagnosis(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceDiagnosis, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceDiagnosis(localEntity, remoteEntity)
}

func DeviceDiagnosisServer(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceDiagnosis, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceDiagnosis(localEntity, remoteEntity)
}

func ElectricalConnection(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.ElectricalConnection, error) {
	localEntity := localCemEntity(service)

	return features.NewElectricalConnection(localEntity, remoteEntity)
}

func Identification(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.Identification, error) {
	localEntity := localCemEntity(service)

	return features.NewIdentification(localEntity, remoteEntity)
}

func Measurement(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.Measurement, error) {
	localEntity := localCemEntity(service)

	return features.NewMeasurement(localEntity, remoteEntity)
}

func LoadControl(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.LoadControl, error) {
	localEntity := localCemEntity(service)

	return features.NewLoadControl(localEntity, remoteEntity)
}

func TimeSeries(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.TimeSeries, error) {
	localEntity := localCemEntity(service)

	return features.NewTimeSeries(localEntity, remoteEntity)
}

func IncentiveTable(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.IncentiveTable, error) {
	localEntity := localCemEntity(service)

	return features.NewIncentiveTable(localEntity, remoteEntity)
}
