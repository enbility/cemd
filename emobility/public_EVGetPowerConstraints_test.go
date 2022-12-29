package emobility

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVGetPowerConstraints(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	minSlots, maxSlots, minDuration, maxDuration, durationStepSize := emobilty.EVGetPowerConstraints()
	assert.Equal(t, uint(0), minSlots)
	assert.Equal(t, uint(0), maxSlots)
	assert.Equal(t, time.Duration(0), minDuration)
	assert.Equal(t, time.Duration(0), maxDuration)
	assert.Equal(t, time.Duration(0), durationStepSize)

	localDevice, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	minSlots, maxSlots, minDuration, maxDuration, durationStepSize = emobilty.EVGetPowerConstraints()
	assert.Equal(t, uint(0), minSlots)
	assert.Equal(t, uint(0), maxSlots)
	assert.Equal(t, time.Duration(0), minDuration)
	assert.Equal(t, time.Duration(0), maxDuration)
	assert.Equal(t, time.Duration(0), durationStepSize)

	emobilty.evTimeSeries = timeSeriesConfiguration(localDevice, emobilty.evEntity)

	minSlots, maxSlots, minDuration, maxDuration, durationStepSize = emobilty.EVGetPowerConstraints()
	assert.Equal(t, uint(0), minSlots)
	assert.Equal(t, uint(0), maxSlots)
	assert.Equal(t, time.Duration(0), minDuration)
	assert.Equal(t, time.Duration(0), maxDuration)
	assert.Equal(t, time.Duration(0), durationStepSize)

	datagram := datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		TimeSeriesConstraintsListData: &model.TimeSeriesConstraintsListDataType{
			TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
				{
					TimeSeriesId:         util.Ptr(model.TimeSeriesIdType(0)),
					SlotCountMin:         util.Ptr(model.TimeSeriesSlotCountType(1)),
					SlotCountMax:         util.Ptr(model.TimeSeriesSlotCountType(10)),
					SlotDurationMin:      model.NewDurationType(1 * time.Minute),
					SlotDurationMax:      model.NewDurationType(60 * time.Minute),
					SlotDurationStepSize: model.NewDurationType(1 * time.Minute),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	minSlots, maxSlots, minDuration, maxDuration, durationStepSize = emobilty.EVGetPowerConstraints()
	assert.Equal(t, uint(1), minSlots)
	assert.Equal(t, uint(10), maxSlots)
	assert.Equal(t, time.Duration(1*time.Minute), minDuration)
	assert.Equal(t, time.Duration(1*time.Hour), maxDuration)
	assert.Equal(t, time.Duration(1*time.Minute), durationStepSize)
}
