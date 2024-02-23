package ucevsecc

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCEVSECCSuite) Test_EVSEManufacturerData() {
	device, serial, err := s.sut.ManufacturerData(nil)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.ManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.ManufacturerData(s.evseEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.ManufacturerData(s.evseEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData = &model.DeviceClassificationManufacturerDataType{
		DeviceName:   eebusutil.Ptr(model.DeviceClassificationStringType("test")),
		SerialNumber: eebusutil.Ptr(model.DeviceClassificationStringType("12345")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.ManufacturerData(s.evseEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", device)
	assert.Equal(s.T(), "12345", serial)
}

func (s *UCEVSECCSuite) Test_EVSEOperatingState() {
	data, errCode, err := s.sut.OperatingState(nil)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), nil, err)

	data, errCode, err = s.sut.OperatingState(s.mockRemoteEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceDiagnosisStateDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), err)

	descData = &model.DeviceDiagnosisStateDataType{
		OperatingState: eebusutil.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
		LastErrorCode:  eebusutil.Ptr(model.LastErrorCodeType("error")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.OperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeStandby, data)
	assert.Equal(s.T(), "error", errCode)
	assert.Nil(s.T(), err)
}
