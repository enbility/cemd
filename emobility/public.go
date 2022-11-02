package emobility

import (
	"github.com/DerAndereAndi/eebus-go-cem/util"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// return the current charge sate of the EV
func (e *EMobilityImpl) EVCurrentChargeState() (EVChargeStateType, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		// no ev entity found means that it is unplugged
		return EVChargeStateTypeUnplugged, nil
	}

	diagnosisState, err := util.GetDeviceDiagnosisState(e.service, evEntity)
	if err != nil {
		return EVChargeStateTypeUnkown, err
	}

	switch diagnosisState.OperatingState {
	case model.DeviceDiagnosisOperatingStateTypeNormalOperation:
		return EVChargeStateTypeActive, nil
	case model.DeviceDiagnosisOperatingStateTypeStandby:
		return EVChargeStateTypePaused, nil
	case model.DeviceDiagnosisOperatingStateTypeFailure:
		return EVChargeStateTypeError, nil
	case model.DeviceDiagnosisOperatingStateTypeFinished:
		return EVChargeStateTypeFinished, nil
	}

	return EVChargeStateTypeUnkown, nil
}

// return the number of ac connected phases of the EV or 0 if it is unknown
func (e *EMobilityImpl) EVConnectedPhases() (uint, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return 0, err
	}

	phases, err := util.GetElectricalConnectedPhases(e.service, evEntity)
	if err != nil {
		return 0, err
	}

	return phases, nil
}

// return the last current measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVCurrents() ([]float64, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return nil, err
	}

	data, err := util.GetMeasurementCurrents(e.service, evEntity)
	if err != nil {
		return nil, err
	}

	phases := []string{"a", "b", "c"}
	var result []float64

	for _, phase := range phases {
		value := 0.0
		if theValue, exists := data[phase]; exists {
			value = theValue
		}
		result = append(result, value)
	}

	return result, nil
}

// return the min, max, default limits for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVCurrentLimits() ([]float64, []float64, []float64, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return nil, nil, nil, err
	}

	dataMin, dataMax, dataDefault, err := util.GetElectricalCurrentsLimits(e.service, evEntity)
	if err != nil {
		return nil, nil, nil, err
	}

	phases := []string{"a", "b", "c"}
	var resultMin, resultMax, resultDefault []float64

	for _, phase := range phases {
		value := 0.0
		if theValue, exists := dataMin[phase]; exists {
			value = theValue
		}
		resultMin = append(resultMin, value)

		value = 0.0
		if theValue, exists := dataMax[phase]; exists {
			value = theValue
		}
		resultMax = append(resultMax, value)

		value = 0.0
		if theValue, exists := dataDefault[phase]; exists {
			value = theValue
		}
		resultDefault = append(resultDefault, value)
	}

	return resultMin, resultMax, resultDefault, nil
}

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
func (e *EMobilityImpl) EVWriteLoadControlLimits(obligations, recommendations []float64) error {

	return nil
}

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
func (e *EMobilityImpl) EVCommunicationStandard() (EVCommunicationStandardType, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}

	// check if device configuration descriptions has an communication standard key name
	support, err := util.GetDeviceConfigurationDescriptionKeyNameSupport(model.DeviceConfigurationKeyNameTypeCommunicationsStandard, e.service, evEntity)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}
	if !support {
		return EVCommunicationStandardTypeUnknown, util.ErrNotSupported
	}

	data, err := util.GetEVCommunicationStandard(e.service, evEntity)
	if err != nil {
		return EVCommunicationStandardTypeUnknown, err
	}

	return EVCommunicationStandardType(*data), err
}

// returns if the EVSE and EV combination support optimzation of self consumption
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EMobilityImpl) EVOptimizationOfSelfConsumptionSupported() (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the Optimization of self consumption usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	// check if loadcontrol limit descriptions contains a recommendation category
	support, err := util.GetLoadControlDescriptionCategorySupport(model.LoadControlCategoryTypeRecommendation, e.service, evEntity)
	if err != nil {
		return false, err
	}
	return support, nil
}

// return if the EVSE and EV combination support providing an SoC
//
// requires EVSoCSupported to return true
// only works with a current ISO15118-2 with VAS or ISO15118-20
// communication between EVSE and EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EMobilityImpl) EVSoCSupported() (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the SoC usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeEVStateOfCharge, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	// check if measurement descriptions has an SoC scope type
	support, err := util.GetMeasurementDescriptionScopeSupport(model.ScopeTypeTypeStateOfCharge, e.service, evEntity)
	if err != nil {
		return false, err
	}

	return support, nil
}

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
func (e *EMobilityImpl) EVSoC() (float64, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return 0.0, err
	}

	// check if the SoC is supported
	support, err := e.EVSoCSupported()
	if err != nil {
		return 0.0, err
	}
	if !support {
		return 0.0, util.ErrNotSupported
	}

	return util.GetMeasurementSoC(e.service, evEntity)
}

// returns if the EVSE and EV combination support coordinated charging
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *EMobilityImpl) EVCoordinatedChargingSupported() (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, e.ski)
	if err != nil {
		return false, err
	}

	// check if the Coordinated charging usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeCoordinatedEVCharging, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	return true, nil
}
