package ucevcem

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCEVCEMSuite) Test_EVConnectedPhases() {
	data, err := s.sut.PhasesConnected(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData = &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				AcConnectedPhases:      eebusutil.Ptr(uint(1)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PhasesConnected(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(1), data)
}

func (s *UCEVCEMSuite) Test_EVCurrentsPerPhase() {
	data, err := s.sut.CurrentsPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.CurrentsPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:              eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentsPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	data, err = s.sut.CurrentsPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])
}

func (s *UCEVCEMSuite) Test_EVPowerPerPhase_Power() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:              eebusutil.Ptr(model.ScopeTypeTypeACPower),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypePower),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data[0])
}

func (s *UCEVCEMSuite) Test_EVPowerPerPhase_Current() {
	data, err := s.sut.PowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	paramDesc := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:              eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeCurrent),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACCurrent),
			},
		},
	}

	rFeature = s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

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

	data, err = s.sut.PowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2300.0, data[0])
}

func (s *UCEVCEMSuite) Test_EVChargedEnergy() {
	data, err := s.sut.EnergyCharged(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measDesc := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeCharge),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, measDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	measData := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EnergyCharged(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data)
}
