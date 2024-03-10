package util

import (
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

func GetLocalElectricalConnectionCharacteristicForContextType(
	service eebusapi.ServiceInterface,
	context model.ElectricalConnectionCharacteristicContextType,
	charType model.ElectricalConnectionCharacteristicTypeType,
) (charData model.ElectricalConnectionCharacteristicDataType) {
	charData = model.ElectricalConnectionCharacteristicDataType{}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	electricalConnection := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	if electricalConnection == nil {
		return
	}

	function := model.FunctionTypeElectricalConnectionCharacteristicListData
	data, err := spine.LocalFeatureDataCopyOfType[*model.ElectricalConnectionCharacteristicListDataType](
		electricalConnection, function)
	if err != nil || data == nil || data.ElectricalConnectionCharacteristicListData == nil {
		return
	}

	for _, item := range data.ElectricalConnectionCharacteristicListData {
		if item.CharacteristicContext != nil && *item.CharacteristicContext == context &&
			item.CharacteristicType != nil && *item.CharacteristicType == charType {
			charData = item
			break
		}
	}

	return
}

func SetLocalElectricalConnectionCharacteristicForContextType(
	service eebusapi.ServiceInterface,
	context model.ElectricalConnectionCharacteristicContextType,
	charType model.ElectricalConnectionCharacteristicTypeType,
	value float64,
) (resultErr error) {
	resultErr = eebusapi.ErrDataNotAvailable

	charData := GetLocalElectricalConnectionCharacteristicForContextType(service, context, charType)
	if charData.CharacteristicId == nil {
		return
	}
	charData.Value = model.NewScaledNumberType(value)

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	electricalConnection := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	if electricalConnection == nil {
		return
	}
	function := model.FunctionTypeElectricalConnectionCharacteristicListData

	listData := &model.ElectricalConnectionCharacteristicListDataType{
		ElectricalConnectionCharacteristicListData: []model.ElectricalConnectionCharacteristicDataType{charData},
	}
	electricalConnection.SetData(function, listData)

	return nil
}
