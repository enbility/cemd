package util

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_GetLocalElectricalConnectionCharacteristicForContextType() {
	context := model.ElectricalConnectionCharacteristicContextTypeEntity
	charType := model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax

	data := GetLocalElectricalConnectionCharacteristicForContextType(s.service, context, charType)
	assert.Nil(s.T(), data.CharacteristicId)

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)

	charData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicListData: []model.ElectricalConnectionCharacteristicDataType{
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
		ElectricalConnectionCharacteristicListData: []model.ElectricalConnectionCharacteristicDataType{
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
