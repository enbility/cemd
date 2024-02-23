package ucoscev

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

// interface for the EVSE Commissioning and Configuration UseCase
type UCOSCEVInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current loadcontrol obligation limits
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such limit is (yet) available
	//   - and others
	LoadControlLimits(entity spineapi.EntityRemoteInterface) ([]float64, error)

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
	WriteLoadControlLimits(entity spineapi.EntityRemoteInterface, limits []api.LoadLimitsPhase) (*model.MsgCounterType, error)

	// Scenario 2

	// this is automatically covered by the SPINE implementation

	// Scenario 3

	// this is covered by the central CEM interface implementation
	// use that one to set the CEM's operation state which will inform all remote devices
}
