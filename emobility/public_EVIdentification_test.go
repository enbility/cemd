package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_EVIdentification(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	mockRemoteEntity := mocks.NewEntityRemoteInterface(t)
	data, err := emobilty.EVIdentification(mockRemoteEntity)
	assert.NotNil(t, err)
	assert.Equal(t, "", data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVIdentification(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, "", data)

	data, err = emobilty.EVIdentification(emobilty.evEntity)
	assert.NotNil(t, err)
	assert.Equal(t, "", data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeIdentification, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		IdentificationListData: &model.IdentificationListDataType{
			IdentificationData: []model.IdentificationDataType{
				{
					IdentificationId:    util.Ptr(model.IdentificationIdType(0)),
					IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
					IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVIdentification(emobilty.evEntity)
	assert.Nil(t, err)
	assert.Equal(t, "test", data)
}
