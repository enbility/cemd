package uclpc

import (
	"time"

	"github.com/enbility/cemd/api"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCLPCSuite) Test_LoadControlLimit() {
	data, err := s.sut.LoadControlLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	data, err = s.sut.LoadControlLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.LoadControlLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: eebusutil.Ptr(true),
				IsLimitActive:     eebusutil.Ptr(false),
				Value:             model.NewScaledNumberType(6000),
				TimePeriod: &model.TimePeriodType{
					EndTime: model.NewAbsoluteOrRelativeTimeType("PT2H"),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.LoadControlLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6000.0, data.Value)
	assert.Equal(s.T(), true, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)
}

func (s *UCLPCSuite) Test_WriteLoadControlLimit() {
	limit := api.LoadLimit{
		Value:    6000,
		IsActive: true,
		Duration: 0,
	}
	_, err := s.sut.WriteLoadControlLimit(s.mockRemoteEntity, limit)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeConsume),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteLoadControlLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				IsLimitChangeable: eebusutil.Ptr(true),
				IsLimitActive:     eebusutil.Ptr(false),
				Value:             model.NewScaledNumberType(6000),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteLoadControlLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)

	limit.Duration = time.Duration(time.Hour * 2)
	_, err = s.sut.WriteLoadControlLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)
}

func (s *UCLPCSuite) Test_FailsafeConsumptionActivePowerLimit() {
	data, err := s.sut.FailsafeConsumptionActivePowerLimit(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(4000),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeConsumptionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 4000.0, data)
}

func (s *UCLPCSuite) Test_WriteFailsafeConsumptionActivePowerLimit() {
	_, err := s.sut.WriteFailsafeConsumptionActivePowerLimit(s.mockRemoteEntity, 6000)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeConsumptionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)
}

func (s *UCLPCSuite) Test_FailsafeDurationMinimum() {
	data, err := s.sut.FailsafeDurationMinimum(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), time.Duration(0), data)

	keyData = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Duration: model.NewDurationType(time.Hour * 2),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeDurationMinimum(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), time.Duration(time.Hour*2), data)
}

func (s *UCLPCSuite) Test_WriteFailsafeDurationMinimum() {
	_, err := s.sut.WriteFailsafeDurationMinimum(s.mockRemoteEntity, time.Duration(time.Hour*2))
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.Nil(s.T(), err)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*2))
	assert.Nil(s.T(), err)

	_, err = s.sut.WriteFailsafeDurationMinimum(s.monitoredEntity, time.Duration(time.Hour*1))
	assert.NotNil(s.T(), err)
}

func (s *UCLPCSuite) Test_PowerConsumptionNominalMax() {
	data, err := s.sut.PowerConsumptionNominalMax(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.PowerConsumptionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicListData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       eebusutil.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  eebusutil.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     eebusutil.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerConsumptionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerConsumptionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 8000.0, data)
}
