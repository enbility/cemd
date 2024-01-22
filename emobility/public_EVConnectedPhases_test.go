package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_EVConnectedPhases(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	data, err := emobilty.EVConnectedPhases()
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVConnectedPhases()
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), data)

	emobilty.evElectricalConnection = electricalConnection(localEntity, emobilty.evEntity)

	data, err = emobilty.EVConnectedPhases()
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

	data, err = emobilty.EVConnectedPhases()
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

	data, err = emobilty.EVConnectedPhases()
	assert.Nil(t, err)
	assert.Equal(t, uint(1), data)
}
