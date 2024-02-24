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

	return features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func DeviceConfiguration(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceConfiguration, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func DeviceDiagnosis(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceDiagnosis, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func DeviceDiagnosisServer(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.DeviceDiagnosis, error) {
	localEntity := localCemEntity(service)

	return features.NewDeviceDiagnosis(model.RoleTypeServer, model.RoleTypeClient, localEntity, remoteEntity)
}

func ElectricalConnection(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.ElectricalConnection, error) {
	localEntity := localCemEntity(service)

	return features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func Identification(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.Identification, error) {
	localEntity := localCemEntity(service)

	return features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func Measurement(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.Measurement, error) {
	localEntity := localCemEntity(service)

	return features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func LoadControl(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.LoadControl, error) {
	localEntity := localCemEntity(service)

	return features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func TimeSeries(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.TimeSeries, error) {
	localEntity := localCemEntity(service)

	return features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func IncentiveTable(service eebusapi.ServiceInterface, remoteEntity api.EntityRemoteInterface) (*features.IncentiveTable, error) {
	localEntity := localCemEntity(service)

	return features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}
