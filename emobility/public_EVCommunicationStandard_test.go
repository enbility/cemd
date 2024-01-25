package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EVCommunicationStandard(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	mockRemoteDevice := mocks.NewDeviceRemoteInterface(t)
	mockRemoteEntity := mocks.NewEntityRemoteInterface(t)
	mockRemoteFeature := mocks.NewFeatureRemoteInterface(t)
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature)
	mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice)
	data, err := emobilty.EVCommunicationStandard(mockRemoteEntity)
	assert.NotNil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeUnknown, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVCommunicationStandard(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeUnknown, data)

	data, err = emobilty.EVCommunicationStandard(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeUnknown, data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		DeviceConfigurationKeyValueDescriptionListData: &model.DeviceConfigurationKeyValueDescriptionListDataType{
			DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
				{
					KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
					KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCommunicationStandard(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeUnknown, data)

	cmd = []model.CmdType{{
		DeviceConfigurationKeyValueDescriptionListData: &model.DeviceConfigurationKeyValueDescriptionListDataType{
			DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
				{
					KeyId:   util.Ptr(model.DeviceConfigurationKeyIdType(0)),
					KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCommunicationStandard(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeUnknown, data)

	cmd = []model.CmdType{{
		DeviceConfigurationKeyValueListData: &model.DeviceConfigurationKeyValueListDataType{
			DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
				{
					KeyId: util.Ptr(model.DeviceConfigurationKeyIdType(0)),
					Value: &model.DeviceConfigurationKeyValueValueType{
						String: util.Ptr(model.DeviceConfigurationKeyValueStringType(EVCommunicationStandardTypeISO151182ED1)),
					},
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCommunicationStandard(emobilty.evEntity)
	assert.Nil(t, err)
	assert.Equal(t, EVCommunicationStandardTypeISO151182ED1, data)
}
