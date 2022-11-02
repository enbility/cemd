package util

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type IdentificationType struct {
	Identifier string
	Type       model.IdentificationTypeType
}

// subscribe to identification
func SubscribeIdentificationForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	return subscribeToFeatureForEntity(service, model.FeatureTypeTypeIdentification, entity)
}

// request FunctionTypeIdentificationListData from a remote entity
func RequestIdentification(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIdentification, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// request FunctionTypeIdentificationListDataType from a remote entity
	msgCounter, err := requestData(featureLocal, featureRemote, model.FunctionTypeIdentificationListData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return msgCounter, nil
}

// return current values for Identification
func GetIdentificationValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]IdentificationType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIdentification, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rData := featureRemote.Data(model.FunctionTypeIdentificationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IdentificationListDataType)
	var resultSet []IdentificationType

	for _, item := range data.IdentificationData {
		if item.IdentificationValue == nil {
			continue
		}

		result := IdentificationType{
			Identifier: string(*item.IdentificationValue),
		}
		if item.IdentificationType != nil {
			result.Type = *item.IdentificationType
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
