package util

import (
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type MeasurementType struct {
	MeasurementId uint
	Value         float64
	ValueMin      float64
	ValueMax      float64
	ValueStep     float64
	Unit          model.UnitOfMeasurementType
	Scope         model.ScopeTypeType
	Timestamp     time.Time
}

// subscribe to measurements
func SubscribeMeasurementsForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	return subscribeToFeatureForEntity(service, model.FeatureTypeTypeMeasurement, entity)
}

// request FunctionTypeMeasurementDescriptionListData from a remote device
func RequestMeasurementDescription(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := requestData(featureLocal, featureRemote, model.FunctionTypeMeasurementDescriptionListData); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func RequestMeasurementConstraints(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := requestData(featureLocal, featureRemote, model.FunctionTypeMeasurementConstraintsListData); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeMeasurementListData from a remote entity
func RequestMeasurementList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// request FunctionTypeMeasurementListData from a remote entity
	msgCounter, err := requestData(featureLocal, featureRemote, model.FunctionTypeMeasurementListData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return msgCounter, nil
}

// return current current values
//
// returns a map with the phase ("a", "b", "c") as a key
func GetMeasurementCurrents(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (map[string]float64, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	descRef, err := measurementDescriptionListData(featureRemote)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	paramRef, _, err := electricalParamDescriptionListData(featureRemote)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	data := featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	resultSet := make(map[string]float64)
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		param, exists := paramRef[*item.MeasurementId]
		if !exists {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil || param.AcMeasuredPhases == nil {
			continue
		}

		if *desc.ScopeType == model.ScopeTypeTypeACCurrent {
			resultSet[string(*param.AcMeasuredPhases)] = item.Value.GetValue()
		}
	}
	if len(resultSet) == 0 {
		return nil, ErrDataNotAvailable
	}

	return resultSet, nil
}

// return number of phases the device is connected with
func GetElectricalConnectedPhases(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (uint, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	data := featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData).(*model.ElectricalConnectionDescriptionListDataType)
	if data == nil {
		return 0, ErrDataNotAvailable
	}

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}

		if item.AcConnectedPhases != nil {
			return *item.AcConnectedPhases, nil
		}
	}

	// default to 3 if the value is not available
	return 3, nil
}

// return current current limit values
//
// returns a map with the phase ("a", "b", "c") as a key for
// minimum, maximum, default/pause values
func GetElectricalCurrentsLimits(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (map[string]float64, map[string]float64, map[string]float64, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	_, paramRef, err := electricalParamDescriptionListData(featureRemote)
	if err != nil {
		return nil, nil, nil, ErrMetadataNotAvailable
	}

	data := featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if data == nil {
		return nil, nil, nil, ErrDataNotAvailable
	}

	resultSetMin := make(map[string]float64)
	resultSetMax := make(map[string]float64)
	resultSetDefault := make(map[string]float64)
	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ElectricalConnectionId == nil || item.PermittedValueSet == nil {
			continue
		}

		param, exists := paramRef[*item.ElectricalConnectionId]
		if !exists {
			continue
		}

		if param.AcMeasuredPhases == nil {
			continue
		}

		for _, set := range item.PermittedValueSet {
			if set.Value != nil && len(set.Value) > 0 {
				resultSetDefault[string(*param.AcMeasuredPhases)] = set.Value[0].GetValue()
			}
			if set.Range != nil {
				for _, rangeItem := range set.Range {
					if rangeItem.Min != nil {
						resultSetMin[string(*param.AcMeasuredPhases)] = rangeItem.Min.GetValue()
					}
					if rangeItem.Max != nil {
						resultSetMax[string(*param.AcMeasuredPhases)] = rangeItem.Max.GetValue()
					}
				}
			}
		}
	}

	if len(resultSetMin) == 0 && len(resultSetMax) == 0 && len(resultSetMax) == 0 {
		return nil, nil, nil, ErrDataNotAvailable
	}

	return resultSetMin, resultSetMax, resultSetDefault, nil
}

// returns if a provided scopetype in the measurement descriptions is available or not
// returns an error if no description data is available yet
func GetMeasurementDescriptionScopeSupport(scope model.ScopeTypeType, service *service.EEBUSService, entity *spine.EntityRemoteImpl) (bool, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	data := featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData).(*model.MeasurementDescriptionListDataType)
	if data == nil {
		return false, ErrDataNotAvailable
	}
	for _, item := range data.MeasurementDescriptionData {
		if item.MeasurementId == nil || item.ScopeType == nil {
			continue
		}
		if *item.ScopeType == scope {
			return true, nil
		}
	}

	return false, ErrDataNotAvailable
}

// return current SoC for measurements
func GetMeasurementSoC(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (float64, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}

	descRef, err := measurementDescriptionListData(featureRemote)
	if err != nil {
		return 0.0, ErrMetadataNotAvailable
	}

	data := featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return 0.0, ErrDataNotAvailable
	}

	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil {
			continue
		}

		if *desc.ScopeType == model.ScopeTypeTypeStateOfCharge {
			return item.Value.GetValue(), nil
		}
	}

	return 0.0, ErrDataNotAvailable
}

type electricatlParamDescriptionMapMeasurementId map[model.MeasurementIdType]model.ElectricalConnectionParameterDescriptionDataType
type electricatlParamDescriptionMapElectricalId map[model.ElectricalConnectionIdType]model.ElectricalConnectionParameterDescriptionDataType

// return a map of ElectricalConnectionParameterDescriptionListDataType with measurementId as key and
// ElectricalConnectionParameterDescriptionListDataType with electricalConnectionId as key
func electricalParamDescriptionListData(featureRemote *spine.FeatureRemoteImpl) (electricatlParamDescriptionMapMeasurementId, electricatlParamDescriptionMapElectricalId, error) {
	data := featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData).(*model.ElectricalConnectionParameterDescriptionListDataType)
	if data == nil {
		return nil, nil, ErrDataNotAvailable
	}
	refMeasurement := make(electricatlParamDescriptionMapMeasurementId)
	refElectrical := make(electricatlParamDescriptionMapElectricalId)
	for _, item := range data.ElectricalConnectionParameterDescriptionData {
		if item.MeasurementId == nil || item.ElectricalConnectionId == nil {
			continue
		}
		refMeasurement[*item.MeasurementId] = item
		refElectrical[*item.ElectricalConnectionId] = item
	}

	return refMeasurement, refElectrical, nil
}

type measurementDescriptionMap map[model.MeasurementIdType]model.MeasurementDescriptionDataType

// return a map of MeasurementDescriptionListDataType with measurementId as key
func measurementDescriptionListData(featureRemote *spine.FeatureRemoteImpl) (measurementDescriptionMap, error) {
	data := featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData).(*model.MeasurementDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}
	ref := make(measurementDescriptionMap)
	for _, item := range data.MeasurementDescriptionData {
		if item.MeasurementId == nil {
			continue
		}
		ref[*item.MeasurementId] = item
	}
	return ref, nil
}

type measurementConstraintMap map[model.MeasurementIdType]model.MeasurementConstraintsDataType

// return a map of MeasurementDescriptionListDataType with measurementId as key
func measurementConstraintsListData(featureRemote *spine.FeatureRemoteImpl) (measurementConstraintMap, error) {
	data := featureRemote.Data(model.FunctionTypeMeasurementConstraintsListData).(*model.MeasurementConstraintsListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}
	ref := make(measurementConstraintMap)
	for _, item := range data.MeasurementConstraintsData {
		if item.MeasurementId == nil {
			continue
		}
		ref[*item.MeasurementId] = item
	}
	return ref, nil
}

// return current values for measurements
func GetMeasurementValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]MeasurementType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// constraints can be optional
	constraintsRef, _ := measurementConstraintsListData(featureRemote)

	descRef, err := measurementDescriptionListData(featureRemote)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	data := featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []MeasurementType
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		result := MeasurementType{
			MeasurementId: uint(*item.MeasurementId),
		}

		if item.Value != nil {
			result.Value = item.Value.GetValue()
		}

		if item.Timestamp != nil {
			if value, err := time.Parse(time.RFC3339, *item.Timestamp); err == nil {
				result.Timestamp = value
			}
		}

		if desc.ScopeType != nil {
			result.Scope = *desc.ScopeType
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		constraint, exists := constraintsRef[*item.MeasurementId]
		if exists {
			if constraint.ValueRangeMin != nil {
				result.ValueMin = constraint.ValueRangeMin.GetValue()
			}
			if constraint.ValueRangeMax != nil {
				result.ValueMax = constraint.ValueRangeMax.GetValue()
			}
			if constraint.ValueStepSize != nil {
				result.ValueStep = constraint.ValueStepSize.GetValue()
			}
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
