package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_EVGetIncentiveConstraints(t *testing.T) {
	emobilty, eebusService := setupEmobility(t)

	mockRemoteDevice := mocks.NewDeviceRemoteInterface(t)
	mockRemoteEntity := mocks.NewEntityRemoteInterface(t)
	mockRemoteFeature := mocks.NewFeatureRemoteInterface(t)
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature)
	mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice)
	constraints, err := emobilty.EVIncentiveConstraints(mockRemoteEntity)
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.NotEqual(t, err, nil)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	constraints, err = emobilty.EVIncentiveConstraints(emobilty.evEntity)
	assert.Equal(t, uint(0), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.NotEqual(t, err, nil)

	constraints, err = emobilty.EVIncentiveConstraints(emobilty.evEntity)
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

	constraints, err = emobilty.EVIncentiveConstraints(emobilty.evEntity)
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

	constraints, err = emobilty.EVIncentiveConstraints(emobilty.evEntity)
	assert.Equal(t, uint(1), constraints.MinSlots)
	assert.Equal(t, uint(0), constraints.MaxSlots)
	assert.Equal(t, err, nil)

}
