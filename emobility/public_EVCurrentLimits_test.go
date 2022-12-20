package emobility

import (
	"testing"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func Test_EVCurrentLimits(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	minData, maxData, defaultData, err := emobilty.EVCurrentLimits()
	assert.NotNil(t, err)
	assert.Nil(t, minData)
	assert.Nil(t, maxData)
	assert.Nil(t, defaultData)

	localDevice, remoteDevice, entites, _ := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	minData, maxData, defaultData, err = emobilty.EVCurrentLimits()
	assert.NotNil(t, err)
	assert.Nil(t, minData)
	assert.Nil(t, maxData)
	assert.Nil(t, defaultData)

	emobilty.evElectricalConnection = electricalConnection(localDevice, emobilty.evEntity)

	minData, maxData, defaultData, err = emobilty.EVCurrentLimits()
	assert.NotNil(t, err)
	assert.Nil(t, minData)
	assert.Nil(t, maxData)
	assert.Nil(t, defaultData)

	datagram := datagramForEntityAndFeatures(false, localDevice, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{
			ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
					MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
					AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				},
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(1)),
					MeasurementId:          util.Ptr(model.MeasurementIdType(1)),
					AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				},
				{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(2)),
					MeasurementId:          util.Ptr(model.MeasurementIdType(2)),
					AcMeasuredPhases:       util.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				},
			},
		}}}

	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	minData, maxData, defaultData, err = emobilty.EVCurrentLimits()
	assert.NotNil(t, err)
	assert.Nil(t, minData)
	assert.Nil(t, maxData)
	assert.Nil(t, defaultData)

	type permittedStruct struct {
		defaultExists                      bool
		defaultValue, expectedDefaultValue float64
		minValue, expectedMinValue         float64
		maxValue, expectedMaxValue         float64
	}

	tests := []struct {
		name      string
		permitted []permittedStruct
	}{
		{
			"1 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2.2, 16, 16},
			},
		},
		{
			"1 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"1 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase ISO15118",
			[]permittedStruct{
				{true, 0.1, 0.1, 2, 2.2, 16, 16},
				{true, 0.1, 0.1, 2, 2.2, 16, 16},
				{true, 0.1, 0.1, 2, 2.2, 16, 16},
			},
		},
		{
			"3 Phase IEC61851",
			[]permittedStruct{
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
				{true, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]permittedStruct{
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
				{false, 0.0, 0.0, 6, 6, 16, 16},
			},
		},
	}

	cmd = []model.CmdType{{
		ElectricalConnectionPermittedValueSetListData: &model.ElectricalConnectionPermittedValueSetListDataType{
			ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
				{},
			},
		}}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			for index, data := range tc.permitted {
				item := model.ScaledNumberSetType{
					Range: []model.ScaledNumberRangeType{
						{
							Min: model.NewScaledNumberType(data.minValue),
							Max: model.NewScaledNumberType(data.maxValue),
						},
					},
				}
				if data.defaultExists {
					item.Value = []model.ScaledNumberType{*model.NewScaledNumberType(data.defaultValue)}
				}
				permittedData = append(permittedData, item)

				permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
					ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(index)),
					PermittedValueSet:      permittedData,
				}
				dataSet = append(dataSet, permittedItem)
			}

			cmd = []model.CmdType{{
				ElectricalConnectionPermittedValueSetListData: &model.ElectricalConnectionPermittedValueSetListDataType{
					ElectricalConnectionPermittedValueSetData: dataSet,
				}}}
			datagram.Payload.Cmd = cmd

			err = localDevice.ProcessCmd(datagram, remoteDevice)
			assert.Nil(t, err)

			minData, maxData, defaultData, err = emobilty.EVCurrentLimits()
			assert.Nil(t, err)

			assert.Nil(t, err)
			assert.Equal(t, len(tc.permitted), len(minData))
			assert.Equal(t, len(tc.permitted), len(maxData))
			assert.Equal(t, len(tc.permitted), len(defaultData))
			for index, item := range tc.permitted {
				assert.Equal(t, item.expectedMinValue, minData[index])
				assert.Equal(t, item.expectedMaxValue, maxData[index])
				assert.Equal(t, item.expectedDefaultValue, defaultData[index])
			}
		})
	}
}
