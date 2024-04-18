package ucopev

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the min, max, default limits for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCOPEV) CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	return util.GetPhaseCurrentLimits(e.service, entity, e.validEntityTypes)
}

// return the current loadcontrol obligation limits
//
// parameters:
//   - entity: the entity of the EV
//
// return values:
//   - limits: per phase data
//
// possible errors:
//   - ErrDataNotAvailable if no such limit is (yet) available
//   - and others
func (e *UCOPEV) LoadControlLimits(entity spineapi.EntityRemoteInterface) (
	limits []api.LoadLimitsPhase, resultErr error) {
	return util.LoadControlLimits(
		e.service,
		entity,
		e.validEntityTypes,
		model.LoadControlLimitTypeTypeMaxValueLimit,
		model.LoadControlCategoryTypeObligation,
		model.ScopeTypeTypeOverloadProtection)
}

// send new LoadControlLimits to the remote EV
//
// parameters:
//   - limits: a set of limits containing phase specific limit data
//
// Sets a maximum A limit for each phase that the EV may not exceed.
// Mainly used for implementing overload protection of the site or limiting the
// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
// and with ISO15118 if the EV does not support the Optimization of Self Consumption
// usecase.
//
// note:
// For obligations to work for optimizing solar excess power, the EV needs to
// have an energy demand. Recommendations work even if the EV does not have an active
// energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
// In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific
// and needs to have specific EVSE support for the specific EV brand.
// In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
func (e *UCOPEV) WriteLoadControlLimits(entity spineapi.EntityRemoteInterface, limits []api.LoadLimitsPhase) (*model.MsgCounterType, error) {
	return util.WriteLoadControlLimits(e.service, entity, e.validEntityTypes, model.LoadControlCategoryTypeObligation, limits)
}
