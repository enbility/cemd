package ucevsecc

import (
	"fmt"
	"testing"
	"time"

	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	eebusmocks "github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
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
	suite.Run(t, new(EvseCCSuite))
}

type EvseCCSuite struct {
	suite.Suite

	sut *UCEvseCC

	service eebusapi.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *mocks.EntityRemoteInterface
	evseEntity       spineapi.EntityRemoteInterface
}

func (s *EvseCCSuite) SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event api.UseCaseEventType) {
}

func (s *EvseCCSuite) BeforeTest(suiteName, testName string) {
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
	s.sut = NewUCEvseCC(s.service, s.service.LocalService(), s)
	s.sut.AddFeatures()
	s.sut.AddUseCase()
	s.evseEntity = entities[0]
}

func (s *EvseCCSuite) Test_EVSEManufacturerData() {
	device, serial, err := s.sut.EVSEManufacturerData(nil)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.EVSEManufacturerData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	device, serial, err = s.sut.EVSEManufacturerData(s.evseEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData := &model.DeviceClassificationManufacturerDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.EVSEManufacturerData(s.evseEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "", device)
	assert.Equal(s.T(), "", serial)

	descData = &model.DeviceClassificationManufacturerDataType{
		DeviceName:   eebusutil.Ptr(model.DeviceClassificationStringType("test")),
		SerialNumber: eebusutil.Ptr(model.DeviceClassificationStringType("12345")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	device, serial, err = s.sut.EVSEManufacturerData(s.evseEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test", device)
	assert.Equal(s.T(), "12345", serial)
}

func (s *EvseCCSuite) Test_EVSEOperatingState() {
	data, errCode, err := s.sut.EVSEOperatingState(nil)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), nil, err)

	data, errCode, err = s.sut.EVSEOperatingState(s.mockRemoteEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	data, errCode, err = s.sut.EVSEOperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.NotNil(s.T(), err)

	descData := &model.DeviceDiagnosisStateDataType{}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evseEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.EVSEOperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeNormalOperation, data)
	assert.Equal(s.T(), "", errCode)
	assert.Nil(s.T(), err)

	descData = &model.DeviceDiagnosisStateDataType{
		OperatingState: eebusutil.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
		LastErrorCode:  eebusutil.Ptr(model.LastErrorCodeType("error")),
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, errCode, err = s.sut.EVSEOperatingState(s.evseEntity)
	assert.Equal(s.T(), model.DeviceDiagnosisOperatingStateTypeStandby, data)
	assert.Equal(s.T(), "error", errCode)
	assert.Nil(s.T(), err)
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService eebusapi.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	localEntity.AddFeature(f)

	writeHandler := shipmocks.NewShipConnectionDataWriterInterface(t)
	writeHandler.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

	var clientRemoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeDeviceClassification,
			[]model.FunctionType{
				model.FunctionTypeDeviceClassificationManufacturerData,
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
					Entity:  []model.AddressEntityType{1},
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
					EntityType: eebusutil.Ptr(model.EntityTypeTypeEVSE),
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
