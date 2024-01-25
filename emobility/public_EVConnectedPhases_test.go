package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EVConnectedPhases(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	mockRemoteDevice := mocks.NewDeviceRemoteInterface(t)
	mockRemoteEntity := mocks.NewEntityRemoteInterface(t)
	mockRemoteFeature := mocks.NewFeatureRemoteInterface(t)
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature)
	mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice)
	data, err := emobilty.EVConnectedPhases(mockRemoteEntity)
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVConnectedPhases(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), data)

	data, err = emobilty.EVConnectedPhases(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		ElectricalConnectionDescriptionListData: &model.ElectricalConnectionDescriptionListDataType{
			ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVConnectedPhases(emobilty.evEntity)
	assert.Nil(t, err)
	assert.Equal(t, uint(3), data)

	cmd = []model.CmdType{{
		ElectricalConnectionDescriptionListData: &model.ElectricalConnectionDescriptionListDataType{
			ElectricalConnectionDescriptionData: []model.ElectricalConnectionDescriptionDataType{
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					AcConnectedPhases:      util.Ptr(uint(1)),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVConnectedPhases(emobilty.evEntity)
	assert.Nil(t, err)
	assert.Equal(t, uint(1), data)
}
