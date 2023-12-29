package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVSoCSupported(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	data, err := emobilty.EVSoCSupported()
	assert.NotNil(t, err)
	assert.Equal(t, false, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVSoCSupported()
	assert.NotNil(t, err)
	assert.Equal(t, false, data)

	emobilty.evMeasurement = measurement(localEntity, emobilty.evEntity)

	data, err = emobilty.EVSoCSupported()
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
							UseCaseName:      util.Ptr(model.UseCaseNameTypeEVStateOfCharge),
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

	data, err = emobilty.EVSoCSupported()
	assert.NotNil(t, err)
	assert.Equal(t, false, data)

	datagram = datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer, model.RoleTypeClient)

	cmd = []model.CmdType{{
		MeasurementDescriptionListData: &model.MeasurementDescriptionListDataType{
			MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(0)),
					ScopeType:     util.Ptr(model.ScopeTypeTypeStateOfCharge),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVSoCSupported()
	assert.Nil(t, err)
	assert.Equal(t, true, data)
}
