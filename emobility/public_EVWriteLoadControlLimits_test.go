package emobility

import (
	"encoding/json"
	"testing"

	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/spine/model"
	eebusUtil "github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func Test_EVWriteLoadControlLimits(t *testing.T) {
	emobilty, eebusService := setupEmobility()

	loadLimits := []EVLoadLimits{}

	err := emobilty.EVWriteLoadControlLimits(loadLimits)
	assert.NotNil(t, err)

	localDevice, localEntity, remoteDevice, entites, writeHandler := setupDevices(eebusService)
	emobilty.evseEntity = entites[0]
	emobilty.evEntity = entites[1]

	err = emobilty.EVWriteLoadControlLimits(loadLimits)
	assert.NotNil(t, err)

	emobilty.evElectricalConnection = electricalConnection(localEntity, emobilty.evEntity)
	emobilty.evLoadControl = loadcontrol(localEntity, emobilty.evEntity)

	err = emobilty.EVWriteLoadControlLimits(loadLimits)
	assert.NotNil(t, err)

	datagram := datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

	cmd := []model.CmdType{{
		ElectricalConnectionParameterDescriptionListData: &model.ElectricalConnectionParameterDescriptionListDataType{
			ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
				{
					ElectricalConnectionId: eebusUtil.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            eebusUtil.Ptr(model.ElectricalConnectionParameterIdType(0)),
					MeasurementId:          eebusUtil.Ptr(model.MeasurementIdType(0)),
					AcMeasuredPhases:       eebusUtil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
				},
				{
					ElectricalConnectionId: eebusUtil.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            eebusUtil.Ptr(model.ElectricalConnectionParameterIdType(1)),
					MeasurementId:          eebusUtil.Ptr(model.MeasurementIdType(1)),
					AcMeasuredPhases:       eebusUtil.Ptr(model.ElectricalConnectionPhaseNameTypeB),
				},
				{
					ElectricalConnectionId: eebusUtil.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            eebusUtil.Ptr(model.ElectricalConnectionParameterIdType(2)),
					MeasurementId:          eebusUtil.Ptr(model.MeasurementIdType(2)),
					AcMeasuredPhases:       eebusUtil.Ptr(model.ElectricalConnectionPhaseNameTypeC),
				},
			},
		}}}
	datagram.Payload.Cmd = cmd

	err = localDevice.ProcessCmd(datagram, remoteDevice)
	assert.Nil(t, err)

	err = emobilty.EVWriteLoadControlLimits(loadLimits)
	assert.NotNil(t, err)

	type dataStruct struct {
		phases                                   int
		permittedDefaultExists                   bool
		permittedDefaultValue                    float64
		permittedMinValue                        float64
		permittedMaxValue                        float64
		obligations, obligationsExpected         []float64
		recommendations, recommendationsExpected []float64
	}

	tests := []struct {
		name string
		data []dataStruct
	}{
		{
			"1 Phase ISO15118",
			[]dataStruct{
				{1, true, 0.1, 2, 16, []float64{0}, []float64{0.1}, []float64{}, []float64{}},
				{1, true, 0.1, 2, 16, []float64{2.2}, []float64{2.2}, []float64{}, []float64{}},
				{1, true, 0.1, 2, 16, []float64{10}, []float64{10}, []float64{}, []float64{}},
				{1, true, 0.1, 2, 16, []float64{16}, []float64{16}, []float64{}, []float64{}},
			},
		},
		{
			"3 Phase ISO15118",
			[]dataStruct{
				{3, true, 0.1, 2, 16, []float64{0, 0, 0}, []float64{0.1, 0.1, 0.1}, []float64{}, []float64{}},
				{3, true, 0.1, 2, 16, []float64{2.2, 2.2, 2.2}, []float64{2.2, 2.2, 2.2}, []float64{}, []float64{}},
				{3, true, 0.1, 2, 16, []float64{10, 10, 10}, []float64{10, 10, 10}, []float64{}, []float64{}},
				{3, true, 0.1, 2, 16, []float64{16, 16, 16}, []float64{16, 16, 16}, []float64{}, []float64{}},
			},
		},
		{
			"1 Phase IEC61851",
			[]dataStruct{
				{1, true, 0, 6, 16, []float64{0}, []float64{0}, []float64{}, []float64{}},
				{1, true, 0, 6, 16, []float64{6}, []float64{6}, []float64{}, []float64{}},
				{1, true, 0, 6, 16, []float64{10}, []float64{10}, []float64{}, []float64{}},
				{1, true, 0, 6, 16, []float64{16}, []float64{16}, []float64{}, []float64{}},
			},
		},
		{
			"3 Phase IEC61851",
			[]dataStruct{
				{3, true, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{}, []float64{}},
				{3, true, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}, []float64{}, []float64{}},
				{3, true, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}, []float64{}, []float64{}},
				{3, true, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}, []float64{}, []float64{}},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]dataStruct{
				{3, false, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}, []float64{}, []float64{}},
				{3, false, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}, []float64{}, []float64{}},
				{3, false, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}, []float64{}, []float64{}},
				{3, false, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}, []float64{}, []float64{}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dataSet := []model.ElectricalConnectionPermittedValueSetDataType{}
			permittedData := []model.ScaledNumberSetType{}
			for _, data := range tc.data {
				for phase := 0; phase < data.phases; phase++ {
					item := model.ScaledNumberSetType{
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(data.permittedMinValue),
								Max: model.NewScaledNumberType(data.permittedMaxValue),
							},
						},
					}
					if data.permittedDefaultExists {
						item.Value = []model.ScaledNumberType{*model.NewScaledNumberType(data.permittedDefaultValue)}
					}
					permittedData = append(permittedData, item)

					permittedItem := model.ElectricalConnectionPermittedValueSetDataType{
						ElectricalConnectionId: eebusUtil.Ptr(model.ElectricalConnectionIdType(0)),
						ParameterId:            eebusUtil.Ptr(model.ElectricalConnectionParameterIdType(phase)),
						PermittedValueSet:      permittedData,
					}
					dataSet = append(dataSet, permittedItem)
				}

				datagram = datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer, model.RoleTypeClient)

				cmd = []model.CmdType{{
					ElectricalConnectionPermittedValueSetListData: &model.ElectricalConnectionPermittedValueSetListDataType{
						ElectricalConnectionPermittedValueSetData: dataSet,
					}}}
				datagram.Payload.Cmd = cmd

				err = localDevice.ProcessCmd(datagram, remoteDevice)
				assert.Nil(t, err)

				err = emobilty.EVWriteLoadControlLimits(loadLimits)
				assert.NotNil(t, err)

				datagram = datagramForEntityAndFeatures(false, localDevice, localEntity, emobilty.evEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer, model.RoleTypeClient)

				limitDesc := []model.LoadControlLimitDescriptionDataType{}
				var limitIdsObligation, limitIdsRecommendation []model.LoadControlLimitIdType
				for index := range data.obligations {
					id := model.LoadControlLimitIdType(index)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       eebusUtil.Ptr(id),
						LimitCategory: eebusUtil.Ptr(model.LoadControlCategoryTypeObligation),
						MeasurementId: eebusUtil.Ptr(model.MeasurementIdType(index)),
					}
					limitDesc = append(limitDesc, limitItem)
					limitIdsObligation = append(limitIdsObligation, id)
				}
				add := len(limitDesc)
				for index := range data.recommendations {
					id := model.LoadControlLimitIdType(index + add)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       eebusUtil.Ptr(id),
						LimitCategory: eebusUtil.Ptr(model.LoadControlCategoryTypeRecommendation),
						MeasurementId: eebusUtil.Ptr(model.MeasurementIdType(index + add)),
					}
					limitDesc = append(limitDesc, limitItem)
					limitIdsRecommendation = append(limitIdsRecommendation, id)
				}

				cmd = []model.CmdType{{
					LoadControlLimitDescriptionListData: &model.LoadControlLimitDescriptionListDataType{
						LoadControlLimitDescriptionData: limitDesc,
					}}}
				datagram.Payload.Cmd = cmd

				err = localDevice.ProcessCmd(datagram, remoteDevice)
				assert.Nil(t, err)

				err = emobilty.EVWriteLoadControlLimits(loadLimits)
				assert.NotNil(t, err)

				limitData := []model.LoadControlLimitDataType{}
				for index := range limitDesc {
					limitItem := model.LoadControlLimitDataType{
						LimitId:           eebusUtil.Ptr(model.LoadControlLimitIdType(index)),
						IsLimitChangeable: eebusUtil.Ptr(true),
					}
					limitData = append(limitData, limitItem)
				}
				sentLimits := len(limitData)

				cmd = []model.CmdType{{
					LoadControlLimitListData: &model.LoadControlLimitListDataType{
						LoadControlLimitData: limitData,
					}}}
				datagram.Payload.Cmd = cmd

				err = localDevice.ProcessCmd(datagram, remoteDevice)
				assert.Nil(t, err)

				err = emobilty.EVWriteLoadControlLimits(loadLimits)
				assert.NotNil(t, err)

				obligations := []EVLoadLimitsPhase{}
				for index, obligation := range data.obligations {
					phase := util.PhaseNameMapping[index]
					obligations = append(obligations, EVLoadLimitsPhase{
						Phase:    phase,
						IsActive: true,
						Value:    obligation,
					})
				}
				recommendations := []EVLoadLimitsPhase{}
				for index, obligation := range data.recommendations {
					phase := util.PhaseNameMapping[index]
					obligations = append(obligations, EVLoadLimitsPhase{
						Phase:    phase,
						IsActive: true,
						Value:    obligation,
					})
				}

				err = emobilty.EVWriteLoadControlLimits([]EVLoadLimits{
					{Category: model.LoadControlCategoryTypeObligation, PhaseData: obligations},
					{Category: model.LoadControlCategoryTypeRecommendation, PhaseData: recommendations},
				})
				assert.Nil(t, err)

				sentDatagram := model.Datagram{}
				sentBytes := writeHandler.LastMessage()
				err := json.Unmarshal(sentBytes, &sentDatagram)
				assert.Nil(t, err)

				sentCmd := sentDatagram.Datagram.Payload.Cmd
				assert.Equal(t, 1, len(sentCmd))

				sentLimitData := sentCmd[0].LoadControlLimitListData.LoadControlLimitData
				assert.Equal(t, sentLimits, len(sentLimitData))

				for _, item := range sentLimitData {
					if index := slices.Index(limitIdsObligation, *item.LimitId); index >= 0 {
						assert.Equal(t, data.obligationsExpected[index], item.Value.GetValue())
					}
					if index := slices.Index(limitIdsRecommendation, *item.LimitId); index >= 0 {
						assert.Equal(t, data.recommendationsExpected[index], item.Value.GetValue())
					}
				}
			}
		})
	}
}
