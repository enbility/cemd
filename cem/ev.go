package cem

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EVCommunicationStandardType string

const (
	EVCommunicationStandardTypeUnknown      EVCommunicationStandardType = "unknown"
	EVCommunicationStandardTypeISO151182ED1 EVCommunicationStandardType = "iso15118-2ed1"
	EVCommunicationStandardTypeISO151182ED2 EVCommunicationStandardType = "iso15118-2ed2"
	EVCommunicationStandardTypeIEC61851     EVCommunicationStandardType = "iec61851"
)

type EVIdentificationType string

const (
	EVIdentificationTypeEUI48 EVIdentificationType = "eui48" // eui48 MAC address
	EVIdentificationTypeEUI64 EVIdentificationType = "eui64" // eui64 MAC address
)

type EVData struct {
	CommunicationStandard       EVCommunicationStandardType
	AsymmetricChargingSupported bool
	IdentificationType          EVIdentificationType
	Identification              string
	ManufacturerDetails         ManufacturerDetails
}

// Interface for the evCC use case for CEM device
type EVDelegate interface {
	// handle device state updates from the remote EV entity
	HandleEVEntityState(ski string, failure bool)
}

// EV Commissioning and Configuration Use Case implementation
type EV struct {
	*spine.UseCaseImpl
	service *service.EEBUSService

	Delegate EVDelegate

	// map of device SKIs to EVData
	data map[string]*EVData
}

// Register the use case and features for handling EVs
// CEM will call this on startup
func AddEVSupport(service *service.EEBUSService) (*EV, error) {
	if service.ServiceDescription.DeviceType != model.DeviceTypeTypeEnergyManagementSystem {
		return nil, errors.New("device type not supported")
	}

	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &EV{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeEVCommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8}),
		service: service,
		data:    make(map[string]*EVData),
	}

	// subscribe to get incoming EV events
	spine.Events.Subscribe(useCase)

	// add the features
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
		entity.AddFeature(f)
	}
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
		entity.AddFeature(f)
	}
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)

		// Set the initial state
		state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
		deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
			OperatingState: &state,
		}
		f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

		entity.AddFeature(f)
	}
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
		entity.AddFeature(f)
	}
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
		entity.AddFeature(f)
	}

	return useCase, nil
}

// get the remote device specific data element
func (e *EV) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *EVData {
	if evdata, ok := e.data[remoteDevice.Ski()]; ok {
		return evdata
	}

	return &EVData{
		CommunicationStandard:       EVCommunicationStandardTypeIEC61851,
		AsymmetricChargingSupported: false,
	}
}

// Invoke to remove an EV entity
// Called when an EV was disconnected
func (e *EV) UnregisterEV() {
	// remove the entity
	e.service.RemoveEntity(e.Entity)
}

// Invoked when an EV entity was added or removed
func (e *EV) TriggerEntityUpdate() {

}

// Internal EventHandler Interface for the CEM
func (e *EV) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			// check if an EV is already connected
			remoteDevice := payload.Device
			if remoteDevice == nil {
				return
			}
			// Attention: We assume an EVSE only has 1 port!
			entity := remoteDevice.Entity([]model.AddressEntityType{1, 1})
			if !e.checkEntityBeingEV(entity) {
				return
			}
			e.evConnected(entity)
		}
	case spine.EventTypeEntityChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			// EV connected
			if !e.checkEntityBeingEV(payload.Entity) {
				return
			}
			e.evConnected(payload.Entity)
		case spine.ElementChangeRemove:
			// EV disconnected
			if !e.checkEntityBeingEV(payload.Entity) {
				return
			}
			fmt.Println("EV DISCONNECTED")
		}
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				if e.Delegate == nil {
					return
				}

				deviceDiagnosisStateData := payload.Data.(*model.DeviceDiagnosisStateDataType)
				failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				e.Delegate.HandleEVEntityState(payload.Ski, failure)
			}
		}
	}
}

// check if the provided entity is an EV
func (e *EV) checkEntityBeingEV(entity *spine.EntityRemoteImpl) bool {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false
	}
	return true
}

// an EV was connected, trigger required communication
func (e *EV) evConnected(entity *spine.EntityRemoteImpl) {
	fmt.Println("EV CONNECTED")

	// get ev configuration data
	e.requestConfigurationKeyValueDescriptionListData(entity)

	// get ev identification data
	e.requestIdentitificationlistData(entity)

	// get manufacturer details
	e.requestManufacturer(entity)

	// get electrical connection parameter
	// we ignore this scenario as it is a scoped request and we'll do
	// full requests in the measurements use case

	// get device diagnosis state
	e.requestDeviceDiagnosisState(entity)
}

// request DeviceConfigurationKeyValueDescriptionListData from a remote entity
func (e *EV) requestConfigurationKeyValueDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceConfigurationKeyValueDescriptionListDataType, 1)
	_, err = featureLocal.RequestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, featureRemote, requestChannel)

	if err != nil {
		fmt.Println(err)
		return
	}

	data := <-requestChannel

	fmt.Printf("DescriptionData: %#v\n", data)

	e.requestConfigurationKeyValueListData(entity)
}

// request DeviceConfigurationKeyValueListDataType from a remote entity
func (e *EV) requestConfigurationKeyValueListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceConfiguration, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceConfigurationKeyValueListDataType, 1)
	_, err = featureLocal.RequestData(model.FunctionTypeDeviceConfigurationKeyValueListData, featureRemote, requestChannel)

	if err != nil {
		fmt.Println(err)
	}

	data := <-requestChannel

	fmt.Printf("KeyValueData: %#v\n", data)

	e.updateDeviceConfigurationData(entity)

	// subscribe to device configuration state updates
	_ = entity.Device().Sender().Subscribe(featureLocal.Address(), featureRemote.Address(), model.FeatureTypeTypeDeviceConfiguration)
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
			if descriptionItem.KeyId == dataItem.KeyId {
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
	}

	fmt.Printf("EV Communication Standard: %s\n", evData.CommunicationStandard)
	fmt.Printf("EV Asymmetric Charging Supported: %t\n", evData.AsymmetricChargingSupported)
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

	requestChannel := make(chan *model.IdentificationListDataType, 1)
	_, err = featureLocal.RequestData(model.FunctionTypeIdentificationListData, featureRemote, nil)

	if err != nil {
		fmt.Println(err)
	}

	data := <-requestChannel

	e.updateIdentificationData(entity)

	fmt.Printf("Identification: %#v\n", data)
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

// request EV manufacturer details from a remote entity
func (e *EV) requestManufacturer(entity *spine.EntityRemoteImpl) {
	response := requestManufacturerDetailsForEntity(e.service, entity)
	if response == nil {
		return
	}

	evData := e.dataForRemoteDevice(entity.Device())
	evData.ManufacturerDetails = *response

	fmt.Printf("Power Source: %s\n", evData.ManufacturerDetails.PowerSource)
}

// request DeviceDiagnosisStateData from a remote device
func (e *EV) requestDeviceDiagnosisState(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := requestDeviceDiagnosisStateForEntity(e.service, entity)

	if response == nil {
		return
	}

	// operationState := *response.OperatingState
	// model.DeviceDiagnosisOperatingStateTypeNormalOperation
	// model.DeviceDiagnosisOperatingStateTypeStandby

	// subscribe to entity diagnosis state updates
	_ = entity.Device().Sender().Subscribe(featureLocal.Address(), featureRemote.Address(), model.FeatureTypeTypeDeviceDiagnosis)
}
