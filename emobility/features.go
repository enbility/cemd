package emobility

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *EMobility) deviceClassification(remoteEntity api.EntityRemoteInterface) (*features.DeviceClassification, error) {
	localEntity := e.localCemEntity()

	return features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) deviceConfiguration(remoteEntity api.EntityRemoteInterface) (*features.DeviceConfiguration, error) {
	localEntity := e.localCemEntity()

	return features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) deviceDiagnosis(remoteEntity api.EntityRemoteInterface) (*features.DeviceDiagnosis, error) {
	localEntity := e.localCemEntity()

	return features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) electricalConnection(remoteEntity api.EntityRemoteInterface) (*features.ElectricalConnection, error) {
	localEntity := e.localCemEntity()

	return features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) measurement(remoteEntity api.EntityRemoteInterface) (*features.Measurement, error) {
	localEntity := e.localCemEntity()

	return features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) loadControl(remoteEntity api.EntityRemoteInterface) (*features.LoadControl, error) {
	localEntity := e.localCemEntity()

	return features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) identification(remoteEntity api.EntityRemoteInterface) (*features.Identification, error) {
	localEntity := e.localCemEntity()

	return features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) timeSeries(remoteEntity api.EntityRemoteInterface) (*features.TimeSeries, error) {
	localEntity := e.localCemEntity()

	return features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}

func (e *EMobility) incentiveTable(remoteEntity api.EntityRemoteInterface) (*features.IncentiveTable, error) {
	localEntity := e.localCemEntity()

	return features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
}
