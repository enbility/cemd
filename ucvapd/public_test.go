package ucvapd

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCVAPDSuite) Test_CurrentProductionPower() {
	data, err := s.sut.CurrentProductionPower(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.CurrentProductionPower(s.pvEntity)
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

	measurementFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := measurementFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentProductionPower(s.pvEntity)
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

	fErr = measurementFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentProductionPower(s.pvEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCVAPDSuite) Test_NominalPeakPower() {
	data, err := s.sut.NominalPeakPower(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.NominalPeakPower(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	confData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem),
			},
		},
	}

	confFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := confFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, confData, nil, nil)
	assert.Nil(s.T(), fErr)

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
	fErr = confFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, keyData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.NominalPeakPower(s.pvEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}

func (s *UCVAPDSuite) Test_TotalPVYield() {
	data, err := s.sut.TotalPVYield(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.TotalPVYield(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId:   eebusutil.Ptr(model.MeasurementIdType(0)),
				MeasurementType: eebusutil.Ptr(model.MeasurementTypeTypeEnergy),
				CommodityType:   eebusutil.Ptr(model.CommodityTypeTypeElectricity),
				ScopeType:       eebusutil.Ptr(model.ScopeTypeTypeACYieldTotal),
			},
		},
	}

	measurementFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := measurementFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.TotalPVYield(s.pvEntity)
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

	fErr = measurementFeature.UpdateData(model.FunctionTypeMeasurementListData, measData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.TotalPVYield(s.pvEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data)
}
