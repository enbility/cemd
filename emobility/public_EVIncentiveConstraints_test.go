package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVGetIncentiveConstraints(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	constraints, err := emobilty.EVIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.NotEqual(t, err, nil)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	constraints, err = emobilty.EVIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.NotEqual(t, err, nil)

	emobilty.evIncentiveTable = incentiveTableConfiguration(localEntity, emobilty.evEntity)

	constraints, err = emobilty.EVIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.NotEqual(t, err, nil)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		IncentiveTableConstraintsData: &model.IncentiveTableConstraintsDataType{
			IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
				{
					IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
						SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
						SlotCountMax: util.Ptr(model.TimeSlotCountType(10)),
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	constraints, err = emobilty.EVIncentiveConstraints()
	assert.Equal(t, uint(1), constraints.MinSlots)
	assert.Equal(t, uint(10), constraints.MaxSlots)
	assert.Equal(t, err, nil)

	cmd = []model.CmdType{{
		IncentiveTableConstraintsData: &model.IncentiveTableConstraintsDataType{
			IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
				{
					IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
						SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	constraints, err = emobilty.EVIncentiveConstraints()
	assert.Equal(t, uint(1), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.Equal(t, err, nil)

}
