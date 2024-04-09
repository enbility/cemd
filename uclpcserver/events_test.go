package uclpcserver

import (
	"fmt"

	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
)

func (s *UCLPCServerSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Device = s.monitoredEntity.Device()
	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.CmdClassifier = eebusutil.Ptr(model.CmdClassifierTypeWrite)
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Function = model.FunctionTypeLoadControlLimitListData
	payload.Data = eebusutil.Ptr(model.LoadControlLimitListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.loadControlFeature
	s.sut.HandleEvent(payload)

	payload.Function = model.FunctionTypeDeviceConfigurationKeyValueListData
	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	payload.LocalFeature = s.deviceConfigurationFeature
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeBindingChange
	payload.ChangeType = spineapi.ElementChangeAdd
	payload.LocalFeature = s.loadControlFeature
	s.sut.HandleEvent(payload)
}

func (s *UCLPCServerSuite) Test_deviceConnected() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}

	s.sut.deviceConnected(payload)

	// no entities
	mockRemoteDevice := mocks.NewDeviceRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().Entities().Return(nil)
	payload.Device = mockRemoteDevice
	s.sut.deviceConnected(payload)

	// one entity with one DeviceDiagnosis server
	payload.Device = s.remoteDevice
	s.sut.deviceConnected(payload)

	s.sut.subscribeHeartbeatWorkaround(payload)
}

func (s *UCLPCServerSuite) Test_multipleDeviceDiagServer() {
	// multiple entities each with DeviceDiagnosis server

	payload := spineapi.EventPayload{
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}

	remoteDeviceName := "remote"

	var remoteFeatures = []struct {
		featureType   model.FeatureTypeType
		role          model.RoleType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeLoadControl,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceConfiguration,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeDeviceDiagnosisHeartbeatData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			model.RoleTypeClient,
			[]model.FunctionType{},
		},
	}
	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	// 4 entites
	for i := 1; i < 5; i++ {
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
						Entity:  []model.AddressEntityType{model.AddressEntityType(i)},
						Feature: eebusutil.Ptr(model.AddressFeatureType(index)),
					},
					FeatureType:       eebusutil.Ptr(feature.featureType),
					Role:              eebusutil.Ptr(feature.role),
					SupportedFunction: supportedFcts,
				},
			}
			featureInformations = append(featureInformations, featureInformation)
		}
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
					EntityType: eebusutil.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{2},
					},
					EntityType: eebusutil.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{3},
					},
					EntityType: eebusutil.Ptr(model.EntityTypeTypeCEM),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: eebusutil.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{4},
					},
					EntityType: eebusutil.Ptr(model.EntityTypeTypeCEM),
				},
			},
		},
		FeatureInformation: featureInformations,
	}

	_, err := s.remoteDevice.AddEntityAndFeatures(true, detailedData)
	if err != nil {
		fmt.Println(err)
	}
	s.remoteDevice.UpdateDevice(detailedData.DeviceInformation.Description)

	s.sut.deviceConnected(payload)

	s.sut.subscribeHeartbeatWorkaround(payload)
}
