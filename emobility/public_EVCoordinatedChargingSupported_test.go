package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_EVCoordinatedChargingSupported(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	data, err := emobilty.EVCoordinatedChargingSupported()
	assert.NotNil(t, err)
	assert.Equal(t, false, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVCoordinatedChargingSupported()
	assert.Nil(t, err)
	assert.Equal(t, false, data)

	datagram := datagramForEntityAndFeatures(true, localDevice, localEntity, nil, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial, model.RoleTypeSpecial)

	cmd := []model.CmdType{{
		NodeManagementUseCaseData: &model.NodeManagementUseCaseDataType{
			UseCaseInformation: []model.UseCaseInformationDataType{
				{
					Actor: util.Ptr(model.UseCaseActorTypeEV),
					UseCaseSupport: []model.UseCaseSupportType{
						{
							UseCaseName:      util.Ptr(model.UseCaseNameTypeCoordinatedEVCharging),
							UseCaseAvailable: util.Ptr(true),
							ScenarioSupport:  []model.UseCaseScenarioSupportType{1},
						},
					},
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCoordinatedChargingSupported()
	assert.Nil(t, err)
	assert.Equal(t, true, data)
}
