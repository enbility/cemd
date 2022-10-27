package emobility

import (
	"github.com/DerAndereAndi/eebus-go-cem/util"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// return the number of ac connected phases of the EV or 0 if it is unknown
//
// ski: the SKI of the remote EVSE device
func (e *EMobilityImpl) EVConnectedPhases(ski string) uint {
	return 0
}

// return the last current measurement for each phase of the connected EV
//
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrDataNotAvailable if no such measurement is (yet) available
// - and others
func (e *EMobilityImpl) EVCurrents(ski string) ([]float64, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
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
		if value, exists := data[phase]; exists {
			result = append(result, value)
		} else {
			result = append(result, 0.0)
		}
	}
	if len(result) == 0 {
		return nil, util.ErrDataNotAvailable
	}

	return result, nil
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
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrDataNotAvailable if that information is not (yet) available
// - ErrNotSupported if getting the communication standard is not supported
// - and others
func (e *EMobilityImpl) EVCommunicationStandard(ski string) (EVCommunicationStandardType, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
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
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrDataNotAvailable if that information is not (yet) available
// - and others
func (e *EMobilityImpl) EVOptimizationOfSelfConsumptionSupported(ski string) (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
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
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrDataNotAvailable if no such measurement is (yet) available
// - and others
func (e *EMobilityImpl) EVSoCSupported(ski string) (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
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
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrNotSupported if support for SoC is not possible
// - ErrDataNotAvailable if no such measurement is (yet) available
// - and others
func (e *EMobilityImpl) EVSoC(ski string) (float64, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
	if err != nil {
		return 0.0, err
	}

	// check if the SoC is supported
	support, err := e.EVSoCSupported(ski)
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
// ski: the SKI of the remote EVSE device that has the EV connected
//
// possible errors:
// - ErrDataNotAvailable if that information is not (yet) available
// - and others
func (e *EMobilityImpl) EVCoordinatedChargingSupported(ski string) (bool, error) {
	evEntity, err := util.EntityOfTypeForSki(e.service, model.EntityTypeTypeEV, ski)
	if err != nil {
		return false, err
	}

	// check if the Coordinated charging usecase is supported
	if !util.IsUsecaseSupported(model.UseCaseNameTypeCoordinatedEVCharging, model.UseCaseActorTypeEV, evEntity.Device()) {
		return false, nil
	}

	return true, nil
}
