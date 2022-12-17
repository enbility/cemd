package emobility

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type EmobilityScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*EMobilityImpl

	mux sync.Mutex
}

var _ scenarios.ScenariosI = (*EmobilityScenarioImpl)(nil)

func NewEMobilityScenario(service *service.EEBUSService) *EmobilityScenarioImpl {
	return &EmobilityScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*EMobilityImpl),
	}
}

// adds all the supported features to the local entity
func (e *EmobilityScenarioImpl) AddFeatures() {
	localEntity := e.Service.LocalEntity()

	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient, "Device Diagnosis Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
		f.AddResultHandler(e)
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
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
		f.AddResultHandler(e)
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient, "LoadControl Client")
		f.AddResultHandler(e)
	}
}

// add supported e-mobility usecases
func (e *EmobilityScenarioImpl) AddUseCases() {
	localEntity := e.Service.LocalEntity()

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
		model.SpecificationVersionType("1.0.1b"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeEVStateOfCharge,
		model.SpecificationVersionType("1.0.0"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
		model.SpecificationVersionType("1.0.1b"),
		[]model.UseCaseScenarioSupportType{1, 2, 3})

	_ = spine.NewUseCase(
		localEntity,
		model.UseCaseNameTypeCoordinatedEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
}

func (e *EmobilityScenarioImpl) RegisterRemoteDevice(details *service.ServiceDetails) any {
	// TODO: emobility should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	e.mux.Lock()
	defer e.mux.Unlock()

	if em, ok := e.remoteDevices[details.SKI()]; ok {
		return em
	}

	emobility := NewEMobility(e.Service, details)
	e.remoteDevices[details.SKI()] = emobility
	return emobility
}

func (e *EmobilityScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	delete(e.remoteDevices, remoteDeviceSki)

	return e.Service.UnpairRemoteService(remoteDeviceSki)
}

func (e *EmobilityScenarioImpl) HandleResult(errorMsg spine.ResultMessage) {
	e.mux.Lock()
	defer e.mux.Unlock()

	if errorMsg.DeviceRemote == nil {
		return
	}

	em, ok := e.remoteDevices[errorMsg.DeviceRemote.Ski()]
	if !ok {
		return
	}

	em.HandleResult(errorMsg)
}
