package util

import (
	eebusapi "github.com/enbility/eebus-go/api"
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

func DeviceConfigurationCheckDataPayloadForKeyName(localServer bool, service eebusapi.ServiceInterface,
	payload spineapi.EventPayload, keyName model.DeviceConfigurationKeyNameType) bool {
	var desc *model.DeviceConfigurationKeyValueDescriptionDataType
	var data *model.DeviceConfigurationKeyValueListDataType

	if payload.Data == nil {
		return false
	}
	data = payload.Data.(*model.DeviceConfigurationKeyValueListDataType)

	if localServer {
		desc = GetLocalDeviceConfigurationDescriptionForKeyName(service, keyName)
	} else {
		deviceconfigF, err := DeviceConfiguration(service, payload.Entity)
		if err != nil || payload.Data == nil {
			return false
		}

		desc, err = deviceconfigF.GetDescriptionForKeyName(keyName)
		if err != nil {
			return false
		}
	}

	for _, item := range data.DeviceConfigurationKeyValueData {
		if item.KeyId == nil || *item.KeyId != *desc.KeyId ||
			item.Value == nil {
			continue
		}

		return true
	}

	return false
}

func GetLocalDeviceConfigurationDescriptionForKeyName(
	service eebusapi.ServiceInterface,
	keyName model.DeviceConfigurationKeyNameType,
) (description *model.DeviceConfigurationKeyValueDescriptionDataType) {
	description = &model.DeviceConfigurationKeyValueDescriptionDataType{}

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
			return &desc
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
