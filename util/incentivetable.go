package util

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// subscribe to time series
func SubscribeIncentiveTableForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	return subscribeToFeatureForEntity(service, model.FeatureTypeTypeIncentiveTable, entity)
}

// request FunctionTypeIncentiveTableDescriptionData from a remote entity
func RequestIncentiveTableDescription(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIncentiveTable, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = requestData(featureLocal, featureRemote, model.FunctionTypeIncentiveTableDescriptionData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func RequestIncentiveTableConstraints(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIncentiveTable, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = requestData(featureLocal, featureRemote, model.FunctionTypeIncentiveTableConstraintsData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableData from a remote entity
func RequestIncentiveTableValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIncentiveTable, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = requestData(featureLocal, featureRemote, model.FunctionTypeIncentiveTableData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
