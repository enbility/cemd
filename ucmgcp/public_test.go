package ucmgcp

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCMGCPSuite) Test_PowerLimitationFactor() {
	data, err := s.sut.PowerLimitationFactor(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypePvCurtailmentLimitFactor),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	keyData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					ScaledNumber: model.NewScaledNumberType(10),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerLimitationFactor(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMGCPSuite) Test_Power() {
	data, err := s.sut.Power(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: eebusutil.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
			},
		},
	}

	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMGCPSuite) Test_EnergyFeedIn() {
	data, err := s.sut.EnergyFeedIn(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeGridFeedIn),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyFeedIn(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMGCPSuite) Test_EnergyConsumed() {
	data, err := s.sut.EnergyConsumed(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeGridConsumption),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMGCPSuite) Test_CurrentPerPhase() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(1)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(2)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	elDescData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId:  eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				PositiveEnergyDirection: eebusutil.Ptr(model.EnergyDirectionTypeConsume),
			},
		},
	}

	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)

}

func (s *UCMGCPSuite) Test_VoltagePerPhase() {
	data, err := s.sut.VoltagePerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(1)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACVoltage),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(2)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeVoltage),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACVoltage),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(230),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(230),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(230),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	elParamData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230, 230, 230}, data)
}

func (s *UCMGCPSuite) Test_Frequency() {
	data, err := s.sut.Frequency(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Frequency(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeFrequency),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACFrequency),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.smgwEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.smgwEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(50),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.smgwEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data)
}
