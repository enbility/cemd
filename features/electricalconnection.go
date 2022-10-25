package features

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Details about the electrical connection
type ElectricalDescriptionType struct {
	ConnectionID            uint
	PowerSupplyType         model.ElectricalConnectionVoltageTypeType
	AcConnectedPhases       uint
	PositiveEnergyDirection model.EnergyDirectionType
}

// Details about the limits of an electrical connection
type ElectricalLimitType struct {
	ConnectionID uint
	Min          float64
	Max          float64
	Default      float64
	Phase        model.ElectricalConnectionPhaseNameType
	Scope        model.ScopeTypeType
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

// return current values for Electrical Description
func GetElectricalDescription(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]ElectricalDescriptionType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData).(*model.ElectricalConnectionDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}

	var resultSet []ElectricalDescriptionType

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}

		result := ElectricalDescriptionType{}

		if item.PowerSupplyType != nil {
			result.PowerSupplyType = *item.PowerSupplyType
		}
		if item.AcConnectedPhases != nil {
			result.AcConnectedPhases = *item.AcConnectedPhases
		}
		if item.PositiveEnergyDirection != nil {
			result.PositiveEnergyDirection = *item.PositiveEnergyDirection
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}

// return current values for Electrical Limits
//
// EV only: Min power data is only provided via IEC61851 or using VAS in ISO15118-2.
func GetElectricalLimitValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]ElectricalLimitType, error) {
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
	paramRef := make(map[model.ElectricalConnectionParameterIdType]model.ElectricalConnectionParameterDescriptionDataType)
	for _, item := range paramDescriptionData.ElectricalConnectionParameterDescriptionData {
		if item.ParameterId == nil {
			continue
		}
		paramRef[*item.ParameterId] = item
	}

	data := featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []ElectricalLimitType

	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ParameterId == nil || item.ElectricalConnectionId == nil {
			continue
		}
		param, exists := paramRef[*item.ParameterId]
		if !exists {
			continue
		}

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
			result := ElectricalLimitType{
				ConnectionID: uint(*item.ElectricalConnectionId),
				Min:          minValue,
				Max:          maxValue,
				Scope:        model.ScopeTypeTypeACPowerTotal,
			}
			resultSet = append(resultSet, result)

		case param.AcMeasuredPhases != nil && hasRange && hasValue:
			// AC Phase Current Limits
			result := ElectricalLimitType{
				ConnectionID: uint(*item.ElectricalConnectionId),
				Min:          minValue,
				Max:          maxValue,
				Default:      value,
				Phase:        *param.AcMeasuredPhases,
				Scope:        model.ScopeTypeTypeACCurrent,
			}
			resultSet = append(resultSet, result)
		}
	}

	return resultSet, nil
}
