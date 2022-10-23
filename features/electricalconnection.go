package features

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type ElectricalLimit struct {
	Min     float64
	Max     float64
	Default float64
	Phase   model.ElectricalConnectionPhaseNameType
	Type    model.ScopeTypeType
}

// request electrical connection data to properly interpret the corresponding data messages
func RequestElectricalConnection(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request ElectricalConnectionParameterDescriptionListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	// request ElectricalConnectionDescriptionListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	// request ElectricalConnectionPermittedValueSetListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	return nil
}

// return current values for Electrical Limits
func GetElectricalLimitValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]ElectricalLimit, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rData := featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}
	paramDescriptionData := rData.(*model.ElectricalConnectionParameterDescriptionListDataType)
	paramRef := make(map[*model.ElectricalConnectionParameterIdType]model.ElectricalConnectionParameterDescriptionDataType)
	for _, item := range paramDescriptionData.ElectricalConnectionParameterDescriptionData {
		if item.ParameterId == nil {
			continue
		}
		paramRef[item.ParameterId] = item
	}

	rData = featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}
	descriptionData := rData.(*model.ElectricalConnectionDescriptionListDataType)
	descRef := make(map[*model.ElectricalConnectionIdType]model.ElectricalConnectionDescriptionDataType)
	for _, item := range descriptionData.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}
		descRef[item.ElectricalConnectionId] = item
	}

	data := featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []ElectricalLimit
	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ParameterId == nil || item.ElectricalConnectionId == nil {
			continue
		}
		param, exists := paramRef[item.ParameterId]
		if !exists {
			continue
		}
		// desc, exists := descRef[item.ElectricalConnectionId]
		// if !exists {
		// 	continue
		// }

		if len(item.PermittedValueSet) == 0 {
			continue
		}

		var value, minValue, maxValue float64
		hasValue := false
		hasRange := false

		for _, element := range item.PermittedValueSet {
			// is a value set
			if element.Value != nil && len(element.Value) > 0 {
				value = element.Value[0].GetValue()
				hasValue = true
			}
			// is a range set
			if element.Range != nil && len(element.Range) > 0 {
				minValue = element.Range[0].Min.GetValue()
				maxValue = element.Range[0].Max.GetValue()
				hasRange = true
			}
		}

		switch {
		// AC Total Power Limits
		case param.ScopeType != nil && *param.ScopeType == model.ScopeTypeTypeACPowerTotal && hasRange:
			result := ElectricalLimit{
				Min:  minValue,
				Max:  maxValue,
				Type: model.ScopeTypeTypeACPowerTotal,
			}
			resultSet = append(resultSet, result)
		case param.AcMeasuredPhases != nil && hasRange && hasValue:
			// AC Phase Current Limits
			result := ElectricalLimit{
				Min:     minValue,
				Max:     maxValue,
				Default: value,
				Phase:   *param.AcMeasuredPhases,
				Type:    model.ScopeTypeTypeACCurrent,
			}
			resultSet = append(resultSet, result)
		}
	}

	return resultSet, nil
}
