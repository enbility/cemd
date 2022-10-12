package cem

import (
	"fmt"

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
func AddEVSupport(service *service.EEBUSService) *EV {
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
	{
		_ = ev.entity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
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
	e.requestConfigurationKeyValueDescriptionListData(entity)

	// get manufacturer details
	e.requestManufacturer(entity)

	// get electrical connection parameter
	// we ignore this scenario as it is a scoped request and we'll do
	// full requests in the measurements use case

	// get device diagnosis state
	e.requestDeviceDiagnosisState(entity)
}

// request EV manufacturer details from a remote entity
func (e *EV) requestManufacturer(entity *spine.EntityRemoteImpl) {
	response := requestManufacturerDetailsForEntity(e.service, entity)
	if response == nil {
		return
	}

	evData := e.dataForRemoteDevice(entity.Device())
	evData.ManufacturerDetails = *response

	fmt.Printf("Brand: %s\n", evData.ManufacturerDetails.BrandName)
	fmt.Printf("Device: %s\n", evData.ManufacturerDetails.DeviceName)
	fmt.Printf("Power Source: %s\n", evData.ManufacturerDetails.PowerSource)
}

// request DeviceDiagnosisStateData from a remote device
func (e *EV) requestDeviceDiagnosisState(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := e.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := requestDeviceDiagnosisStateForEntity(e.service, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	if response == nil {
		return
	}

	// operationState := *response.OperatingState
	// model.DeviceDiagnosisOperatingStateTypeNormalOperation
	// model.DeviceDiagnosisOperatingStateTypeStandby

	// subscribe to entity diagnosis state updates
	fErr := featureLocal.SubscribeAndWait(featureRemote.Device(), featureRemote.Address())
	if fErr != nil {
		fmt.Println(fErr.String())
	}
}
