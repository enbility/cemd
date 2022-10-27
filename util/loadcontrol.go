package util

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type LoadControlLimitType struct {
	LimitId       uint
	MeasurementId uint
	Category      model.LoadControlCategoryType
	Unit          model.UnitOfMeasurementType
	Scope         model.ScopeTypeType
	IsChangeable  bool
	IsActive      bool
	Value         float64
}

// request FunctionTypeLoadControlLimitDescriptionListData from a remote device
func RequestLoadControlLimitDescription(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := requestData(featureLocal, featureRemote, model.FunctionTypeLoadControlLimitDescriptionListData); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeLoadControlLimitConstraintsListData from a remote device
func RequestLoadControlLimitConstraints(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if _, err := requestData(featureLocal, featureRemote, model.FunctionTypeLoadControlLimitConstraintsListData); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeLoadControlLimitListData from a remote device
func RequestLoadControlLimitList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, err := requestData(featureLocal, featureRemote, model.FunctionTypeLoadControlLimitListData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return msgCounter, nil
}

// returns if a provided category in the load control descriptions is available or not
// returns an error if no description data is available yet
func GetLoadControlDescriptionCategorySupport(category model.LoadControlCategoryType, service *service.EEBUSService, entity *spine.EntityRemoteImpl) (bool, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	data := featureRemote.Data(model.FunctionTypeLoadControlLimitDescriptionListData).(*model.LoadControlLimitDescriptionListDataType)
	if data == nil {
		return false, ErrDataNotAvailable
	}
	for _, item := range data.LoadControlLimitDescriptionData {
		if item.LimitId == nil || item.LimitCategory == nil {
			continue
		}
		if *item.LimitCategory == category {
			return true, nil
		}
	}

	return false, ErrDataNotAvailable
}

func GetLoadControlLimitValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]LoadControlLimitType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	descriptionData := featureRemote.Data(model.FunctionTypeLoadControlLimitDescriptionListData).(*model.LoadControlLimitDescriptionListDataType)
	if descriptionData == nil {
		return nil, ErrMetadataNotAvailable
	}
	descRef := make(map[model.LoadControlLimitIdType]model.LoadControlLimitDescriptionDataType)
	for _, item := range descriptionData.LoadControlLimitDescriptionData {
		if item.MeasurementId == nil {
			continue
		}
		descRef[*item.LimitId] = item
	}

	data := featureRemote.Data(model.FunctionTypeLoadControlLimitListData).(*model.LoadControlLimitListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []LoadControlLimitType
	for _, item := range data.LoadControlLimitData {
		if item.LimitId == nil {
			continue
		}

		desc, exists := descRef[*item.LimitId]
		if !exists {
			continue
		}

		result := LoadControlLimitType{
			LimitId: uint(*item.LimitId),
		}

		if desc.MeasurementId != nil {
			result.MeasurementId = uint(*desc.MeasurementId)
		}
		if desc.LimitCategory != nil {
			result.Category = *desc.LimitCategory
		}
		if desc.ScopeType != nil {
			result.Scope = *desc.ScopeType
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		if item.IsLimitActive != nil {
			result.IsActive = *item.IsLimitActive
		}
		if item.IsLimitChangeable != nil {
			result.IsChangeable = *item.IsLimitChangeable
		}
		if item.Value != nil {
			result.Value = item.Value.GetValue()
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
