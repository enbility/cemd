package ucevcc

import (
	"fmt"
	"testing"
	"time"

	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusmocks "github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/util"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/ship-go/cert"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEvCCSuite(t *testing.T) {
	suite.Run(t, new(EvCCSuite))
}

type EvCCSuite struct {
	suite.Suite

	sut *UCEvCC

	service eebusapi.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *mocks.EntityRemoteInterface
	evEntity         spineapi.EntityRemoteInterface
}

func (s *EvCCSuite) SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event api.UseCaseEventType) {
}

func (s *EvCCSuite) BeforeTest(suiteName, testName string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := eebusapi.NewConfiguration(
		"test", "test", "test", "test",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, 230.0, time.Second*4)

	serviceHandler := eebusmocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()

	mockRemoteDevice := mocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = mocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := mocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()

	var entities []spineapi.EntityRemoteInterface

	s.remoteDevice, entities = setupDevices(s.service, s.T())
	s.sut = NewUCEvCC(s.service, s.service.LocalService(), s)
	s.sut.AddFeatures()
	s.sut.AddUseCase()
	s.evEntity = entities[1]
}

func (s *EvCCSuite) Test_EVConnected() {
	data := s.sut.EVConnected(nil)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.mockRemoteEntity)
	assert.Equal(s.T(), false, data)

	data = s.sut.EVConnected(s.evEntity)
	assert.Equal(s.T(), true, data)
}

func (s *EvCCSuite) Test_EVCommunicationStandard() {
	data, err := s.sut.EVCommunicationStandard(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), UcEVCCUnknownCommunicationStandard, data)

	data, err = s.sut.EVCommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), UcEVCCUnknownCommunicationStandard, data)

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

	data, err = s.sut.EVCommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), UcEVCCUnknownCommunicationStandard, data)

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

	data, err = s.sut.EVCommunicationStandard(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), UcEVCCUnknownCommunicationStandard, data)

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

	data, err = s.sut.EVCommunicationStandard(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(model.DeviceConfigurationKeyValueStringTypeISO151182ED2), data)
}

func (s *EvCCSuite) Test_EVAsymmetricChargingSupported() {
	data, err := s.sut.EVAsymmetricChargingSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.EVAsymmetricChargingSupported(s.evEntity)
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

	data, err = s.sut.EVAsymmetricChargingSupported(s.evEntity)
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

	data, err = s.sut.EVAsymmetricChargingSupported(s.evEntity)
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

	data, err = s.sut.EVAsymmetricChargingSupported(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}

func (s *EvCCSuite) Test_EVIdentification() {
	data, err := s.sut.EVIdentifications(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []IdentificationItem(nil), data)

	data, err = s.sut.EVIdentifications(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), []IdentificationItem(nil), data)

	data, err = s.sut.EVIdentifications(s.evEntity)
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

	data, err = s.sut.EVIdentifications(s.evEntity)
	assert.Nil(s.T(), err)
	resultData := []IdentificationItem{{Value: "test", ValueType: model.IdentificationTypeTypeEui64}}
	assert.Equal(s.T(), resultData, data)
}

func (s *EvCCSuite) Test_EVManufacturerData() {
	device, serial, err := s.sut.EVManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.EVManufacturerData(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.EVManufacturerData(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.EVManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData = &model.DeviceClassificationManufacturerDataType{
		DeviceName:   eebusutil.Ptr(model.DeviceClassificationStringType("test")),
		SerialNumber: eebusutil.Ptr(model.DeviceClassificationStringType("12345")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.EVManufacturerData(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", device)
	assert.Equal(s.T(), "12345", serial)
}

func (s *EvCCSuite) Test_EVConnectedPhases() {
	data, err := s.sut.EVConnectedPhases(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	data, err = s.sut.EVConnectedPhases(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	data, err = s.sut.EVConnectedPhases(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData := &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EVConnectedPhases(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(0), data)

	descData = &model.ElectricalConnectionDescriptionListDataType{
		ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
			{
				ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				AcConnectedPhases:      util.Ptr(uint(1)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.EVConnectedPhases(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uint(1), data)
}

func (s *EvCCSuite) Test_EVCurrentLimits() {
	minData, maxData, defaultData, err := s.sut.EVCurrentLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.sut.EVCurrentLimits(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = s.sut.EVCurrentLimits(s.evEntity)
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

	minData, maxData, defaultData, err = s.sut.EVCurrentLimits(s.evEntity)
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

			minData, maxData, defaultData, err = s.sut.EVCurrentLimits(s.evEntity)
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

func (s *EvCCSuite) Test_EVInSleepMode() {
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

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService eebusapi.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeIdentification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(3, localEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(4, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(5, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	localEntity.AddFeature(f)

	writeHandler := shipmocks.NewShipConnectionDataWriterInterface(t)
	writeHandler.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

	var clientRemoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeDeviceConfiguration,
			[]model.FunctionType{
				model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
				model.FunctionTypeDeviceConfigurationKeyValueListData,
			},
		},
		{model.FeatureTypeTypeIdentification,
			[]model.FunctionType{
				model.FunctionTypeIdentificationListData,
			},
		},
		{model.FeatureTypeTypeDeviceClassification,
			[]model.FunctionType{
				model.FunctionTypeDeviceClassificationManufacturerData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			[]model.FunctionType{
				model.FunctionTypeElectricalConnectionDescriptionListData,
			},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			[]model.FunctionType{
				model.FunctionTypeDeviceDiagnosisStateData,
			},
		},
	}

	remoteDeviceName := "remote"

	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for index, feature := range clientRemoteFeatures {
		supportedFcts := []model.FunctionPropertyType{}
		for _, fct := range feature.supportedFcts {
			supportedFct := model.FunctionPropertyType{
				Function: eebusutil.Ptr(fct),
				PossibleOperations: &model.PossibleOperationsType{
					Read: &model.PossibleOperationsReadType{},
				},
			}
			supportedFcts = append(supportedFcts, supportedFct)
		}

		featureInformation := model.NodeManagementDetailedDiscoveryFeatureInformationType{
			Description: &model.NetworkManagementFeatureDescriptionDataType{
				FeatureAddress: &model.FeatureAddressType{
					Device:  eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
					Entity:  []model.AddressEntityType{1, 1},
					Feature: eebusutil.Ptr(model.AddressFeatureType(index)),
				},
				FeatureType:       eebusutil.Ptr(feature.featureType),
				Role:              eebusutil.Ptr(model.RoleTypeServer),
				SupportedFunction: supportedFcts,
			},
		}
		featureInformations = append(featureInformations, featureInformation)
	}

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: eebusutil.Ptr(model.EntityTypeTypeEV),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1, 1},
					},
					EntityType: eebusutil.Ptr(model.EntityTypeTypeEV),
				},
			},
		},
		FeatureInformation: featureInformations,
	}

	entities, err := remoteDevice.AddEntityAndFeatures(true, detailedData)
	if err != nil {
		fmt.Println(err)
	}

	localDevice.AddRemoteDeviceForSki(remoteSki, remoteDevice)

	return remoteDevice, entities
}
