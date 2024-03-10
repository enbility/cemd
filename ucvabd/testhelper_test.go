package ucvabd

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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestVABDSuite(t *testing.T) {
	suite.Run(t, new(UCVABDSuite))
}

type UCVABDSuite struct {
	suite.Suite

	sut *UCVABD

	service eebusapi.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *mocks.EntityRemoteInterface
	batteryEntity    spineapi.EntityRemoteInterface
}

func (s *UCVABDSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}

func (s *UCVABDSuite) BeforeTest(suiteName, testName string) {
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

	s.sut = NewUCVABD(s.service, s.Event)
	s.sut.AddFeatures()
	s.sut.AddUseCase()

	s.remoteDevice, s.batteryEntity = setupDevices(s.service, s.T())
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService eebusapi.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()

	writeHandler := shipmocks.NewShipConnectionDataWriterInterface(t)
	writeHandler.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

	remoteDeviceName := "remote"

	var remoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeMeasurement,
			[]model.FunctionType{
				model.FunctionTypeMeasurementDescriptionListData,
				model.FunctionTypeMeasurementConstraintsListData,
				model.FunctionTypeMeasurementListData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			[]model.FunctionType{
				model.FunctionTypeElectricalConnectionParameterDescriptionListData,
				model.FunctionTypeElectricalConnectionDescriptionListData,
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
					EntityType: eebusutil.Ptr(model.EntityTypeTypeElectricityStorageSystem),
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

	return remoteDevice, entities[0]
}
