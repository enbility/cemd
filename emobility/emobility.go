package emobility

import (
	"github.com/DerAndereAndi/eebus-go/features"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/service/util"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EmobilityI interface {
	// return the current charge sate of the EV
	EVCurrentChargeState() (EVChargeStateType, error)

	// return the number of ac connected phases of the EV or 0 if it is unknown
	EVConnectedPhases() (uint, error)

	// return the charged energy measurement in Wh of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVChargedEnergy() (float64, error)

	// return the last power measurement for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVPowerPerPhase() ([]float64, error)

	// return the last current measurement for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVCurrentsPerPhase() ([]float64, error)

	// return the min, max, default limits for each phase of the connected EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVCurrentLimits() ([]float64, []float64, []float64, error)

	// send new LoadControlLimits to the remote EV
	//
	// parameters:
	//   - obligations: Overload Protection Limits per phase in A
	//   - recommendations: Self Consumption recommendations per phase in A
	//
	// obligations:
	// Sets a maximum A limit for each phase that the EV may not exceed.
	// Mainly used for implementing overload protection of the site or limiting the
	// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
	// and with ISO15118 if the EV does not support the Optimization of Self Consumption
	// usecase.
	//
	// recommendations:
	// Sets a recommended charge power in A for each phase. This is mainly
	// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
	// The EV either needs to support the Optimization of Self Consumption usecase or
	// the EVSE needs to be able map the recommendations into oligation limits which then
	// works for all EVs communication either via IEC61851 or ISO15118.
	//
	// note:
	// For obligations to work for optimizing solar excess power, the EV needs to
	// have an energy demand. Recommendations work even if the EV does not have an active
	// energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
	// In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific
	// and needs to have specific EVSE support for the specific EV brand.
	// In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
	EVWriteLoadControlLimits(obligations, recommendations []float64) error

	// return the current communication standard type used to communicate between EVSE and EV
	//
	// if an EV is connected via IEC61851, no ISO15118 specific data can be provided!
	// sometimes the connection starts with IEC61851 before it switches
	// to ISO15118, and sometimes it falls back again. so the error return is
	// never absolut for the whole connection time, except if the use case
	// is not supported
	//
	// the values are not constant and can change due to communication problems, bugs, and
	// sometimes communication starts with IEC61851 before it switches to ISO
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - ErrNotSupported if getting the communication standard is not supported
	//   - and others
	EVCommunicationStandard() (EVCommunicationStandardType, error)

	// returns the identification of the currently connected EV or nil if not available
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVIdentification() (string, error)

	// returns if the EVSE and EV combination support optimzation of self consumption
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVOptimizationOfSelfConsumptionSupported() (bool, error)

	// return if the EVSE and EV combination support providing an SoC
	//
	// requires EVSoCSupported to return true
	// only works with a current ISO15118-2 with VAS or ISO15118-20
	// communication between EVSE and EV
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVSoCSupported() (bool, error)

	// return the last known SoC of the connected EV
	//
	// requires EVSoCSupported to return true
	// only works with a current ISO15118-2 with VAS or ISO15118-20
	// communication between EVSE and EV
	//
	// possible errors:
	//   - ErrNotSupported if support for SoC is not possible
	//   - ErrDataNotAvailable if no such measurement is (yet) available
	//   - and others
	EVSoC() (float64, error)

	// returns if the EVSE and EV combination support coordinated charging
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	EVCoordinatedChargingSupported() (bool, error)
}

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

var _ EmobilityI = (*EMobilityImpl)(nil)

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
