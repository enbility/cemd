package cem

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Delegate Interface for the EVSE
type EVSEDelegate interface {
	// handle device state updates from the remote EVSE device
	HandleEVSEDeviceState(ski string, failure bool, errorCode string)
}

type EVSE struct {
	*spine.UseCaseImpl

	service *service.EEBUSService

	Delegate EVSEDelegate
}

// Add EVSE support
func AddEVSESupport(service *service.EEBUSService) EVSE {
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

	return *useCase
}

// Internal EventHandler Interface for the CEM
func (r *EVSE) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		if payload.ChangeType == spine.ElementChangeAdd {
			r.requestManufacturer(payload.Device)
			r.requestDeviceDiagnosisState(payload.Device)
		}
		// case spine.EventTypeSubscriptionChange:
		// 	switch payload.Data.(type) {
		// 	case model.SubscriptionManagementRequestCallType:
		// 		data := payload.Data.(model.SubscriptionManagementRequestCallType)
		// 		if model.FunctionType(*data.ServerFeatureType) == model.FunctionTypeDeviceDiagnosisStateData {
		// 			// only a CEM should send heartbeats to its remotes
		// 			if r.service.ServiceDescription.DeviceType != model.DeviceTypeTypeEnergyManagementSystem {
		// 				return
		// 			}
		// 			remoteDevice := r.service.RemoteDeviceForSki(payload.Ski)
		// 			if remoteDevice == nil {
		// 				fmt.Println("No remote device found for SKI:", payload.Ski)
		// 				return
		// 			}
		// 			switch payload.ChangeType {
		// 			case spine.ElementChangeAdd:
		// 				// start sending heartbeats
		// 				senderAddr := r.Entity.Device().FeatureByTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer).Address()
		// 				destinationAddr := remoteDevice.FeatureByTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient).Address()
		// 				if senderAddr == nil || destinationAddr == nil {
		// 					fmt.Println("No sender or destination address found for SKI:", payload.Ski)
		// 					return
		// 				}
		// 				remoteDevice.StartHearbeat(senderAddr, destinationAddr)
		// 			case spine.ElementChangeRemove:
		// 				// Stop sending heartbeats
		// 				remoteDevice.StopHeartbeat()
		// 			}
		// 		}
		// 	}
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				if r.Delegate == nil {
					return
				}

				deviceDiagnosisStateData := payload.Data.(model.DeviceDiagnosisStateDataType)
				failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				r.Delegate.HandleEVSEDeviceState(payload.Ski, failure, string(*deviceDiagnosisStateData.LastErrorCode))
			}
		}
	}
}

// request DeviceClassificationManufacturerData from a remote device
func (r *EVSE) requestManufacturer(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal, featureRemote, err := r.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, remoteDevice)
	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote, requestChannel)
}

// request DeviceDiagnosisStateData from a remote device
func (r *EVSE) requestDeviceDiagnosisState(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal, featureRemote, err := r.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, remoteDevice)
	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceDiagnosisStateDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote, requestChannel)

	// subscribe to device diagnosis state updates
	remoteDevice.Sender().Subscribe(featureLocal.Address(), featureRemote.Address(), model.FeatureTypeTypeDeviceDiagnosis)
}

// notify remote devices about the new device diagnosis state
func (r *EVSE) notifyDeviceDiagnosisState(operatingState *model.DeviceDiagnosisStateDataType) {
	remoteDevice := r.service.RemoteDeviceOfType(model.DeviceTypeTypeEnergyManagementSystem)
	if remoteDevice == nil {
		return
	}

	featureLocal, featureRemote, err := r.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, remoteDevice)
	if err != nil {
		fmt.Println(err)
		return
	}

	featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)

	featureLocal.NotifyData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote)
}
