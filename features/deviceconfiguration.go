package features

import (
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type DeviceConfigurationType struct {
	Key           string
	ValueBool     bool
	ValueDate     time.Time
	ValueDatetime time.Time
	ValueDuration time.Duration
	ValueString   string
	ValueTime     time.Time
	ValueFloat    float64
	Type          model.DeviceConfigurationKeyValueTypeType
	Unit          string
}

// request DeviceConfiguration data from a remote entity
func RequestDeviceConfiguration(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
	if _, err := requestData(featureLocal, featureRemote, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData); err != nil {
		fmt.Println(err)
		return err
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

	// request FunctionTypeDeviceConfigurationKeyValueListData from a remote entity
	msgCounter, err := requestData(featureLocal, featureRemote, model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil {
		fmt.Println(err)
		return nil, err
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
			Key: *desc.KeyName,
		}
		if desc.ValueType == nil {
			continue
		}
		result.Type = *desc.ValueType
		switch *desc.ValueType {
		case model.DeviceConfigurationKeyValueTypeTypeBoolean:
			if item.Value.Boolean != nil {
				result.ValueBool = bool(*item.Value.Boolean)
			}
		case model.DeviceConfigurationKeyValueTypeTypeDate:
			if item.Value.Date != nil {
				if value, err := model.GetDateFromString(*item.Value.Date); err == nil {
					result.ValueDate = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeDateTime:
			if item.Value.DateTime != nil {
				if value, err := model.GetDateTimeFromString(*item.Value.DateTime); err == nil {
					result.ValueDatetime = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeDuration:
			if item.Value.Duration != nil {
				if value, err := model.GetDurationFromString(*item.Value.Duration); err == nil {
					result.ValueDuration = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeString:
			if item.Value.String != nil {
				result.ValueString = string(*item.Value.String)
			}
		case model.DeviceConfigurationKeyValueTypeTypeTime:
			if item.Value.Time != nil {
				if value, err := model.GetTimeFromString(*item.Value.Time); err != nil {
					result.ValueTime = value
				}
			}
		case model.DeviceConfigurationKeyValueTypeTypeScalednumber:
			if item.Value.ScaledNumber != nil {
				result.ValueFloat = item.Value.ScaledNumber.GetValue()
			}
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
