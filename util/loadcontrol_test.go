package util

import (
	"testing"

	"github.com/enbility/cemd/api"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *LoadControlSuite) Test_EVWriteLoadControlLimits() {
	loadLimits := []api.LoadLimitsPhase{}

	msgCounter, err := WriteLoadControlLimits(s.service, s.mockRemoteEntity, model.LoadControlCategoryTypeObligation, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	type dataStruct struct {
		phases                 int
		permittedDefaultExists bool
		permittedDefaultValue  float64
		permittedMinValue      float64
		permittedMaxValue      float64
		limits, limitsExpected []float64
	}

	tests := []struct {
		name string
		data []dataStruct
	}{
		{
			"1 Phase ISO15118",
			[]dataStruct{
				{1, true, 0.1, 2, 16, []float64{0}, []float64{0.1}},
				{1, true, 0.1, 2, 16, []float64{2.2}, []float64{2.2}},
				{1, true, 0.1, 2, 16, []float64{10}, []float64{10}},
				{1, true, 0.1, 2, 16, []float64{16}, []float64{16}},
			},
		},
		{
			"3 Phase ISO15118",
			[]dataStruct{
				{3, true, 0.1, 2, 16, []float64{0, 0, 0}, []float64{0.1, 0.1, 0.1}},
				{3, true, 0.1, 2, 16, []float64{2.2, 2.2, 2.2}, []float64{2.2, 2.2, 2.2}},
				{3, true, 0.1, 2, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, true, 0.1, 2, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
		{
			"1 Phase IEC61851",
			[]dataStruct{
				{1, true, 0, 6, 16, []float64{0}, []float64{0}},
				{1, true, 0, 6, 16, []float64{6}, []float64{6}},
				{1, true, 0, 6, 16, []float64{10}, []float64{10}},
				{1, true, 0, 6, 16, []float64{16}, []float64{16}},
			},
		},
		{
			"3 Phase IEC61851",
			[]dataStruct{
				{3, true, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}},
				{3, true, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}},
				{3, true, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, true, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
		{
			"3 Phase IEC61851 Elli",
			[]dataStruct{
				{3, false, 0, 6, 16, []float64{0, 0, 0}, []float64{0, 0, 0}},
				{3, false, 0, 6, 16, []float64{6, 6, 6}, []float64{6, 6, 6}},
				{3, false, 0, 6, 16, []float64{10, 10, 10}, []float64{10, 10, 10}},
				{3, false, 0, 6, 16, []float64{16, 16, 16}, []float64{16, 16, 16}},
			},
		},
	}

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
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
						ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
						ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(phase)),
						PermittedValueSet:      permittedData,
					}
					dataSet = append(dataSet, permittedItem)
				}

				permData := &model.ElectricalConnectionPermittedValueSetListDataType{
					ElectricalConnectionPermittedValueSetData: dataSet,
				}

				fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err := WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, loadLimits)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				limitDesc := []model.LoadControlLimitDescriptionDataType{}
				for index := range data.limits {
					id := model.LoadControlLimitIdType(index)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       eebusutil.Ptr(id),
						LimitCategory: eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
						MeasurementId: eebusutil.Ptr(model.MeasurementIdType(index)),
					}
					limitDesc = append(limitDesc, limitItem)
				}
				add := len(limitDesc)
				for index := range data.limits {
					id := model.LoadControlLimitIdType(index + add)
					limitItem := model.LoadControlLimitDescriptionDataType{
						LimitId:       eebusutil.Ptr(id),
						LimitCategory: eebusutil.Ptr(model.LoadControlCategoryTypeRecommendation),
						MeasurementId: eebusutil.Ptr(model.MeasurementIdType(index)),
					}
					limitDesc = append(limitDesc, limitItem)
				}

				descData := &model.LoadControlLimitDescriptionListDataType{
					LoadControlLimitDescriptionData: limitDesc,
				}

				rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, loadLimits)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				limitData := []model.LoadControlLimitDataType{}
				for index := range limitDesc {
					limitItem := model.LoadControlLimitDataType{
						LimitId:           eebusutil.Ptr(model.LoadControlLimitIdType(index)),
						IsLimitChangeable: eebusutil.Ptr(true),
						IsLimitActive:     eebusutil.Ptr(false),
						Value:             model.NewScaledNumberType(data.permittedMaxValue),
					}
					limitData = append(limitData, limitItem)
				}

				limitListData := &model.LoadControlLimitListDataType{
					LoadControlLimitData: limitData,
				}

				fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitListData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, loadLimits)
				assert.NotNil(t, err)
				assert.Nil(t, msgCounter)

				phaseLimitValues := []api.LoadLimitsPhase{}
				for index, limit := range data.limits {
					phase := PhaseNameMapping[index]
					phaseLimitValues = append(phaseLimitValues, api.LoadLimitsPhase{
						Phase:    phase,
						IsActive: true,
						Value:    limit,
					})

				}
				msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeObligation, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)

				msgCounter, err = WriteLoadControlLimits(s.service, s.evEntity, model.LoadControlCategoryTypeRecommendation, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)
			}
		})
	}
}
