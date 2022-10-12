package cem

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

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
