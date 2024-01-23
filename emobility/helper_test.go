package emobility

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/stretchr/testify/mock"
)

type WriteMessageHandler struct {
	sentMessages [][]byte

	mux sync.Mutex
}

var _ shipapi.ShipConnectionDataWriterInterface = (*WriteMessageHandler)(nil)

func (t *WriteMessageHandler) WriteShipMessageWithPayload(message []byte) {
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
func NewTestEMobility(service api.ServiceInterface, details *shipapi.ServiceDetails) *EMobility {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().Entity([]model.AddressEntityType{1})
	emobility := &EMobility{
		service: service,
		entity:  localEntity,
		ski:     ski,
	}

	service.RegisterRemoteSKI(ski, false)

	return emobility
}

func setupEmobility(t *testing.T) (*EMobility, api.ServiceInterface) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, 230.0, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(t)
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	eebusService := service.NewService(configuration, serviceHandler)
	_ = eebusService.Setup()
	details := shipapi.NewServiceDetails(remoteSki)
	emobility := NewTestEMobility(eebusService, details)
	return emobility, eebusService
}

func setupDevices(
	eebusService api.ServiceInterface) (
	spineapi.DeviceLocalInterface,
	spineapi.EntityLocalInterface,
	spineapi.DeviceRemoteInterface,
	[]spineapi.EntityRemoteInterface,
	*WriteMessageHandler) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(3, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(4, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(5, localEntity, model.FeatureTypeTypeIdentification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(6, localEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(6, localEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(6, localEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeClient)
	localEntity.AddFeature(f)

	writeHandler := &WriteMessageHandler{}
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

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

func datagramForEntityAndFeatures(
	notify bool,
	localDevice spineapi.DeviceLocalInterface,
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	featureType model.FeatureTypeType,
	remoteRole, localRole model.RoleType) model.DatagramType {
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

func featureOfTypeAndRole(
	entity spineapi.EntityRemoteInterface,
	featureType model.FeatureTypeType,
	role model.RoleType) spineapi.FeatureRemoteInterface {
	for _, f := range entity.Features() {
		if f.Type() == featureType && f.Role() == role {
			return f
		}
	}
	return nil
}

func deviceDiagnosis(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.DeviceDiagnosis {
	feature, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func electricalConnection(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.ElectricalConnection {
	feature, err := features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func measurement(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.Measurement {
	feature, err := features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func deviceConfiguration(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.DeviceConfiguration {
	feature, err := features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func identificationConfiguration(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.Identification {
	feature, err := features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func loadcontrol(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.LoadControl {
	feature, err := features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func timeSeriesConfiguration(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.TimeSeries {
	feature, err := features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}

func incentiveTableConfiguration(
	localEntity spineapi.EntityLocalInterface,
	entity spineapi.EntityRemoteInterface) *features.IncentiveTable {
	feature, err := features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		fmt.Println(err)
	}
	return feature
}
