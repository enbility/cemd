package uclpp

import (
	"time"

	"github.com/enbility/cemd/api"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCLPPSuite) Test_LoadControlLimit() {
	data, err := s.sut.ProductionLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data.Value)
	assert.Equal(s.T(), false, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
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
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
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

	data, err = s.sut.ProductionLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 6000.0, data.Value)
	assert.Equal(s.T(), true, data.IsChangeable)
	assert.Equal(s.T(), false, data.IsActive)
}

func (s *UCLPPSuite) Test_WriteLoadControlLimit() {
	limit := api.LoadLimit{
		Value:    6000,
		IsActive: true,
		Duration: 0,
	}
	_, err := s.sut.WriteProductionLimit(s.mockRemoteEntity, limit)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit)
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

	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)

	limit.Duration = time.Duration(time.Hour * 2)
	_, err = s.sut.WriteProductionLimit(s.monitoredEntity, limit)
	assert.NotNil(s.T(), err)
}

func (s *UCLPPSuite) Test_FailsafeProductionActivePowerLimit() {
	data, err := s.sut.FailsafeProductionActivePowerLimit(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
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

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
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

	data, err = s.sut.FailsafeProductionActivePowerLimit(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 4000.0, data)
}

func (s *UCLPPSuite) Test_WriteFailsafeProductionActivePowerLimit() {
	_, err := s.sut.WriteFailsafeProductionActivePowerLimit(s.mockRemoteEntity, 6000)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
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

	_, err = s.sut.WriteFailsafeProductionActivePowerLimit(s.monitoredEntity, 6000)
	assert.Nil(s.T(), err)
}

func (s *UCLPPSuite) Test_FailsafeDurationMinimum() {
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

func (s *UCLPPSuite) Test_WriteFailsafeDurationMinimum() {
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

func (s *UCLPPSuite) Test_PowerProductionNominalMax() {
	data, err := s.sut.PowerProductionNominalMax(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.PowerProductionNominalMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				CharacteristicId:       eebusutil.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  eebusutil.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
				CharacteristicType:     eebusutil.Ptr(model.ElectricalConnectionCharacteristicTypeTypePowerProductionNominalMax),
				Value:                  model.NewScaledNumberType(8000),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionCharacteristicListData, charData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerProductionNominalMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 8000.0, data)
}
