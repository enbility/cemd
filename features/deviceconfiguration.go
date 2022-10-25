package features

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceConfigurationType struct {
	Key   string
	Value any
	Type  model.DeviceConfigurationKeyValueTypeType
	Unit  string
}

// request DeviceConfiguration data from a remote entity
func RequestDeviceConfiguration(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	return nil
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func RequestDeviceConfigurationKeyValueList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeDeviceConfigurationKeyValueListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// return current values for Device Configuration
func GetDeviceConfigurationValues(service *service.EEBUSService, entity *spine.EntityRemoteImpl) ([]DeviceConfigurationType, error) {
	_, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rDescData := featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if rDescData == nil {
		return nil, ErrMetadataNotAvailable
	}
	descData := rDescData.(*model.DeviceConfigurationKeyValueDescriptionListDataType)

	ref := make(map[model.DeviceConfigurationKeyIdType]model.DeviceConfigurationKeyValueDescriptionDataType)
	for _, item := range descData.DeviceConfigurationKeyValueDescriptionData {
		if item.KeyName == nil || item.KeyId == nil {
			continue
		}
		ref[*item.KeyId] = item
	}

	rData := featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceConfigurationKeyValueListDataType)
	var resultSet []DeviceConfigurationType

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil {
			continue
		}
		desc, exists := ref[*item.KeyId]
		if !exists || desc.KeyName == nil {
			continue
		}

		result := DeviceConfigurationType{
			Key:   *desc.KeyName,
			Value: item.Value,
		}
		if desc.ValueType != nil {
			result.Type = *desc.ValueType
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
