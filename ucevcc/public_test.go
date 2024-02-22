package ucevcc

import (
	"testing"

	"github.com/enbility/cemd/api"
	"github.com/enbility/eebus-go/util"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *EVCCSuite) Test_EVCurrentChargeState() {
	data, err := s.sut.CurrentChargeState(s.mockRemoteEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeUnplugged, data)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeUnknown, data)

	stateData := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeActive, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypePaused, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFailure),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeError, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFinished),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeFinished, data)

	stateData = &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeInAlarm),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CurrentChargeState(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.EVChargeStateTypeUnknown, data)
}

func (s *EVCCSuite) Test_EVConnected() {
	data := s.sut.EVConnected(nil)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.mockRemoteEntity)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.evEntity)
	assert.Equal(s.T(), false, data)

	stateData := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, stateData, nil, nil)
	assert.Nil(s.T(), fErr)

	data = s.sut.EVConnected(s.evEntity)
	assert.Equal(s.T(), true, data)
}

func (s *EVCCSuite) Test_EVCommunicationStandard() {
	data, err := s.sut.CommunicationStandard(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.UCEVCCCommunicationStandardUnknown, data)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.UCEVCCCommunicationStandardUnknown, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.UCEVCCCommunicationStandardUnknown, data)

	descData = &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), api.UCEVCCCommunicationStandardUnknown, data)

	devData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					String: eebusutil.Ptr(model.DeviceConfigurationKeyValueStringTypeISO151182ED2),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, devData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.CommunicationStandard(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(model.DeviceConfigurationKeyValueStringTypeISO151182ED2), data)
}

func (s *EVCCSuite) Test_EVAsymmetricChargingSupported() {
	data, err := s.sut.AsymmetricChargingSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.AsymmetricChargingSupported(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupported(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData = &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupported(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	devData := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: &model.DeviceConfigurationKeyValueValueType{
					Boolean: eebusutil.Ptr(true),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, devData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.AsymmetricChargingSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}

func (s *EVCCSuite) Test_EVIdentification() {
	data, err := s.sut.Identifications(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []IdentificationItem(nil), data)

	data, err = s.sut.Identifications(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []IdentificationItem(nil), data)

	data, err = s.sut.Identifications(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []IdentificationItem(nil), data)

	idData := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(0)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIdentification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeIdentificationListData, idData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.Identifications(s.evEntity)
	assert.Nil(s.T(), err)
	resultData := []IdentificationItem{{Value: "test", ValueType: model.IdentificationTypeTypeEui64}}
	assert.Equal(s.T(), resultData, data)
}

func (s *EVCCSuite) Test_EVManufacturerData() {
	device, serial, err := s.sut.ManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.ManufacturerData(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.ManufacturerData(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.ManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData = &model.DeviceClassificationManufacturerDataType{
		DeviceName:   eebusutil.Ptr(model.DeviceClassificationStringType("test")),
		SerialNumber: eebusutil.Ptr(model.DeviceClassificationStringType("12345")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.ManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", device)
	assert.Equal(s.T(), "12345", serial)
}

func (s *EVCCSuite) Test_EVCurrentLimits() {
	minData, maxData, defaultData, err := s.sut.CurrentLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	minData, maxData, defaultData, err = s.sut.CurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	type permittedStruct struct {
		defaultExists                      bool
		defaultValue, expectedDefaultValue float64
		minValue, expectedMinValue         float64
		maxValue, expectedMaxValue         float64
	}

	tests := []struct {
		name      string
		permitted []permittedStruct
	}{
		{
			"1 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2, 16, 16},
			},
		},
		{
			"1 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"1 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
			},
		},
		{
			"3 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			for index, data := range tc.permitted {
				item := model.ScaledNumberSetType{
					Range: []model.ScaledNumberRangeType{
						{
							Min: model.NewScaledNumberType(data.minValue),
							Max: model.NewScaledNumberType(data.maxValue),
						},
					},
				}
				if data.defaultExists {
					item.Value = []model.ScaledNumberType{*model.NewScaledNumberType(data.defaultValue)}
				}
				permittedData = append(permittedData, item)

				permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(index)),
					PermittedValueSet:      permittedData,
				}
				dataSet = append(dataSet, permittedItem)
			}

			permData := &model.ElectricalConnectionPermittedValueSetListDataType{
				ElectricalConnectionPermittedValueSetData: dataSet,
			}

			fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
			assert.Nil(s.T(), fErr)

			minData, maxData, defaultData, err = s.sut.CurrentLimits(s.evEntity)
			assert.Nil(s.T(), err)

			assert.Nil(s.T(), err)
			assert.Equal(s.T(), len(tc.permitted), len(minData))
			assert.Equal(s.T(), len(tc.permitted), len(maxData))
			assert.Equal(s.T(), len(tc.permitted), len(defaultData))
			for index, item := range tc.permitted {
				assert.Equal(s.T(), item.expectedMinValue, minData[index])
				assert.Equal(s.T(), item.expectedMaxValue, maxData[index])
				assert.Equal(s.T(), item.expectedDefaultValue, defaultData[index])
			}
		})
	}
}

func (s *EVCCSuite) Test_EVInSleepMode() {
	data, err := s.sut.EVInSleepMode(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.EVInSleepMode(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData := &model.DeviceDiagnosisStateDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EVInSleepMode(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	descData = &model.DeviceDiagnosisStateDataType{
		OperatingState: eebusutil.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EVInSleepMode(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
