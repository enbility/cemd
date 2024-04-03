package util

import (
	"testing"

	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_EVCurrentLimits() {
	entityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}

	minData, maxData, defaultData, err := GetPhaseCurrentLimits(s.service, s.mockRemoteEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = GetPhaseCurrentLimits(s.service, s.monitoredEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

	minData, maxData, defaultData, err = GetPhaseCurrentLimits(s.service, s.monitoredEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

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

	minData, maxData, defaultData, err = GetPhaseCurrentLimits(s.service, s.monitoredEntity, entityTypes)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), minData)
	assert.Nil(s.T(), maxData)
	assert.Nil(s.T(), defaultData)

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
				{true, 0.1, 0.1, 2, 2, 16, 16},
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
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
				{true, 0.1, 0.1, 2, 2, 16, 16},
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

	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
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
					ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
					ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(index)),
					PermittedValueSet:      permittedData,
				}
				dataSet = append(dataSet, permittedItem)
			}

			permData := &model.ElectricalConnectionPermittedValueSetListDataType{
				ElectricalConnectionPermittedValueSetData: dataSet,
			}

			fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
			assert.Nil(s.T(), fErr)

			minData, maxData, defaultData, err = GetPhaseCurrentLimits(s.service, s.monitoredEntity, entityTypes)
			assert.Nil(s.T(), err)

			assert.Nil(s.T(), err)
			assert.Equal(s.T(), len(tc.permitted), len(minData))
			assert.Equal(s.T(), len(tc.permitted), len(maxData))
			assert.Equal(s.T(), len(tc.permitted), len(defaultData))
			for index, item := range tc.permitted {
				assert.Equal(s.T(), item.expectedMinValue, minData[index])
				assert.Equal(s.T(), item.expectedMaxValue, maxData[index])
				assert.Equal(s.T(), item.expectedDefaultValue, defaultData[index])
			}
		})
	}
}

func (s *UtilSuite) Test_GetLocalElectricalConnectionCharacteristicForContextType() {
	context := model.ElectricalConnectionCharacteristicContextTypeEntity
	charType := model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax

	data := GetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType)
	assert.Nil(s.T(), data.CharacteristicId)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       eebusutil.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  eebusutil.Ptr(context),
				CharacteristicType:     eebusutil.Ptr(charType),
			},
		},
	}
	feature.SetData(model.FunctionTypeElectricalConnectionCharacteristicListData, charData)

	data = GetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType)
	assert.NotNil(s.T(), data.CharacteristicId)
}

func (s *UtilSuite) Test_SetLocalElectricalConnectionCharacteristicForContextType() {
	context := model.ElectricalConnectionCharacteristicContextTypeEntity
	charType := model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax
	value := 10.0

	err := SetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType, value)
	assert.NotNil(s.T(), err)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicData: []model.ElectricalConnectionCharacteristicDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				CharacteristicId:       eebusutil.Ptr(model.ElectricalConnectionCharacteristicIdType(0)),
				CharacteristicContext:  eebusutil.Ptr(context),
				CharacteristicType:     eebusutil.Ptr(charType),
			},
		},
	}
	feature.SetData(model.FunctionTypeElectricalConnectionCharacteristicListData, charData)

	err = SetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType, value)
	assert.Nil(s.T(), err)

	data := GetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType)
	assert.NotNil(s.T(), data.CharacteristicId)
	assert.Equal(s.T(), uint(0), uint(*data.CharacteristicId))
	assert.NotNil(s.T(), data.Value)
	assert.Equal(s.T(), 10.0, data.Value.GetValue())
}
