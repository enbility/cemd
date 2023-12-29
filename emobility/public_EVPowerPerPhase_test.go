package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVPowerPerPhase_Power(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	data, err := emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	emobilty.evElectricalConnection = electricalConnection(localEntity, emobilty.evEntity)
	emobilty.evMeasurement = measurement(localEntity, emobilty.evEntity)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{
			ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
					MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
					ScopeType:              util.Ptr(model.ScopeTypeTypeACPower),
					AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	datagram = datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer, model.RoleTypeClient)

	cmd = []model.CmdType{{
		MeasurementDescriptionListData: &model.MeasurementDescriptionListDataType{
			MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
				{
					MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
					MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
					CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
					ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	cmd = []model.CmdType{{
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(0)),
					Value:         model.NewScaledNumberType(80),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.Nil(t, err)
	assert.Equal(t, 80.0, data[0])
}

func Test_EVPowerPerPhase_Current(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	data, err := emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	emobilty.evElectricalConnection = electricalConnection(localEntity, emobilty.evEntity)
	emobilty.evMeasurement = measurement(localEntity, emobilty.evEntity)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{
			ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
					MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
					ScopeType:              util.Ptr(model.ScopeTypeTypeACCurrent),
					AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	datagram = datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer, model.RoleTypeClient)

	cmd = []model.CmdType{{
		MeasurementDescriptionListData: &model.MeasurementDescriptionListDataType{
			MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
				{
					MeasurementId:   util.Ptr(model.MeasurementIdType(0)),
					MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
					CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
					ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	cmd = []model.CmdType{{
		MeasurementListData: &model.MeasurementListDataType{
			MeasurementData: []model.MeasurementDataType{
				{
					MeasurementId: util.Ptr(model.MeasurementIdType(0)),
					Value:         model.NewScaledNumberType(10),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.EVPowerPerPhase()
	assert.Nil(t, err)
	assert.Equal(t, 2300.0, data[0])
}
