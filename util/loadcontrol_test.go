package util

import (
	"testing"

	"github.com/enbility/cemd/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScopeLocal() {
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation
	direction := model.EnergyDirectionType("")

	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}

	exists := LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	payload.Entity = s.monitoredEntity

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
		},
	}

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	feature.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, descData)

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

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

	elFeature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	elFeature.SetData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData)

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	payload.Data = limitData
	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	limitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	payload.Data = limitData
	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(true, s.service, payload, limitType, category, direction, scope)
	assert.True(s.T(), exists)
}

func (s *UtilSuite) Test_LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope() {
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation
	direction := model.EnergyDirectionType("")

	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}

	exists := LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	payload.Entity = s.monitoredEntity

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	payload.Data = limitData
	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.False(s.T(), exists)

	limitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	payload.Data = limitData
	exists = LoadControlLimitsCheckPayloadDataForTypeCategoryDirectionScope(false, s.service, payload, limitType, category, direction, scope)
	assert.True(s.T(), exists)
}

func (s *UtilSuite) Test_LoadControlLimits() {
	var data []api.LoadLimitsPhase
	var err error
	limitType := model.LoadControlLimitTypeTypeMaxValueLimit
	scope := model.ScopeTypeTypeSelfConsumption
	category := model.LoadControlCategoryTypeObligation
	entityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}

	data, err = LoadControlLimits(s.service, s.mockRemoteEntity, entityTypes, limitType, category, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	data, err = LoadControlLimits(s.service, s.monitoredEntity, entityTypes, limitType, category, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: eebusutil.Ptr(category),
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				LimitType:     eebusutil.Ptr(limitType),
				ScopeType:     eebusutil.Ptr(scope),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.service, s.monitoredEntity, entityTypes, limitType, category, scope)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), 0.0, data[0].Value)

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

	rElFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.service, s.monitoredEntity, entityTypes, limitType, category, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitListData, limitData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.service, s.monitoredEntity, entityTypes, limitType, category, scope)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(2)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(6),
								Max: model.NewScaledNumberType(16),
							},
						},
					},
				},
			},
		},
	}

	fErr = rElFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = LoadControlLimits(s.service, s.monitoredEntity, entityTypes, limitType, category, scope)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 3, len(data))
	assert.Equal(s.T(), 16.0, data[0].Value)
}

func (s *UtilSuite) Test_WriteLoadControlLimits() {
	loadLimits := []api.LoadLimitsPhase{}

	category := model.LoadControlCategoryTypeObligation
	entityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}

	msgCounter, err := WriteLoadControlLimits(s.service, s.mockRemoteEntity, entityTypes, category, loadLimits)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), msgCounter)

	msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, loadLimits)
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, loadLimits)
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
				// clean up data
				remoteLoadControlF := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				assert.NotNil(s.T(), remoteLoadControlF)

				emptyLimits := model.LoadControlLimitListDataType{}
				errT := remoteLoadControlF.UpdateData(model.FunctionTypeLoadControlLimitListData, &emptyLimits, nil, nil)
				assert.Nil(s.T(), errT)

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

				msgCounter, err := WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, loadLimits)
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

				rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
				fErr = rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
				assert.Nil(s.T(), fErr)

				msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, loadLimits)
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

				msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, loadLimits)
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

				msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)

				msgCounter, err = WriteLoadControlLimits(s.service, s.monitoredEntity, entityTypes, category, phaseLimitValues)
				assert.Nil(t, err)
				assert.NotNil(t, msgCounter)
			}
		})
	}
}

func (s *UtilSuite) Test_GetLocalLimitDescriptionsForTypeCategoryDirectionScope() {
	limitType := model.LoadControlLimitTypeTypeSignDependentAbsValueLimit
	limitCategory := model.LoadControlCategoryTypeObligation
	limitDirection := model.EnergyDirectionTypeConsume
	limitScopeType := model.ScopeTypeTypeActivePowerLimit

	data := GetLocalLimitDescriptionsForTypeCategoryDirectionScope(s.service, limitType, limitCategory, limitDirection, limitScopeType)
	assert.Equal(s.T(), 0, len(data))

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	desc := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:      eebusutil.Ptr(limitType),
				LimitCategory:  eebusutil.Ptr(limitCategory),
				LimitDirection: eebusutil.Ptr(limitDirection),
				ScopeType:      eebusutil.Ptr(limitScopeType),
			},
		},
	}
	feature.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, desc)

	data = GetLocalLimitDescriptionsForTypeCategoryDirectionScope(s.service, limitType, limitCategory, limitDirection, limitScopeType)
	assert.Equal(s.T(), 1, len(data))
	assert.NotNil(s.T(), data[0].LimitId)
}

func (s *UtilSuite) Test_GetLocalLimitValueForLimitId() {
	limitId := model.LoadControlLimitIdType(0)

	data := GetLocalLimitValueForLimitId(s.service, limitId)
	assert.Nil(s.T(), data.LimitId)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	desc := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
			},
		},
	}
	feature.SetData(model.FunctionTypeLoadControlLimitListData, desc)

	data = GetLocalLimitValueForLimitId(s.service, limitId)
	assert.NotNil(s.T(), data.LimitId)
}
