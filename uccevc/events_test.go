package uccevc

import (
	"time"

	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCCEVCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = eebusutil.Ptr(model.TimeSeriesDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.TimeSeriesListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.IncentiveTableDescriptionDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.IncentiveDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCCEVCSuite) Test_evTimeSeriesDescriptionDataUpdate() {
	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.evEntity)

	timeDesc := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:   eebusutil.Ptr(model.TimeSeriesIdType(0)),
				TimeSeriesType: eebusutil.Ptr(model.TimeSeriesTypeTypeConstraints),
				UpdateRequired: eebusutil.Ptr(true),
			},
			{
				TimeSeriesId:   eebusutil.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesType: eebusutil.Ptr(model.TimeSeriesTypeTypeSingleDemand),
			},
		},
	}

	rTimeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	fErr := rTimeFeature.UpdateData(model.FunctionTypeTimeSeriesDescriptionListData, timeDesc, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.evEntity)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: eebusutil.Ptr(model.TimeSeriesIdType(1)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: eebusutil.Ptr(model.TimeSeriesSlotIdType(0)),
						MinValue:         model.NewScaledNumberType(1000),
						Value:            model.NewScaledNumberType(10000),
						MaxValue:         model.NewScaledNumberType(100000),
					},
				},
			},
		},
	}

	fErr = rTimeFeature.UpdateData(model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	demand, err := s.sut.EnergyDemand(s.evEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1000.0, demand.MinDemand)
	assert.Equal(s.T(), 10000.0, demand.OptDemand)
	assert.Equal(s.T(), 100000.0, demand.MaxDemand)
	assert.Equal(s.T(), 0.0, demand.DurationUntilStart)
	assert.Equal(s.T(), 0.0, demand.DurationUntilEnd)

	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.evEntity)

	constData := &model.TimeSeriesConstraintsListDataType{
		TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
			{
				TimeSeriesId:         eebusutil.Ptr(model.TimeSeriesIdType(0)),
				SlotCountMin:         eebusutil.Ptr(model.TimeSeriesSlotCountType(1)),
				SlotCountMax:         eebusutil.Ptr(model.TimeSeriesSlotCountType(10)),
				SlotDurationMin:      model.NewDurationType(1 * time.Minute),
				SlotDurationMax:      model.NewDurationType(60 * time.Minute),
				SlotDurationStepSize: model.NewDurationType(1 * time.Minute),
			},
		},
	}

	fErr = rTimeFeature.UpdateData(model.FunctionTypeTimeSeriesConstraintsListData, constData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.evEntity)

	incConstData := &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMin: eebusutil.Ptr(model.TimeSlotCountType(1)),
					SlotCountMax: eebusutil.Ptr(model.TimeSlotCountType(10)),
				},
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	fErr = rFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, incConstData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evTimeSeriesDescriptionDataUpdate(remoteSki, s.evEntity)

}
