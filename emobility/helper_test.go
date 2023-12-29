package emobility

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type WriteMessageHandler struct {
	sentMessages [][]byte

	mux sync.Mutex
}

var _ spine.SpineDataConnection = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteSpineMessage(message []byte) {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.sentMessages = append(t.sentMessages, message)
}

func (t *WriteMessageHandler) LastMessage() []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	if len(t.sentMessages) == 0 {
		return nil
	}

	return t.sentMessages[len(t.sentMessages)-1]
}

func (t *WriteMessageHandler) MessageWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData != nil {
			continue
		}

		return msg
	}

	return nil
}

func (t *WriteMessageHandler) ResultWithReference(msgCounterReference *model.MsgCounterType) []byte {
	t.mux.Lock()
	defer t.mux.Unlock()

	var datagram model.Datagram

	for _, msg := range t.sentMessages {
		if err := json.Unmarshal(msg, &datagram); err != nil {
			return nil
		}
		if datagram.Datagram.Header.MsgCounterReference == nil {
			continue
		}
		if uint(*datagram.Datagram.Header.MsgCounterReference) != uint(*msgCounterReference) {
			continue
		}
		if datagram.Datagram.Payload.Cmd[0].ResultData == nil {
			continue
		}

		return msg
	}

	return nil
}

const remoteSki string = "testremoteski"

// we don't want to handle events in these tests for now, so we don't use NewEMobility(...)
func NewTestEMobility(service *service.EEBUSService, details *service.ServiceDetails) *EMobilityImpl {
	ski := util.NormalizeSKI(details.SKI)

	localEntity := service.LocalDevice().Entity([]model.AddressEntityType{1})
	emobility := &EMobilityImpl{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}

	service.RegisterRemoteSKI(ski, true)

	return emobility
}

func setupEmobility() (*EMobilityImpl, *service.EEBUSService) {
	cert, _ := service.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := service.NewConfiguration(
		"test", "test", "test", "test",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, 230.0, time.Second*4)
	eebusService := service.NewEEBUSService(configuration, nil)
	_ = eebusService.Setup()
	details := service.NewServiceDetails(remoteSki)
	emobility := NewTestEMobility(eebusService, details)
	return emobility, eebusService
}

func setupDevices(eebusService *service.EEBUSService) (*spine.DeviceLocalImpl, *spine.EntityLocalImpl, *spine.DeviceRemoteImpl, []*spine.EntityRemoteImpl, *WriteMessageHandler) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(2, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(3, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(4, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(5, localEntity, model.FeatureTypeTypeIdentification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(6, localEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(6, localEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(6, localEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeClient)
	localEntity.AddFeature(f)

	writeHandler := &WriteMessageHandler{}
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemoteImpl(localDevice, remoteSki, sender)

	var clientRemoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{
			model.FeatureTypeTypeDeviceDiagnosis,
			[]model.FunctionType{},
		},
		{
			model.FeatureTypeTypeDeviceConfiguration,
			[]model.FunctionType{
				model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
				model.FunctionTypeDeviceConfigurationKeyValueListData,
			},
		},
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
		{
			model.FeatureTypeTypeLoadControl,
			[]model.FunctionType{
				model.FunctionTypeLoadControlLimitDescriptionListData,
				model.FunctionTypeLoadControlLimitListData,
			},
		},
		{
			model.FeatureTypeTypeIdentification,
			[]model.FunctionType{
				model.FunctionTypeIdentificationListData,
			},
		},
		{model.FeatureTypeTypeTimeSeries,
			[]model.FunctionType{
				model.FunctionTypeTimeSeriesDescriptionListData,
				model.FunctionTypeTimeSeriesListData,
				model.FunctionTypeTimeSeriesConstraintsListData,
			},
		},
		{model.FeatureTypeTypeIncentiveTable,
			[]model.FunctionType{
				model.FunctionTypeIncentiveTableConstraintsData,
			},
		},
	}

	remoteDeviceName := "remote"

	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for index, feature := range clientRemoteFeatures {
		supportedFcts := []model.FunctionPropertyType{}
		for _, fct := range feature.supportedFcts {
			supportedFct := model.FunctionPropertyType{
				Function: util.Ptr(fct),
				PossibleOperations: &model.PossibleOperationsType{
					Read: &model.PossibleOperationsReadType{},
				},
			}
			supportedFcts = append(supportedFcts, supportedFct)
		}

		featureInformation := model.NodeManagementDetailedDiscoveryFeatureInformationType{
			Description: &model.NetworkManagementFeatureDescriptionDataType{
				FeatureAddress: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
					Entity:  []model.AddressEntityType{1, 1},
					Feature: util.Ptr(model.AddressFeatureType(index)),
				},
				FeatureType:       util.Ptr(feature.featureType),
				Role:              util.Ptr(model.RoleTypeServer),
				SupportedFunction: supportedFcts,
			},
		}
		featureInformations = append(featureInformations, featureInformation)
	}

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEVSE),
				},
			},
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1, 1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEV),
				},
			},
		},
		FeatureInformation: featureInformations,
	}
	localDevice.AddRemoteDeviceForSki(remoteSki, remoteDevice)

	entities, err := remoteDevice.AddEntityAndFeatures(true, detailedData)
	if err != nil {
		fmt.Println(err)
	}

	return localDevice, localEntity, remoteDevice, entities, writeHandler
}

func datagramForEntityAndFeatures(notify bool, localDevice *spine.DeviceLocalImpl, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl, featureType model.FeatureTypeType, remoteRole, localRole model.RoleType) model.DatagramType {
	var addressSource, addressDestination *model.FeatureAddressType
	if remoteEntity == nil {
		// NodeManagement
		addressSource = &model.FeatureAddressType{
			Entity:  []model.AddressEntityType{0},
			Feature: util.Ptr(model.AddressFeatureType(0)),
		}
		addressDestination = &model.FeatureAddressType{
			Device:  localDevice.Address(),
			Entity:  []model.AddressEntityType{0},
			Feature: util.Ptr(model.AddressFeatureType(0)),
		}
	} else {
		rFeature := featureOfTypeAndRole(remoteEntity, featureType, remoteRole)
		addressSource = rFeature.Address()

		lFeature := localEntity.FeatureOfTypeAndRole(featureType, localRole)
		addressDestination = lFeature.Address()
	}
	datagram := model.DatagramType{
		Header: model.HeaderType{
			AddressSource:       addressSource,
			AddressDestination:  addressDestination,
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: util.Ptr(model.MsgCounterType(1)),
			CmdClassifier:       util.Ptr(model.CmdClassifierTypeReply),
		},
		Payload: model.PayloadType{
			Cmd: []model.CmdType{},
		},
	}
	if notify {
		datagram.Header.CmdClassifier = util.Ptr(model.CmdClassifierTypeNotify)
	}

	return datagram
}

func featureOfTypeAndRole(entity *spine.EntityRemoteImpl, featureType model.FeatureTypeType, role model.RoleType) *spine.FeatureRemoteImpl {
	for _, f := range entity.Features() {
		if f.Type() == featureType && f.Role() == role {
			return f
		}
	}
	return nil
}

func deviceDiagnosis(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.DeviceDiagnosis {
	feature, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func electricalConnection(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.ElectricalConnection {
	feature, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func measurement(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.Measurement {
	feature, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func deviceConfiguration(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.DeviceConfiguration {
	feature, err := features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func identificationConfiguration(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.Identification {
	feature, err := features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func loadcontrol(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.LoadControl {
	feature, err := features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func timeSeriesConfiguration(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.TimeSeries {
	feature, err := features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func incentiveTableConfiguration(localEntity *spine.EntityLocalImpl, entity *spine.EntityRemoteImpl) *features.IncentiveTable {
	feature, err := features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}
