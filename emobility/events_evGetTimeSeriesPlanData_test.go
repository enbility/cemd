package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_evGetTimeSeriesPlanData(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	data, err := emobilty.evGetTimeSeriesPlanData()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	localDevice, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	ctrl := gomock.NewController(t)

	dataProviderMock := NewMockEmobilityDataProvider(ctrl)
	emobilty.dataProvider = dataProviderMock

	data, err = emobilty.evGetTimeSeriesPlanData()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	emobilty.evTimeSeries = timeSeriesConfiguration(localDevice, emobilty.evEntity)

	data, err = emobilty.evGetTimeSeriesPlanData()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	datagram := datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{
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
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.evGetTimeSeriesPlanData()
	assert.NotNil(t, err)
	assert.Nil(t, data)

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
					TimePeriod: &model.TimePeriodType{
						StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					},
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
							Duration:         util.Ptr(model.DurationType("PT5M36S")),
							MaxValue:         model.NewScaledNumberType(4201),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P1D")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	data, err = emobilty.evGetTimeSeriesPlanData()
	assert.Nil(t, err)
	assert.NotNil(t, data)

}
