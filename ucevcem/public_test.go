package ucevcem

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

func TestEvCEMSuite(t *testing.T) {
	suite.Run(t, new(EvCEMSuite))
}

type EvCEMSuite struct {
	suite.Suite

	sut *UCEvCEM

	service eebusapi.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *mocks.EntityRemoteInterface
	evEntity         spineapi.EntityRemoteInterface
}

func (s *EvCEMSuite) SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event api.UseCaseEventType) {
}

func (s *EvCEMSuite) BeforeTest(suiteName, testName string) {
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
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()

	var entities []spineapi.EntityRemoteInterface

	s.remoteDevice, entities = setupDevices(s.service, s.T())
	s.sut = NewUCEvCEM(s.service, s.service.LocalService(), s)
	s.sut.AddFeatures()
	s.sut.AddUseCase()
	s.evEntity = entities[1]
}

func (s *EvCEMSuite) Test_EVCurrentsPerPhase() {
	data, err := s.sut.EVCurrentsPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.EVCurrentsPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.EVCurrentsPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVCurrentsPerPhase(s.evEntity)
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

	data, err = s.sut.EVCurrentsPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 10.0, data[0])
}

func (s *EvCEMSuite) Test_EVPowerPerPhase_Power() {
	data, err := s.sut.EVPowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data[0])
}

func (s *EvCEMSuite) Test_EVPowerPerPhase_Current() {
	data, err := s.sut.EVPowerPerPhase(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
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

	data, err = s.sut.EVPowerPerPhase(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 2300.0, data[0])
}

func (s *EvCEMSuite) Test_EVChargedEnergy() {
	data, err := s.sut.EVChargedEnergy(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EVChargedEnergy(s.evEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, data)

	data, err = s.sut.EVChargedEnergy(s.evEntity)
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

	data, err = s.sut.EVChargedEnergy(s.evEntity)
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

	data, err = s.sut.EVChargedEnergy(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, data)
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService eebusapi.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)

	writeHandler := shipmocks.NewShipConnectionDataWriterInterface(t)
	writeHandler.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

	var clientRemoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeElectricalConnection,
			[]model.FunctionType{
				model.FunctionTypeElectricalConnectionDescriptionListData,
				model.FunctionTypeElectricalConnectionParameterDescriptionListData,
				model.FunctionTypeElectricalConnectionPermittedValueSetListData,
			},
		},
		{
			model.FeatureTypeTypeMeasurement,
			[]model.FunctionType{
				model.FunctionTypeMeasurementDescriptionListData,
				model.FunctionTypeMeasurementListData,
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
