package features

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
func RequestDeviceConfigurationKeyValueDescriptionList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
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
	/*
	   // subscribe to device configuration state updates
	   fErr = featureLocal.SubscribeAndWait(featureRemote.Device(), featureRemote.Address())

	   	if fErr != nil {
	   		fmt.Println(fErr.String())
	   	}
	*/
}

/*
// set the new device configuration data
func updateDeviceConfigurationData(entity *spine.EntityRemoteImpl) {
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


// request IdentificationListDataType from a remote entity
func (e *EV) requestIdentitificationlistData(entity *spine.EntityRemoteImpl) {
	knownEVData, ok := e.data[entity.Device().Ski()]
	if !ok || knownEVData.CommunicationStandard == EVCommunicationStandardTypeUnknown || knownEVData.CommunicationStandard == EVCommunicationStandardTypeIEC61851 {
		// identification requests only work with ISO connections to the EV
		return
	}

	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIdentification, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeIdentificationListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
	}

	e.updateIdentificationData(entity)
}

// set the new identification data
func (e *EV) updateIdentificationData(entity *spine.EntityRemoteImpl) {
	_, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIdentification, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := featureRemote.Data(model.FunctionTypeIdentificationListData).(*model.IdentificationListDataType)
	if data == nil {
		return
	}

	evData := e.dataForRemoteDevice(entity.Device())

	for _, dataItem := range data.IdentificationData {
		if dataItem.IdentificationType == nil {
			continue
		}

		evData.IdentificationType = EVIdentificationType(*dataItem.IdentificationType)
		evData.Identification = string(*dataItem.IdentificationValue)
	}

	fmt.Printf("EV Identification Type: %s\n", evData.IdentificationType)
	fmt.Printf("EV Identification: %s\n", evData.Identification)
}

*/
