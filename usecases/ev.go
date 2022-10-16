package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Interface for the evCC use case for CEM device
type EVDelegate interface {
	// handle device state updates from the remote EV entity
	HandleEVEntityState(ski string, failure bool)
}

// EV Commissioning and Configuration Use Case implementation
type EV struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	Delegate EVDelegate

	// map of device SKIs to EVData
	data map[string]*EVData
}

// Register the use case and features for handling EVs
// CEM will call this on startup
func NewEVCommissioningAndConfiguration(service *service.EEBUSService) *EV {
	// add the use case
	ev := &EV{
		service: service,
		entity:  service.LocalEntity(),
		data:    make(map[string]*EVData),
	}

	// subscribe to get incoming EV events
	spine.Events.Subscribe(ev)

	// add use cases
	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
		model.SpecificationVersionType("1.0.0"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeEVStateOfCharge,
		model.SpecificationVersionType("1.0.0"),
		[]model.UseCaseScenarioSupportType{1})

	_ = spine.NewUseCase(
		ev.entity,
		model.UseCaseNameTypeCoordinatedEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})

	// add the features
	{
		_ = ev.entity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
	}
	{
		_ = ev.entity.GetOrAddFeature(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
	}
	{
		f := ev.entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)

		// Set the initial state
		state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
		deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
			OperatingState: &state,
		}
		f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	}
	{
		_ = ev.entity.GetOrAddFeature(model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
	}

	return ev
}

// Invoke to remove an EV entity
// Called when an EV was disconnected
func (e *EV) UnregisterEV() {
	// remove the entity
	e.service.RemoveEntity(e.entity)
}

// Invoked when an EV entity was added or removed
func (e *EV) TriggerEntityUpdate() {

}

// an EV was connected, trigger required communication
func (e *EV) evConnected(entity *spine.EntityRemoteImpl) {
	fmt.Println("EV CONNECTED")

	// get ev configuration data
	_, err := features.RequestDeviceConfigurationKeyValueDescriptionList(e.service, entity)
	if err != nil {
		return
	}

	// get manufacturer details
	_, err = features.RequestManufacturerDetailsForEntity(e.service, entity)
	if err != nil {
		return
	}

	// get device diagnosis state
	_, err = features.RequestDiagnosisStateForEntity(e.service, entity)

	// get electrical connection parameter
}
