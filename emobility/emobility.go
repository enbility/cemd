package emobility

import (
	"github.com/DerAndereAndi/eebus-go/features"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EMobilityImpl struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	evseEntity *spine.EntityRemoteImpl
	evEntity   *spine.EntityRemoteImpl

	deviceClassification map[*spine.EntityRemoteImpl]*features.DeviceClassification
	deviceDiagnosis      map[*spine.EntityRemoteImpl]*features.DeviceDiagnosis

	evDeviceConfiguration  *features.DeviceConfiguration
	evElectricalConnection *features.ElectricalConnection
	evMeasurement          *features.Measurement
	evIdentification       *features.Identification
	evLoadControl          *features.LoadControl

	ski string
}

// Add E-Mobility support
func NewEMobility(service *service.EEBUSService, ski string) *EMobilityImpl {
	ski = util.NormalizeSKI(ski)

	emobility := &EMobilityImpl{
		service:              service,
		entity:               service.LocalEntity(),
		ski:                  ski,
		deviceClassification: make(map[*spine.EntityRemoteImpl]*features.DeviceClassification),
		deviceDiagnosis:      make(map[*spine.EntityRemoteImpl]*features.DeviceDiagnosis),
	}
	spine.Events.Subscribe(emobility)

	return emobility
}

// adds all the supported features to the local entity
func AddEmobilityFeatures(service *service.EEBUSService) {
	localEntity := service.LocalEntity()

	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient, "Device Diagnosis Client")
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
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
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient, "LoadControl Client")
	}
}

// add supported e-mobility usecases
func AddEmobilityUseCases(service *service.EEBUSService) {
	localEntity := service.LocalEntity()

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
