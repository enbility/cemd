package util

import (
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

func GetLocalDeviceConfigurationDescriptionForKeyName(
	service eebusapi.ServiceInterface,
	keyName model.DeviceConfigurationKeyNameType,
) (description model.DeviceConfigurationKeyValueDescriptionDataType) {
	description = model.DeviceConfigurationKeyValueDescriptionDataType{}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	deviceConfiguration := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	if deviceConfiguration == nil {
		return
	}

	data, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueDescriptionListDataType](
		deviceConfiguration, model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData)
	if err != nil || data == nil || data.DeviceConfigurationKeyValueDescriptionData == nil {
		return
	}

	for _, desc := range data.DeviceConfigurationKeyValueDescriptionData {
		if desc.KeyName != nil && *desc.KeyName == keyName {
			description = desc
			break
		}
	}

	return
}

func GetLocalDeviceConfigurationKeyValueForKeyName(
	service eebusapi.ServiceInterface,
	keyName model.DeviceConfigurationKeyNameType,
) (keyData model.DeviceConfigurationKeyValueDataType) {
	keyData = model.DeviceConfigurationKeyValueDataType{}

	description := GetLocalDeviceConfigurationDescriptionForKeyName(service, keyName)
	if description.KeyId == nil {
		return
	}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	deviceConfiguration := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	if deviceConfiguration == nil {
		return
	}

	data, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueListDataType](
		deviceConfiguration, model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil || data == nil || data.DeviceConfigurationKeyValueData == nil {
		return
	}

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId != nil && *item.KeyId == *description.KeyId {
			keyData = item
			break
		}
	}

	return
}

func SetLocalDeviceConfigurationKeyValue(
	service eebusapi.ServiceInterface,
	keyName model.DeviceConfigurationKeyNameType,
	changeable bool,
	value model.DeviceConfigurationKeyValueValueType,
) (resultErr error) {
	resultErr = eebusapi.ErrDataNotAvailable

	description := GetLocalDeviceConfigurationDescriptionForKeyName(service, keyName)
	if description.KeyId == nil {
		return
	}

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	deviceConfiguration := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	if deviceConfiguration == nil {
		return
	}

	data, err := spine.LocalFeatureDataCopyOfType[*model.DeviceConfigurationKeyValueListDataType](deviceConfiguration, model.FunctionTypeDeviceConfigurationKeyValueListData)
	if err != nil {
		data = &model.DeviceConfigurationKeyValueListDataType{}
	}

	found := false
	for index, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil || *item.KeyId != *description.KeyId {
			continue
		}

		item.IsValueChangeable = eebusutil.Ptr(changeable)
		item.Value = eebusutil.Ptr(value)

		data.DeviceConfigurationKeyValueData[index] = item
		found = true
	}

	if !found {
		item := model.DeviceConfigurationKeyValueDataType{
			KeyId:             eebusutil.Ptr(*description.KeyId),
			IsValueChangeable: eebusutil.Ptr(changeable),
			Value:             eebusutil.Ptr(value),
		}
		data.DeviceConfigurationKeyValueData = append(data.DeviceConfigurationKeyValueData, item)
	}

	deviceConfiguration.SetData(model.FunctionTypeDeviceConfigurationKeyValueListData, data)

	return nil
}
