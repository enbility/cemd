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

	paramRef, err := electricalParamDescriptionListData(featureRemote)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	data := featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	resultSet := make(map[string]float64)
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil {
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

		if desc.ScopeType == nil || param.AcMeasuredPhases == nil || item.Value == nil {
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
		if item.MeasurementId == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil || item.Value == nil {
			continue
		}

		if *desc.ScopeType == model.ScopeTypeTypeStateOfCharge {
			return item.Value.GetValue(), nil
		}
	}

	return 0.0, ErrDataNotAvailable
}

type electricatlParamDescriptionMap map[model.MeasurementIdType]model.ElectricalConnectionParameterDescriptionDataType

// return a map of ElectricalConnectionParameterDescriptionListDataType with measurementId as key
func electricalParamDescriptionListData(featureRemote *spine.FeatureRemoteImpl) (electricatlParamDescriptionMap, error) {
	data := featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData).(*model.ElectricalConnectionParameterDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}
	ref := make(electricatlParamDescriptionMap)
	for _, item := range data.ElectricalConnectionParameterDescriptionData {
		if item.MeasurementId == nil {
			continue
		}
		ref[*item.MeasurementId] = item
	}

	return ref, nil
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
	data := featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData).(*model.MeasurementConstraintsListDataType)
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

	constraintsRef, err := measurementConstraintsListData(featureRemote)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

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
