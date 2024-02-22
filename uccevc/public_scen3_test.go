package uccevc

import (
	"testing"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCCEVCSuite) Test_EVGetIncentiveConstraints() {
	constraints, err := s.sut.IncentiveConstraints(s.mockRemoteEntity)
	assert.Equal(s.T(), uint(0), constraints.MinSlots)
	assert.Equal(s.T(), uint(0), constraints.MaxSlots)
	assert.NotEqual(s.T(), err, nil)

	constraints, err = s.sut.IncentiveConstraints(s.evEntity)
	assert.Equal(s.T(), uint(0), constraints.MinSlots)
	assert.Equal(s.T(), uint(0), constraints.MaxSlots)
	assert.NotEqual(s.T(), err, nil)

	constData := &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
					SlotCountMax: util.Ptr(model.TimeSlotCountType(10)),
				},
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	constraints, err = s.sut.IncentiveConstraints(s.evEntity)
	assert.Equal(s.T(), uint(1), constraints.MinSlots)
	assert.Equal(s.T(), uint(10), constraints.MaxSlots)
	assert.Equal(s.T(), err, nil)

	constData = &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	constraints, err = s.sut.IncentiveConstraints(s.evEntity)
	assert.Equal(s.T(), uint(1), constraints.MinSlots)
	assert.Equal(s.T(), uint(0), constraints.MaxSlots)
	assert.Equal(s.T(), err, nil)
}

func (s *UCCEVCSuite) Test_EVWriteIncentives() {
	data := []api.DurationSlotValue{}

	err := s.sut.WriteIncentives(s.mockRemoteEntity, data)
	assert.NotNil(s.T(), err)

	err = s.sut.WriteIncentives(s.evEntity, data)
	assert.NotNil(s.T(), err)

	constData := &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMin: util.Ptr(model.TimeSlotCountType(1)),
					SlotCountMax: util.Ptr(model.TimeSlotCountType(10)),
				},
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	err = s.sut.WriteIncentives(s.evEntity, data)
	assert.Nil(s.T(), err)

	type dataStruct struct {
		error              bool
		minSlots, maxSlots uint
		slots              []api.DurationSlotValue
	}

	tests := []struct {
		name string
		data []dataStruct
	}{
		{
			"too few slots",
			[]dataStruct{
				{
					true, 2, 2,
					[]api.DurationSlotValue{
						{Duration: time.Hour, Value: 0.1},
					},
				},
			},
		}, {
			"too many slots",
			[]dataStruct{
				{
					true, 1, 1,
					[]api.DurationSlotValue{
						{Duration: time.Hour, Value: 0.1},
						{Duration: time.Hour, Value: 0.1},
					},
				},
			},
		},
		{
			"1 slot",
			[]dataStruct{
				{
					false, 1, 1,
					[]api.DurationSlotValue{
						{Duration: time.Hour, Value: 0.1},
					},
				},
			},
		},
		{
			"2 slots",
			[]dataStruct{
				{
					false, 1, 2,
					[]api.DurationSlotValue{
						{Duration: time.Hour, Value: 0.1},
						{Duration: 30 * time.Minute, Value: 0.2},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			for _, data := range tc.data {

				constData = &model.IncentiveTableConstraintsDataType{
					IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
						{
							IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
								SlotCountMin: util.Ptr(model.TimeSlotCountType(data.minSlots)),
								SlotCountMax: util.Ptr(model.TimeSlotCountType(data.maxSlots)),
							},
						},
					},
				}

				fErr := rFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, constData, nil, nil)
				assert.Nil(s.T(), fErr)

				err = s.sut.WriteIncentives(s.evEntity, data.slots)
				if data.error {
					assert.NotNil(t, err)
					continue
				} else {
					assert.Nil(t, err)
				}

			}
		})
	}
}
