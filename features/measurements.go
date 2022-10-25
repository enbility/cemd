package features

import (
	"errors"
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

// request measurement data to properly interpret the corresponding data messages
func RequestMeasurement(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request MeasurementDescriptionListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	// request FunctionTypeMeasurementConstraintsListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementConstraintsListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	// request FunctionTypeMeasurementListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	return nil
}

// return current values for measurements
func GetMeasurementValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]MeasurementType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rData := featureRemote.Data(model.FunctionTypeMeasurementConstraintsListData)
	// Constraints are optional, data may be empty
	constraintsRef := make(map[model.MeasurementIdType]model.MeasurementConstraintsDataType)
	switch constraintsData := rData.(type) {
	case *model.MeasurementConstraintsListDataType:
		if constraintsData != nil {
			for _, item := range constraintsData.MeasurementConstraintsData {
				if item.MeasurementId == nil {
					continue
				}
				constraintsRef[*item.MeasurementId] = item
			}
		}
	}

	rData = featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}
	descriptionData := rData.(*model.MeasurementDescriptionListDataType)
	descRef := make(map[model.MeasurementIdType]model.MeasurementDescriptionDataType)
	for _, item := range descriptionData.MeasurementDescriptionData {
		if item.MeasurementId == nil {
			continue
		}
		descRef[*item.MeasurementId] = item
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
