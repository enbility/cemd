package util

import (
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_ManufacturerData() {
	entityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}

	data, err := ManufacturerData(s.service, s.mockRemoteEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = ManufacturerData(s.service, s.monitoredEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.DeviceClassificationManufacturerDataType{

		DeviceName:   util.Ptr(model.DeviceClassificationStringType("deviceName")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("deviceCode")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serialNumber")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	assert.NotNil(s.T(), rFeature)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)
	data, err = ManufacturerData(s.service, s.monitoredEntity, entityTypes)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), "deviceName", data.DeviceName)
	assert.Equal(s.T(), "deviceCode", data.DeviceCode)
	assert.Equal(s.T(), "serialNumber", data.SerialNumber)
	assert.Equal(s.T(), "", data.SoftwareRevision)
}
