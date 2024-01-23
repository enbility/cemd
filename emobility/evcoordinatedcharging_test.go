package emobility

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func Test_CoordinatedChargingScenarios(t *testing.T) {
	emobility, eebusService := setupEmobility(t)

	data, err := emobility.EVChargedEnergy()
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, data)

	localDevice, localEntity, remoteDevice, entites, _ := setupDevices(eebusService)
	emobility.evseEntity = entites[0]
	emobility.evEntity = entites[1]

	ctrl := gomock.NewController(t)

	dataProviderMock := NewMockEmobilityDataProvider(ctrl)
	emobility.dataProvider = dataProviderMock

	emobility.evTimeSeries = timeSeriesConfiguration(localEntity, emobility.evEntity)
	emobility.evIncentiveTable = incentiveTableConfiguration(localEntity, emobility.evEntity)

	datagramtt := datagramForEntityAndFeatures(false, localDevice, localEntity, emobility.evEntity, model.FeatureTypeTypeTimeSeries, model.RoleTypeServer, model.RoleTypeClient)
	datagramit := datagramForEntityAndFeatures(false, localDevice, localEntity, emobility.evEntity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer, model.RoleTypeClient)

	setupTimeSeries(t, datagramtt, localDevice, remoteDevice)
	setupIncentiveTable(t, datagramit, localDevice, remoteDevice)

	// demand, No Profile No Timer demand

	cmd := []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
					TimePeriod: &model.TimePeriodType{
						StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					},
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Value:            model.NewScaledNumberType(0),
							MaxValue:         model.NewScaledNumberType(74690),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)

	demand, err := emobility.EVEnergyDemand()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, demand.MinDemand)
	assert.Equal(t, 0.0, demand.OptDemand)
	assert.Equal(t, 74690.0, demand.MaxDemand)
	assert.Equal(t, 0.0, demand.DurationUntilStart)
	assert.Equal(t, 0.0, demand.DurationUntilEnd)

	// the final plan

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
							Duration:         util.Ptr(model.DurationType("PT18H3M7S")),
							MaxValue:         model.NewScaledNumberType(4163),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("PT42M")),
							MaxValue:         model.NewScaledNumberType(2736),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P1D")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)

	// demand, profile + timer with 80% target and no climate, minSoC reached

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
					TimePeriod: &model.TimePeriodType{
						StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					},
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P2DT4H40M36S")),
							Value:            model.NewScaledNumberType(53400),
							MaxValue:         model.NewScaledNumberType(74690),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P1D")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)

	demand, err = emobility.EVEnergyDemand()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, demand.MinDemand)
	assert.Equal(t, 53400.0, demand.OptDemand)
	assert.Equal(t, 74690.0, demand.MaxDemand)
	assert.Equal(t, 0.0, demand.DurationUntilStart)
	assert.Equal(t, time.Duration(time.Hour*52+time.Minute*40+time.Second*36).Seconds(), demand.DurationUntilEnd)

	// the final plan

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
							Duration:         util.Ptr(model.DurationType("P1DT15H24M24S")),
							MaxValue:         model.NewScaledNumberType(0),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("PT12H35M50S")),
							MaxValue:         model.NewScaledNumberType(4163),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(2)),
							Duration:         util.Ptr(model.DurationType("PT40M22S")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)

	// demand, profile with 25% min SoC, minSoC not reached, no timer

	cmd = []model.CmdType{{
		TimeSeriesListData: &model.TimeSeriesListDataType{
			TimeSeriesData: []model.TimeSeriesDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
				},
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(2)),
					TimePeriod: &model.TimePeriodType{
						StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					},
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(0)),
							Duration:         util.Ptr(model.DurationType("PT8M42S")),
							MaxValue:         model.NewScaledNumberType(4212),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P1D")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(3)),
					TimePeriod: &model.TimePeriodType{
						StartTime: model.NewAbsoluteOrRelativeTimeType("PT0S"),
					},
					TimeSeriesSlot: []model.TimeSeriesSlotType{
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Value:            model.NewScaledNumberType(600),
							MinValue:         model.NewScaledNumberType(600),
							MaxValue:         model.NewScaledNumberType(75600),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)

	demand, err = emobility.EVEnergyDemand()
	assert.Nil(t, err)
	assert.Equal(t, 600.0, demand.MinDemand)
	assert.Equal(t, 600.0, demand.OptDemand)
	assert.Equal(t, 75600.0, demand.MaxDemand)
	assert.Equal(t, 0.0, demand.DurationUntilStart)
	assert.Equal(t, 0.0, demand.DurationUntilEnd)

	// the final plan

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
							Duration:         util.Ptr(model.DurationType("PT8M42S")),
							MaxValue:         model.NewScaledNumberType(4212),
						},
						{
							TimeSeriesSlotId: util.Ptr(model.TimeSeriesSlotIdType(1)),
							Duration:         util.Ptr(model.DurationType("P1D")),
							MaxValue:         model.NewScaledNumberType(0),
						},
					},
				},
			},
		},
	}}

	datagramtt.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagramtt, remoteDevice)
	assert.Nil(t, err)
}

func setupTimeSeries(
	t *testing.T,
	datagram model.DatagramType,
	localDevice api.DeviceLocalInterface,
	remoteDevice api.DeviceRemoteInterface) {
	cmd := []model.CmdType{{
		TimeSeriesConstraintsListData: &model.TimeSeriesConstraintsListDataType{
			TimeSeriesConstraintsData: []model.TimeSeriesConstraintsDataType{
				{
					TimeSeriesId: util.Ptr(model.TimeSeriesIdType(1)),
					SlotCountMax: util.Ptr(model.TimeSeriesSlotCountType(30)),
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	cmd = []model.CmdType{{
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
		},
	}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}

func setupIncentiveTable(
	t *testing.T,
	datagram model.DatagramType,
	localDevice api.DeviceLocalInterface,
	remoteDevice api.DeviceRemoteInterface) {
	cmd := []model.CmdType{{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: []model.IncentiveTableDescriptionType{
				{
					TariffDescription: &model.TariffDescriptionDataType{
						TariffId:        util.Ptr(model.TariffIdType(1)),
						TariffWriteable: util.Ptr(true),
						UpdateRequired:  util.Ptr(false),
						ScopeType:       util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
					},
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}

/*
func requestIncentiveUpdate(t *testing.T, datagram model.DatagramType, localDevice api.DeviceLocal, remoteDevice api.DeviceRemote) {
	cmd := []model.CmdType{{
		IncentiveTableDescriptionData: &model.IncentiveTableDescriptionDataType{
			IncentiveTableDescription: []model.IncentiveTableDescriptionType{
				{
					TariffDescription: &model.TariffDescriptionDataType{
						TariffId:        util.Ptr(model.TariffIdType(1)),
						TariffWriteable: util.Ptr(true),
						UpdateRequired:  util.Ptr(true),
						ScopeType:       util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
					},
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}

func requestPowerTableUpdate(t *testing.T, datagram model.DatagramType, localDevice api.DeviceLocal, remoteDevice api.DeviceRemote) {
	cmd := []model.CmdType{{
		TimeSeriesDescriptionListData: &model.TimeSeriesDescriptionListDataType{
			TimeSeriesDescriptionData: []model.TimeSeriesDescriptionDataType{
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(1)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
					TimeSeriesWriteable: util.Ptr(true),
					UpdateRequired:      util.Ptr(true),
				},
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(2)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypePlan),
					TimeSeriesWriteable: util.Ptr(false),
					Unit:                util.Ptr(model.UnitOfMeasurementTypeW),
				},
				{
					TimeSeriesId:        util.Ptr(model.TimeSeriesIdType(3)),
					TimeSeriesType:      util.Ptr(model.TimeSeriesTypeTypeConstraints),
					TimeSeriesWriteable: util.Ptr(false),
					Unit:                util.Ptr(model.UnitOfMeasurementTypeWh),
				},
			},
		},
	}}

	datagram.Payload.Cmd = cmd

	err := localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)
}
*/
