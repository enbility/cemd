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

// bind to load control so we can write limits
func BindLoadControlLimit(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	return bindToFeatureForEntity(service, model.FeatureTypeTypeLoadControl, entity)
}

// subscribe to load control
func SubscribeLoadControlForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	return subscribeToFeatureForEntity(service, model.FeatureTypeTypeLoadControl, entity)
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

type loadControlLimitDescriptionMap map[model.LoadControlLimitIdType]model.LoadControlLimitDescriptionDataType

// returns the load control descriptions
// returns an error if no description data is available yet
func GetLoadControlLimitDescription(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (loadControlLimitDescriptionMap, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	data := featureRemote.Data(model.FunctionTypeLoadControlLimitDescriptionListData).(*model.LoadControlLimitDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}

	ref := make(loadControlLimitDescriptionMap)
	for _, item := range data.LoadControlLimitDescriptionData {
		if item.LimitId == nil {
			continue
		}
		ref[*item.LimitId] = item
	}

	return ref, nil
}

// returns if a provided category in the load control limit descriptions is available or not
// returns an error if no description data is available yet
func GetLoadControlLimitDescriptionCategorySupport(category model.LoadControlCategoryType, service *service.EEBUSService, entity *spine.EntityRemoteImpl) (bool, error) {
	data, err := GetLoadControlLimitDescription(service, entity)
	if err != nil {
		return false, err
	}

	for _, item := range data {
		if item.LimitId == nil || item.LimitCategory == nil {
			continue
		}
		if *item.LimitCategory == category {
			return true, nil
		}
	}

	return false, ErrDataNotAvailable
}

// write load control limits
// returns an error if this failed
func WriteLoadControlLimitValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl, data []model.LoadControlLimitDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, ErrMissingData
	}

	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeLoadControl, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cmd := []model.CmdType{{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: data,
		},
	}}

	return featureRemote.Sender().Write(featureLocal.Address(), featureRemote.Address(), cmd)
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
