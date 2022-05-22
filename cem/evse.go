package cem

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type ManufacturerDetails struct {
	BrandName    string
	VendorName   string
	VendorCode   string
	DeviceName   string
	SerialNumber string
	PowerSource  string
}

type EVSEData struct {
	OperatingState      model.DeviceDiagnosisOperatingStateType
	ManufacturerDetails ManufacturerDetails
}

// Delegate Interface for the EVSE
type EVSEDelegate interface {
	// handle device state updates from the remote EVSE device
	HandleEVSEDeviceState(ski string, failure bool)

	// handle device manufacturer data updates from the remote EVSE device
	HandleEVSEDeviceManufacturerData(ski string, details ManufacturerDetails)
}

type EVSE struct {
	*spine.UseCaseImpl

	service *service.EEBUSService

	Delegate EVSEDelegate

	// map of device SKIs to EVData
	data map[string]*EVSEData
}

// Add EVSE support
func AddEVSESupport(service *service.EEBUSService) *EVSE {
	entity := service.LocalEntity()

	// add the use case
	useCase := &EVSE{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeEVSECommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2}),
		service: service,
	}
	spine.Events.Subscribe(useCase)

	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
		entity.AddFeature(f)
	}
	{
		f := service.EntityFeature(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient, "Device Diagnosis Client")
		entity.AddFeature(f)
	}

	return useCase
}

// get the remote device specific data element
func (e *EVSE) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *EVSEData {
	if evsedata, ok := e.data[remoteDevice.Ski()]; ok {
		return evsedata
	}

	return &EVSEData{
		OperatingState:      model.DeviceDiagnosisOperatingStateTypeNormalOperation,
		ManufacturerDetails: ManufacturerDetails{},
	}
}

// Internal EventHandler Interface for the CEM
func (e *EVSE) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			e.requestManufacturer(payload.Device)
			e.requestDeviceDiagnosisState(payload.Device)
		}
	case spine.EventTypeSubscriptionChange:
		switch payload.Data.(type) {
		case model.SubscriptionManagementRequestCallType:
			data := payload.Data.(model.SubscriptionManagementRequestCallType)
			if *data.ServerFeatureType == model.FeatureTypeTypeDeviceDiagnosis {
				remoteDevice := e.service.RemoteDeviceForSki(payload.Ski)
				if remoteDevice == nil {
					fmt.Println("No remote device found for SKI:", payload.Ski)
					return
				}
				switch payload.ChangeType {
				case spine.ElementChangeAdd:
					// start sending heartbeats
					senderAddr := e.Entity.Device().FeatureByTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer).Address()
					rEntity := remoteDevice.Entity([]model.AddressEntityType{1})
					destinationAddr := remoteDevice.FeatureByEntityTypeAndRole(rEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient).Address()
					if senderAddr == nil || destinationAddr == nil {
						fmt.Println("No sender or destination address found for SKI:", payload.Ski)
						return
					}
					remoteDevice := e.service.RemoteDeviceForSki(payload.Ski)
					remoteDevice.StartHeartbeatSend(senderAddr, destinationAddr)
				}
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				if e.Delegate == nil {
					return
				}

				deviceDiagnosisStateData := payload.Data.(model.DeviceDiagnosisStateDataType)
				failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				e.Delegate.HandleEVSEDeviceState(payload.Ski, failure)
			}
		}
	}
}

// request DeviceClassificationManufacturerData from a remote evse device
func (e *EVSE) requestManufacturer(remoteDevice *spine.DeviceRemoteImpl) {
	rEntity := remoteDevice.Entity([]model.AddressEntityType{1})
	response := requestManufacturerDetailsForEntity(e.service, rEntity)
	if response == nil {
		return
	}

	evseData := e.dataForRemoteDevice(remoteDevice)
	evseData.ManufacturerDetails = *response

	if e.Delegate != nil {
		e.Delegate.HandleEVSEDeviceManufacturerData(remoteDevice.Ski(), *response)
	}
}

// request DeviceClassificationManufacturerData from a remote entity
// is re-used in the EV use case
func requestManufacturerDetailsForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) *ManufacturerDetails {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, entity)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	_, _ = featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote, requestChannel)

	// wait for the response
	response := <-requestChannel

	details := &ManufacturerDetails{}

	if response.BrandName != nil {
		details.BrandName = string(*response.BrandName)
	}
	if response.VendorName != nil {
		details.VendorName = string(*response.VendorName)
	}
	if response.VendorCode != nil {
		details.VendorCode = string(*response.VendorCode)
	}
	if response.DeviceName != nil {
		details.DeviceName = string(*response.DeviceName)
	}
	if response.SerialNumber != nil {
		details.SerialNumber = string(*response.SerialNumber)
	}
	if response.PowerSource != nil {
		details.PowerSource = string(*response.PowerSource)
	}

	return details
}

// request DeviceDiagnosisStateData from a remote device
func (e *EVSE) requestDeviceDiagnosisState(remoteDevice *spine.DeviceRemoteImpl) {
	rEntity := remoteDevice.Entity([]model.AddressEntityType{1})

	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, rEntity)
	if err != nil {
		fmt.Println(err)
		return
	}

	response := requestDeviceDiagnosisStateForEntity(e.service, rEntity)

	failure := *response.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
	if e.Delegate != nil {
		e.Delegate.HandleEVSEDeviceState(remoteDevice.Ski(), failure)
	}

	// subscribe to device diagnosis state updates
	_ = remoteDevice.Sender().Subscribe(featureLocal.Address(), featureRemote.Address(), model.FeatureTypeTypeDeviceDiagnosis)
}

// request DeviceDiagnosisStateData from a remote entity
func requestDeviceDiagnosisStateForEntity(service *service.EEBUSService, entity *spine.EntityRemoteImpl) *model.DeviceDiagnosisStateDataType {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, entity)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	requestChannel := make(chan *model.DeviceDiagnosisStateDataType)
	_, _ = featureLocal.RequestData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote, requestChannel)

	// wait for the response
	response := <-requestChannel

	return response
}

/*
// notify remote devices about the new device diagnosis state
func (e *EVSE) notifyDeviceDiagnosisState(operatingState *model.DeviceDiagnosisStateDataType) {
	remoteDevice := e.service.RemoteDeviceOfType(model.DeviceTypeTypeEnergyManagementSystem)
	if remoteDevice == nil {
		return
	}

	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, remoteDevice)
	if err != nil {
		fmt.Println(err)
		return
	}

	featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)

	_, _ = featureLocal.NotifyData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote)
}
*/
