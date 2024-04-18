package ucoscev

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
func (e *UCOSCEV) CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	return util.GetPhaseCurrentLimits(e.service, entity, e.validEntityTypes)
}

// return the current loadcontrol recommendation limits
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
func (e *UCOSCEV) LoadControlLimits(entity spineapi.EntityRemoteInterface) (limits []api.LoadLimitsPhase, resultErr error) {
	return util.LoadControlLimits(
		e.service,
		entity,
		e.validEntityTypes,
		model.LoadControlLimitTypeTypeMaxValueLimit,
		model.LoadControlCategoryTypeRecommendation,
		model.ScopeTypeTypeSelfConsumption)
}

// send new LoadControlLimits to the remote EV
//
// parameters:
//   - limits: a set of limits containing phase specific limit data
//
// recommendations:
// Sets a recommended charge power in A for each phase. This is mainly
// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
// The EV either needs to support the Optimization of Self Consumption usecase or
// the EVSE needs to be able map the recommendations into oligation limits which then
// works for all EVs communication either via IEC61851 or ISO15118.
func (e *UCOSCEV) WriteLoadControlLimits(entity spineapi.EntityRemoteInterface, limits []api.LoadLimitsPhase) (*model.MsgCounterType, error) {
	return util.WriteLoadControlLimits(e.service, entity, e.validEntityTypes, model.LoadControlCategoryTypeRecommendation, limits)
}
