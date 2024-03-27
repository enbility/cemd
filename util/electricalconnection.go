package util

import (
	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

func GetPhaseCurrentLimits(
	service eebusapi.ServiceInterface,
	entity spineapi.EntityRemoteInterface,
	entityTypes []model.EntityTypeType) (
	resultMin []float64, resultMax []float64, resultDefault []float64, resultErr error) {
	if !IsCompatibleEntity(entity, entityTypes) {
		return nil, nil, nil, api.ErrNoCompatibleEntity
	}

	evElectricalConnection, err := ElectricalConnection(service, entity)
	if err != nil {
		return nil, nil, nil, eebusapi.ErrDataNotAvailable
	}

	for _, phaseName := range PhaseNameMapping {
		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.ParameterId == nil {
			continue
		}

		dataMin, dataMax, dataDefault, err := evElectricalConnection.GetLimitsForParameterId(*elParamDesc.ParameterId)
		if err != nil {
			continue
		}

		// Min current data should be derived from min power data
		// but as this value is only properly provided via VAS the
		// currrent min values can not be trusted.

		resultMin = append(resultMin, dataMin)
		resultMax = append(resultMax, dataMax)
		resultDefault = append(resultDefault, dataDefault)
	}

	if len(resultMin) == 0 {
		return nil, nil, nil, eebusapi.ErrDataNotAvailable
	}

	return resultMin, resultMax, resultDefault, nil
}

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
