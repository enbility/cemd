package emobility

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
)

type EmobilityScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*EMobilityImpl

	mux sync.Mutex

	currency      model.CurrencyType
	configuration EmobilityConfiguration
}

var _ scenarios.ScenariosI = (*EmobilityScenarioImpl)(nil)

func NewEMobilityScenario(service *service.EEBUSService, currency model.CurrencyType, configuration EmobilityConfiguration) *EmobilityScenarioImpl {
	return &EmobilityScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*EMobilityImpl),
		currency:      currency,
		configuration: configuration,
	}
}

// adds all the supported features to the local entity
func (e *EmobilityScenarioImpl) AddFeatures() {
	localEntity := e.Service.LocalEntity()

	// server features
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
		f.AddResultHandler(e)
		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)

		// Set the initial state
		deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
			OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		}
		f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	}

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceDiagnosis,
		model.FeatureTypeTypeDeviceClassification,
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
		model.FeatureTypeTypeLoadControl,
		model.FeatureTypeTypeIdentification,
	}

	if !e.configuration.CoordinatedChargingDisabled {
		clientFeatures = append(clientFeatures, model.FeatureTypeTypeTimeSeries)
		clientFeatures = append(clientFeatures, model.FeatureTypeTypeIncentiveTable)
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
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

	if !e.configuration.CoordinatedChargingDisabled {
		_ = spine.NewUseCase(
			localEntity,
			model.UseCaseNameTypeCoordinatedEVCharging,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8})
	}
}

func (e *EmobilityScenarioImpl) RegisterRemoteDevice(details *service.ServiceDetails, dataProvider any) any {
	// TODO: emobility should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	e.mux.Lock()
	defer e.mux.Unlock()

	if em, ok := e.remoteDevices[details.SKI()]; ok {
		return em
	}

	emobility := NewEMobility(e.Service, details, e.currency, dataProvider.(EmobilityDataProvider))
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
