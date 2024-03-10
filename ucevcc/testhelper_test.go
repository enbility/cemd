package ucevcc

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
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestEVCCSuite(t *testing.T) {
	suite.Run(t, new(UCEVCCSuite))
}

type UCEVCCSuite struct {
	suite.Suite

	sut *UCEVCC

	service eebusapi.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockSender       *mocks.SenderInterface
	mockRemoteEntity *mocks.EntityRemoteInterface
	evEntity         spineapi.EntityRemoteInterface
}

func (s *UCEVCCSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}

func (s *UCEVCCSuite) BeforeTest(suiteName, testName string) {
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
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()
	mockRemoteFeature.EXPECT().Address().Return(&model.FeatureAddressType{}).Maybe()
	mockRemoteFeature.EXPECT().Operations().Return(nil).Maybe()

	s.sut = NewUCEVCC(s.service, s.Event)
	s.sut.AddFeatures()
	s.sut.AddUseCase()

	var entities []spineapi.EntityRemoteInterface
	s.remoteDevice, s.mockSender, entities = setupDevices(s.service, s.T())
	s.evEntity = entities[1]
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService eebusapi.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	*mocks.SenderInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()

	mockSender := mocks.NewSenderInterface(t)
	defaultMsgCounter := model.MsgCounterType(100)
	mockSender.
		EXPECT().
		Request(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&defaultMsgCounter, nil).
		Maybe()
	mockSender.
		EXPECT().
		Subscribe(mock.Anything, mock.Anything, mock.Anything).
		Return(&defaultMsgCounter, nil).
		Maybe()
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, mockSender)

	remoteDeviceName := "remote"

	var remoteFeatures = []struct {
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
				model.FunctionTypeElectricalConnectionParameterDescriptionListData,
				model.FunctionTypeElectricalConnectionPermittedValueSetListData,
			},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			[]model.FunctionType{
				model.FunctionTypeDeviceDiagnosisStateData,
			},
		},
	}

	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for index, feature := range remoteFeatures {
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
	remoteDevice.UpdateDevice(detailedData.DeviceInformation.Description)

	localDevice.AddRemoteDeviceForSki(remoteSki, remoteDevice)

	return remoteDevice, mockSender, entities
}
