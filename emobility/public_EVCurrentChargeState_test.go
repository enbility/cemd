package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_EVCurrentChargeState(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	data, err := emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeUnplugged, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeUnplugged, data)

	emobilty.evDeviceDiagnosis = deviceDiagnosis(localEntity, emobilty.evEntity)

	data, err = emobilty.EVCurrentChargeState()
	assert.NotNil(t, err)
	assert.Equal(t, EVChargeStateTypeUnknown, data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeActive, data)

	cmd = []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeStandby),
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypePaused, data)

	cmd = []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFailure),
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeError, data)

	cmd = []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeFinished),
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeFinished, data)

	cmd = []model.CmdType{{
		DeviceDiagnosisStateData: &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeInAlarm),
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVCurrentChargeState()
	assert.Nil(t, err)
	assert.Equal(t, EVChargeStateTypeUnknown, data)
}
