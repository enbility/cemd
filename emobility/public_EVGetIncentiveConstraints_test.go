package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVGetIncentiveConstraints(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	constraints := emobilty.EVGetIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)

	localDevice, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	constraints = emobilty.EVGetIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)

	emobilty.evIncentiveTable = incentiveTableConfiguration(localDevice, emobilty.evEntity)

	constraints = emobilty.EVGetIncentiveConstraints()
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)

	datagram := datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer, model.RoleTypeClient)

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

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	constraints = emobilty.EVGetIncentiveConstraints()
	assert.Equal(t, uint(1), constraints.MinSlots)
	assert.Equal(t, uint(10), constraints.MaxSlots)

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

	constraints = emobilty.EVGetIncentiveConstraints()
	assert.Equal(t, uint(1), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)

}
