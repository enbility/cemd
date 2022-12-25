package emobility

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVEnergySingleDemand(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	data, duration, err := emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	localDevice, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	emobilty.evDeviceConfiguration = deviceConfiguration(localDevice, emobilty.evEntity)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	datagram := datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
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

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

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

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	emobilty.evTimeSeries = timeSeriesConfiguration(localDevice, emobilty.evEntity)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	datagram = datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer, model.RoleTypeClient)

	cmd = []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{
			TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
				{
					TimeSeriesId:   util.Ptr(model.TimeSeriesIdType(0)),
					TimeSeriesType: util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
						},
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)
	assert.Equal(t, time.Duration(0), duration)

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
							Value:            model.NewScaledNumberType(10000),
						},
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.Nil(t, err)
	assert.Equal(t, 10000.0, data)
	assert.Equal(t, time.Duration(0), duration)

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(0)),
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
							Value:            model.NewScaledNumberType(10000),
							Duration:         model.NewDurationType(2 * time.Hour),
						},
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, duration, err = emobilty.EVEnergyDemand()
	assert.Nil(t, err)
	assert.Equal(t, 10000.0, data)
	assert.Equal(t, time.Duration(2*time.Hour), duration)
}
