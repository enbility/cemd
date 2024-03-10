package util

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_GetLocalDeviceConfigurationDescriptionForKeyName() {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit

	data := GetLocalDeviceConfigurationDescriptionForKeyName(s.service, keyName)
	assert.Nil(s.T(), data.KeyId)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(keyName),
			},
		},
	}
	feature.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData)

	data = GetLocalDeviceConfigurationDescriptionForKeyName(s.service, keyName)
	assert.NotNil(s.T(), data.KeyId)
}

func (s *UtilSuite) Test_GetLocalDeviceConfigurationKeyValueForId() {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit

	data := GetLocalDeviceConfigurationKeyValueForKeyName(s.service, keyName)
	assert.Nil(s.T(), data.KeyId)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(keyName),
			},
		},
	}
	feature.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData)

	data = GetLocalDeviceConfigurationKeyValueForKeyName(s.service, keyName)
	assert.Nil(s.T(), data.KeyId)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
			},
		},
	}

	feature.SetData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData)

	data = GetLocalDeviceConfigurationKeyValueForKeyName(s.service, keyName)
	assert.NotNil(s.T(), data.KeyId)
}

func (s *UtilSuite) Test_SetLocalDeviceConfigurationKeyValue() {
	keyName := model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit
	changeable := false
	value := model.DeviceConfigurationKeyValueValueType{
		ScaledNumber: model.NewScaledNumberType(10),
	}

	err := SetLocalDeviceConfigurationKeyValue(s.service, keyName, changeable, value)
	assert.NotNil(s.T(), err)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(keyName),
			},
		},
	}
	feature.SetData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData)

	err = SetLocalDeviceConfigurationKeyValue(s.service, keyName, changeable, value)
	assert.Nil(s.T(), err)

	data := GetLocalDeviceConfigurationKeyValueForKeyName(s.service, keyName)
	assert.NotNil(s.T(), data.KeyId)
	assert.Equal(s.T(), uint(0), uint(*data.KeyId))
	assert.NotNil(s.T(), data.Value)
	assert.NotNil(s.T(), data.Value.ScaledNumber)
	assert.Equal(s.T(), 10.0, data.Value.ScaledNumber.GetValue())
}
