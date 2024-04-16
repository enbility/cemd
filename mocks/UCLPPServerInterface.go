// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	cemdapi "github.com/enbility/cemd/api"
	api "github.com/enbility/spine-go/api"

	mock "github.com/stretchr/testify/mock"

	model "github.com/enbility/spine-go/model"

	time "time"
)

// UCLPPServerInterface is an autogenerated mock type for the UCLPPServerInterface type
type UCLPPServerInterface struct {
	mock.Mock
}

type UCLPPServerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *UCLPPServerInterface) EXPECT() *UCLPPServerInterface_Expecter {
	return &UCLPPServerInterface_Expecter{mock: &_m.Mock}
}

// AddFeatures provides a mock function with given fields:
func (_m *UCLPPServerInterface) AddFeatures() {
	_m.Called()
}

// UCLPPServerInterface_AddFeatures_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFeatures'
type UCLPPServerInterface_AddFeatures_Call struct {
	*mock.Call
}

// AddFeatures is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) AddFeatures() *UCLPPServerInterface_AddFeatures_Call {
	return &UCLPPServerInterface_AddFeatures_Call{Call: _e.mock.On("AddFeatures")}
}

func (_c *UCLPPServerInterface_AddFeatures_Call) Run(run func()) *UCLPPServerInterface_AddFeatures_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_AddFeatures_Call) Return() *UCLPPServerInterface_AddFeatures_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCLPPServerInterface_AddFeatures_Call) RunAndReturn(run func()) *UCLPPServerInterface_AddFeatures_Call {
	_c.Call.Return(run)
	return _c
}

// AddUseCase provides a mock function with given fields:
func (_m *UCLPPServerInterface) AddUseCase() {
	_m.Called()
}

// UCLPPServerInterface_AddUseCase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUseCase'
type UCLPPServerInterface_AddUseCase_Call struct {
	*mock.Call
}

// AddUseCase is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) AddUseCase() *UCLPPServerInterface_AddUseCase_Call {
	return &UCLPPServerInterface_AddUseCase_Call{Call: _e.mock.On("AddUseCase")}
}

func (_c *UCLPPServerInterface_AddUseCase_Call) Run(run func()) *UCLPPServerInterface_AddUseCase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_AddUseCase_Call) Return() *UCLPPServerInterface_AddUseCase_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCLPPServerInterface_AddUseCase_Call) RunAndReturn(run func()) *UCLPPServerInterface_AddUseCase_Call {
	_c.Call.Return(run)
	return _c
}

// ContractualProductionNominalMax provides a mock function with given fields:
func (_m *UCLPPServerInterface) ContractualProductionNominalMax() (float64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ContractualProductionNominalMax")
	}

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func() (float64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCLPPServerInterface_ContractualProductionNominalMax_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContractualProductionNominalMax'
type UCLPPServerInterface_ContractualProductionNominalMax_Call struct {
	*mock.Call
}

// ContractualProductionNominalMax is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) ContractualProductionNominalMax() *UCLPPServerInterface_ContractualProductionNominalMax_Call {
	return &UCLPPServerInterface_ContractualProductionNominalMax_Call{Call: _e.mock.On("ContractualProductionNominalMax")}
}

func (_c *UCLPPServerInterface_ContractualProductionNominalMax_Call) Run(run func()) *UCLPPServerInterface_ContractualProductionNominalMax_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_ContractualProductionNominalMax_Call) Return(_a0 float64, _a1 error) *UCLPPServerInterface_ContractualProductionNominalMax_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCLPPServerInterface_ContractualProductionNominalMax_Call) RunAndReturn(run func() (float64, error)) *UCLPPServerInterface_ContractualProductionNominalMax_Call {
	_c.Call.Return(run)
	return _c
}

// FailsafeDurationMinimum provides a mock function with given fields:
func (_m *UCLPPServerInterface) FailsafeDurationMinimum() (time.Duration, bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FailsafeDurationMinimum")
	}

	var r0 time.Duration
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func() (time.Duration, bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UCLPPServerInterface_FailsafeDurationMinimum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FailsafeDurationMinimum'
type UCLPPServerInterface_FailsafeDurationMinimum_Call struct {
	*mock.Call
}

// FailsafeDurationMinimum is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) FailsafeDurationMinimum() *UCLPPServerInterface_FailsafeDurationMinimum_Call {
	return &UCLPPServerInterface_FailsafeDurationMinimum_Call{Call: _e.mock.On("FailsafeDurationMinimum")}
}

func (_c *UCLPPServerInterface_FailsafeDurationMinimum_Call) Run(run func()) *UCLPPServerInterface_FailsafeDurationMinimum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_FailsafeDurationMinimum_Call) Return(duration time.Duration, isChangeable bool, resultErr error) *UCLPPServerInterface_FailsafeDurationMinimum_Call {
	_c.Call.Return(duration, isChangeable, resultErr)
	return _c
}

func (_c *UCLPPServerInterface_FailsafeDurationMinimum_Call) RunAndReturn(run func() (time.Duration, bool, error)) *UCLPPServerInterface_FailsafeDurationMinimum_Call {
	_c.Call.Return(run)
	return _c
}

// FailsafeProductionActivePowerLimit provides a mock function with given fields:
func (_m *UCLPPServerInterface) FailsafeProductionActivePowerLimit() (float64, bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FailsafeProductionActivePowerLimit")
	}

	var r0 float64
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func() (float64, bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() float64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func() bool); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FailsafeProductionActivePowerLimit'
type UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call struct {
	*mock.Call
}

// FailsafeProductionActivePowerLimit is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) FailsafeProductionActivePowerLimit() *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call {
	return &UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call{Call: _e.mock.On("FailsafeProductionActivePowerLimit")}
}

func (_c *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call) Run(run func()) *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call) Return(value float64, isChangeable bool, resultErr error) *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call {
	_c.Call.Return(value, isChangeable, resultErr)
	return _c
}

func (_c *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call) RunAndReturn(run func() (float64, bool, error)) *UCLPPServerInterface_FailsafeProductionActivePowerLimit_Call {
	_c.Call.Return(run)
	return _c
}

// IsUseCaseSupported provides a mock function with given fields: remoteEntity
func (_m *UCLPPServerInterface) IsUseCaseSupported(remoteEntity api.EntityRemoteInterface) (bool, error) {
	ret := _m.Called(remoteEntity)

	if len(ret) == 0 {
		panic("no return value specified for IsUseCaseSupported")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (bool, error)); ok {
		return rf(remoteEntity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) bool); ok {
		r0 = rf(remoteEntity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(remoteEntity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCLPPServerInterface_IsUseCaseSupported_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsUseCaseSupported'
type UCLPPServerInterface_IsUseCaseSupported_Call struct {
	*mock.Call
}

// IsUseCaseSupported is a helper method to define mock.On call
//   - remoteEntity api.EntityRemoteInterface
func (_e *UCLPPServerInterface_Expecter) IsUseCaseSupported(remoteEntity interface{}) *UCLPPServerInterface_IsUseCaseSupported_Call {
	return &UCLPPServerInterface_IsUseCaseSupported_Call{Call: _e.mock.On("IsUseCaseSupported", remoteEntity)}
}

func (_c *UCLPPServerInterface_IsUseCaseSupported_Call) Run(run func(remoteEntity api.EntityRemoteInterface)) *UCLPPServerInterface_IsUseCaseSupported_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCLPPServerInterface_IsUseCaseSupported_Call) Return(_a0 bool, _a1 error) *UCLPPServerInterface_IsUseCaseSupported_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCLPPServerInterface_IsUseCaseSupported_Call) RunAndReturn(run func(api.EntityRemoteInterface) (bool, error)) *UCLPPServerInterface_IsUseCaseSupported_Call {
	_c.Call.Return(run)
	return _c
}

// ProductionLimit provides a mock function with given fields:
func (_m *UCLPPServerInterface) ProductionLimit() (cemdapi.LoadLimit, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ProductionLimit")
	}

	var r0 cemdapi.LoadLimit
	var r1 error
	if rf, ok := ret.Get(0).(func() (cemdapi.LoadLimit, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() cemdapi.LoadLimit); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(cemdapi.LoadLimit)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCLPPServerInterface_ProductionLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProductionLimit'
type UCLPPServerInterface_ProductionLimit_Call struct {
	*mock.Call
}

// ProductionLimit is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) ProductionLimit() *UCLPPServerInterface_ProductionLimit_Call {
	return &UCLPPServerInterface_ProductionLimit_Call{Call: _e.mock.On("ProductionLimit")}
}

func (_c *UCLPPServerInterface_ProductionLimit_Call) Run(run func()) *UCLPPServerInterface_ProductionLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_ProductionLimit_Call) Return(_a0 cemdapi.LoadLimit, _a1 error) *UCLPPServerInterface_ProductionLimit_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCLPPServerInterface_ProductionLimit_Call) RunAndReturn(run func() (cemdapi.LoadLimit, error)) *UCLPPServerInterface_ProductionLimit_Call {
	_c.Call.Return(run)
	return _c
}

// SetContractualProductionNominalMax provides a mock function with given fields: value
func (_m *UCLPPServerInterface) SetContractualProductionNominalMax(value float64) error {
	ret := _m.Called(value)

	if len(ret) == 0 {
		panic("no return value specified for SetContractualProductionNominalMax")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(float64) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UCLPPServerInterface_SetContractualProductionNominalMax_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetContractualProductionNominalMax'
type UCLPPServerInterface_SetContractualProductionNominalMax_Call struct {
	*mock.Call
}

// SetContractualProductionNominalMax is a helper method to define mock.On call
//   - value float64
func (_e *UCLPPServerInterface_Expecter) SetContractualProductionNominalMax(value interface{}) *UCLPPServerInterface_SetContractualProductionNominalMax_Call {
	return &UCLPPServerInterface_SetContractualProductionNominalMax_Call{Call: _e.mock.On("SetContractualProductionNominalMax", value)}
}

func (_c *UCLPPServerInterface_SetContractualProductionNominalMax_Call) Run(run func(value float64)) *UCLPPServerInterface_SetContractualProductionNominalMax_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(float64))
	})
	return _c
}

func (_c *UCLPPServerInterface_SetContractualProductionNominalMax_Call) Return(resultErr error) *UCLPPServerInterface_SetContractualProductionNominalMax_Call {
	_c.Call.Return(resultErr)
	return _c
}

func (_c *UCLPPServerInterface_SetContractualProductionNominalMax_Call) RunAndReturn(run func(float64) error) *UCLPPServerInterface_SetContractualProductionNominalMax_Call {
	_c.Call.Return(run)
	return _c
}

// SetFailsafeDurationMinimum provides a mock function with given fields: duration, changeable
func (_m *UCLPPServerInterface) SetFailsafeDurationMinimum(duration time.Duration, changeable bool) error {
	ret := _m.Called(duration, changeable)

	if len(ret) == 0 {
		panic("no return value specified for SetFailsafeDurationMinimum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(time.Duration, bool) error); ok {
		r0 = rf(duration, changeable)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UCLPPServerInterface_SetFailsafeDurationMinimum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFailsafeDurationMinimum'
type UCLPPServerInterface_SetFailsafeDurationMinimum_Call struct {
	*mock.Call
}

// SetFailsafeDurationMinimum is a helper method to define mock.On call
//   - duration time.Duration
//   - changeable bool
func (_e *UCLPPServerInterface_Expecter) SetFailsafeDurationMinimum(duration interface{}, changeable interface{}) *UCLPPServerInterface_SetFailsafeDurationMinimum_Call {
	return &UCLPPServerInterface_SetFailsafeDurationMinimum_Call{Call: _e.mock.On("SetFailsafeDurationMinimum", duration, changeable)}
}

func (_c *UCLPPServerInterface_SetFailsafeDurationMinimum_Call) Run(run func(duration time.Duration, changeable bool)) *UCLPPServerInterface_SetFailsafeDurationMinimum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration), args[1].(bool))
	})
	return _c
}

func (_c *UCLPPServerInterface_SetFailsafeDurationMinimum_Call) Return(resultErr error) *UCLPPServerInterface_SetFailsafeDurationMinimum_Call {
	_c.Call.Return(resultErr)
	return _c
}

func (_c *UCLPPServerInterface_SetFailsafeDurationMinimum_Call) RunAndReturn(run func(time.Duration, bool) error) *UCLPPServerInterface_SetFailsafeDurationMinimum_Call {
	_c.Call.Return(run)
	return _c
}

// SetFailsafeProductionActivePowerLimit provides a mock function with given fields: value, changeable
func (_m *UCLPPServerInterface) SetFailsafeProductionActivePowerLimit(value float64, changeable bool) error {
	ret := _m.Called(value, changeable)

	if len(ret) == 0 {
		panic("no return value specified for SetFailsafeProductionActivePowerLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(float64, bool) error); ok {
		r0 = rf(value, changeable)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFailsafeProductionActivePowerLimit'
type UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call struct {
	*mock.Call
}

// SetFailsafeProductionActivePowerLimit is a helper method to define mock.On call
//   - value float64
//   - changeable bool
func (_e *UCLPPServerInterface_Expecter) SetFailsafeProductionActivePowerLimit(value interface{}, changeable interface{}) *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call {
	return &UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call{Call: _e.mock.On("SetFailsafeProductionActivePowerLimit", value, changeable)}
}

func (_c *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call) Run(run func(value float64, changeable bool)) *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(float64), args[1].(bool))
	})
	return _c
}

func (_c *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call) Return(resultErr error) *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call {
	_c.Call.Return(resultErr)
	return _c
}

func (_c *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call) RunAndReturn(run func(float64, bool) error) *UCLPPServerInterface_SetFailsafeProductionActivePowerLimit_Call {
	_c.Call.Return(run)
	return _c
}

// SetProductionLimit provides a mock function with given fields: limit
func (_m *UCLPPServerInterface) SetProductionLimit(limit cemdapi.LoadLimit) error {
	ret := _m.Called(limit)

	if len(ret) == 0 {
		panic("no return value specified for SetProductionLimit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(cemdapi.LoadLimit) error); ok {
		r0 = rf(limit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UCLPPServerInterface_SetProductionLimit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetProductionLimit'
type UCLPPServerInterface_SetProductionLimit_Call struct {
	*mock.Call
}

// SetProductionLimit is a helper method to define mock.On call
//   - limit cemdapi.LoadLimit
func (_e *UCLPPServerInterface_Expecter) SetProductionLimit(limit interface{}) *UCLPPServerInterface_SetProductionLimit_Call {
	return &UCLPPServerInterface_SetProductionLimit_Call{Call: _e.mock.On("SetProductionLimit", limit)}
}

func (_c *UCLPPServerInterface_SetProductionLimit_Call) Run(run func(limit cemdapi.LoadLimit)) *UCLPPServerInterface_SetProductionLimit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(cemdapi.LoadLimit))
	})
	return _c
}

func (_c *UCLPPServerInterface_SetProductionLimit_Call) Return(resultErr error) *UCLPPServerInterface_SetProductionLimit_Call {
	_c.Call.Return(resultErr)
	return _c
}

func (_c *UCLPPServerInterface_SetProductionLimit_Call) RunAndReturn(run func(cemdapi.LoadLimit) error) *UCLPPServerInterface_SetProductionLimit_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUseCaseAvailability provides a mock function with given fields: available
func (_m *UCLPPServerInterface) UpdateUseCaseAvailability(available bool) {
	_m.Called(available)
}

// UCLPPServerInterface_UpdateUseCaseAvailability_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUseCaseAvailability'
type UCLPPServerInterface_UpdateUseCaseAvailability_Call struct {
	*mock.Call
}

// UpdateUseCaseAvailability is a helper method to define mock.On call
//   - available bool
func (_e *UCLPPServerInterface_Expecter) UpdateUseCaseAvailability(available interface{}) *UCLPPServerInterface_UpdateUseCaseAvailability_Call {
	return &UCLPPServerInterface_UpdateUseCaseAvailability_Call{Call: _e.mock.On("UpdateUseCaseAvailability", available)}
}

func (_c *UCLPPServerInterface_UpdateUseCaseAvailability_Call) Run(run func(available bool)) *UCLPPServerInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *UCLPPServerInterface_UpdateUseCaseAvailability_Call) Return() *UCLPPServerInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCLPPServerInterface_UpdateUseCaseAvailability_Call) RunAndReturn(run func(bool)) *UCLPPServerInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return(run)
	return _c
}

// UseCaseName provides a mock function with given fields:
func (_m *UCLPPServerInterface) UseCaseName() model.UseCaseNameType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for UseCaseName")
	}

	var r0 model.UseCaseNameType
	if rf, ok := ret.Get(0).(func() model.UseCaseNameType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.UseCaseNameType)
	}

	return r0
}

// UCLPPServerInterface_UseCaseName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseCaseName'
type UCLPPServerInterface_UseCaseName_Call struct {
	*mock.Call
}

// UseCaseName is a helper method to define mock.On call
func (_e *UCLPPServerInterface_Expecter) UseCaseName() *UCLPPServerInterface_UseCaseName_Call {
	return &UCLPPServerInterface_UseCaseName_Call{Call: _e.mock.On("UseCaseName")}
}

func (_c *UCLPPServerInterface_UseCaseName_Call) Run(run func()) *UCLPPServerInterface_UseCaseName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCLPPServerInterface_UseCaseName_Call) Return(_a0 model.UseCaseNameType) *UCLPPServerInterface_UseCaseName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UCLPPServerInterface_UseCaseName_Call) RunAndReturn(run func() model.UseCaseNameType) *UCLPPServerInterface_UseCaseName_Call {
	_c.Call.Return(run)
	return _c
}

// NewUCLPPServerInterface creates a new instance of UCLPPServerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUCLPPServerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UCLPPServerInterface {
	mock := &UCLPPServerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}