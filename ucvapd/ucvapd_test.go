package ucvapd

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCVAPDSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *UCVAPDSuite) Test_IsUseCaseSupported() {
	data, err := s.sut.IsUseCaseSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: eebusutil.Ptr(model.UseCaseActorTypePVSystem),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      eebusutil.Ptr(model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData),
						UseCaseAvailable: eebusutil.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1, 2, 3},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	confData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem),
			},
		},
	}

	confFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr = confFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, confData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	elData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	elFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = elFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, elData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	fErr = elFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeACYieldTotal),
			},
		},
	}

	measurementFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.pvEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr = measurementFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.pvEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
