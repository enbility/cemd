package uccevc

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/ship-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCCEVCSuite) Test_EVChargePlan() {
	_, err := s.sut.ChargePlan(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.NotNil(s.T(), err)

	descData := &model.TimeSeriesDescriptionListDataType{
		TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(1)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
				TimeSeriesWriteable: util.Ptr(true),
				UpdateRequired:      util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(2)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypePlan),
				TimeSeriesWriteable: util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
			},
			{
				TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(3)),
				TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeSingleDemand),
				TimeSeriesWriteable: util.Ptr(false),
				Unit:                util.Ptr(model.UnitOfMeasurementTypeWh),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeTimeSeriesDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.NotNil(s.T(), err)

	timeData := &model.TimeSeriesListDataType{
		TimeSeriesData: []model.TimeSeriesDataType{
			{
				TimeSeriesId: eebusutil.Ptr(model.TimeSeriesIdType(2)),
				TimePeriod: &model.TimePeriodType{
					StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
				},
				TimeSeriesSlot: []model.TimeSeriesSlotType{
					{
						TimeSeriesSlotId: eebusutil.Ptr(model.TimeSeriesSlotIdType(0)),
						Duration:         eebusutil.Ptr(model.DurationType("PT5M36S")),
						MaxValue:         model.NewScaledNumberType(4201),
					},
					{
						TimeSeriesSlotId: eebusutil.Ptr(model.TimeSeriesSlotIdType(1)),
						Duration:         eebusutil.Ptr(model.DurationType("P1D")),
						MaxValue:         model.NewScaledNumberType(0),
					},
				},
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeTimeSeriesListData, timeData, nil, nil)
	assert.Nil(s.T(), fErr)

	_, err = s.sut.ChargePlan(s.evEntity)
	assert.Nil(s.T(), err)
}
