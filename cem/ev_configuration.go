package cem

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
func (e *EV) requestConfigurationKeyValueDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	e.requestConfigurationKeyValueListData(entity)
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func (e *EV) requestConfigurationKeyValueListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeDeviceConfigurationKeyValueListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	e.updateDeviceConfigurationData(entity)

	// subscribe to device configuration state updates
	fErr = featureLocal.SubscribeAndWait(featureRemote.Device(), featureRemote.Address())
	if fErr != nil {
		fmt.Println(fErr.String())
	}
}

// set the new device configuration data
func (e *EV) updateDeviceConfigurationData(entity *spine.EntityRemoteImpl) {
	_, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	descriptionData := featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData).(*model.DeviceConfigurationKeyValueDescriptionListDataType)
	data := featureRemote.Data(model.FunctionTypeDeviceConfigurationKeyValueListData).(*model.DeviceConfigurationKeyValueListDataType)
	if descriptionData == nil || data == nil {
		return
	}

	evData := e.dataForRemoteDevice(entity.Device())

	for _, descriptionItem := range descriptionData.DeviceConfigurationKeyValueDescriptionData {
		for _, dataItem := range data.DeviceConfigurationKeyValueData {
			if *descriptionItem.KeyId != *dataItem.KeyId {
				continue
			}

			if descriptionItem.KeyName == nil {
				continue
			}

			switch *descriptionItem.KeyName {
			case string(model.DeviceConfigurationKeyNameTypeCommunicationsStandard):
				evData.CommunicationStandard = EVCommunicationStandardType(*dataItem.Value.String)
			case string(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported):
				evData.AsymmetricChargingSupported = (*dataItem.Value.Boolean)
			}
		}
	}

	fmt.Printf("EV Communication Standard: %s\n", evData.CommunicationStandard)
	fmt.Printf("EV Asymmetric Charging Supported: %t\n", evData.AsymmetricChargingSupported)

	// get ev identification data
	e.requestIdentitificationlistData(entity)
}
