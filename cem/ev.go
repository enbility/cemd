package cem

import (
	"errors"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/usecase"
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

// Interface for the evCC use case for CEM device
type EVDelegate interface {
	// handle device state updates from the remote EV entity
	HandleEVEntityState(ski string, failure bool, errorCode string)
}

// EV Commissioning and Configuration Use Case implementation
type EV struct {
	*usecase.UseCaseImpl
	service *service.EEBUSService

	Delegate EVDelegate
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
		UseCaseImpl: usecase.NewUseCase(
			entity,
			model.UseCaseNameTypeEVCommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8}),
		service: service,
	}

	// subscribe to get incoming EV events
	spine.Events.Subscribe(useCase)

	// add the features
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
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
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeLoadControl, model.RoleTypeClient, "LoadControl Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeTimeSeries, model.RoleTypeClient, "TimeSeries Client")
		entity.AddFeature(f)
	}
	{
		f := useCase.EntityFeature(entity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeClient, "IncentiveTable Client")
		entity.AddFeature(f)
	}

	return useCase, nil
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
	case spine.EventTypeEntityChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			// EV connected
		case spine.ElementChangeRemove:
			// EV disconnected
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
				e.Delegate.HandleEVEntityState(payload.Ski, failure, string(*deviceDiagnosisStateData.LastErrorCode))
			}
		}
	}
}
