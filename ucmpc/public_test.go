package ucmpc

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCMPCSuite) Test_Power() {
	data, err := s.sut.Power(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Power(s.monitoredEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.monitoredEntity)
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

	data, err = s.sut.Power(s.monitoredEntity)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, elDescData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Power(s.monitoredEntity)
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

	data, err = s.sut.Power(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMPCSuite) Test_PowerPerPhase() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(1)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(2)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
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

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
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

	data, err = s.sut.PowerPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)
}

func (s *UCMPCSuite) Test_EnergyConsumed() {
	data, err := s.sut.EnergyConsumed(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACEnergyConsumed),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
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

	data, err = s.sut.EnergyConsumed(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMPCSuite) Test_EnergyProduced() {
	data, err := s.sut.EnergyProduced(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACEnergyProduced),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
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

	data, err = s.sut.EnergyProduced(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCMPCSuite) Test_CurrentPerPhase() {
	data, err := s.sut.CurrentPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
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

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
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

	data, err = s.sut.CurrentPerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{10, 10, 10}, data)
}

func (s *UCMPCSuite) Test_VoltagePerPhase() {
	data, err := s.sut.VoltagePerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
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

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, elParamData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.VoltagePerPhase(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []float64{230, 230, 230}, data)
}

func (s *UCMPCSuite) Test_Frequency() {
	data, err := s.sut.Frequency(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.Frequency(s.monitoredEntity)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Frequency(s.monitoredEntity)
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

	data, err = s.sut.Frequency(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 50.0, data)
}
