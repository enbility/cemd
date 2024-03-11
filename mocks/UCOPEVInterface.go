// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	cemdapi "github.com/enbility/cemd/api"
	api "github.com/enbility/spine-go/api"

	mock "github.com/stretchr/testify/mock"

	model "github.com/enbility/spine-go/model"
)

// UCOPEVInterface is an autogenerated mock type for the UCOPEVInterface type
type UCOPEVInterface struct {
	mock.Mock
}

type UCOPEVInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *UCOPEVInterface) EXPECT() *UCOPEVInterface_Expecter {
	return &UCOPEVInterface_Expecter{mock: &_m.Mock}
}

// AddFeatures provides a mock function with given fields:
func (_m *UCOPEVInterface) AddFeatures() {
	_m.Called()
}

// UCOPEVInterface_AddFeatures_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFeatures'
type UCOPEVInterface_AddFeatures_Call struct {
	*mock.Call
}

// AddFeatures is a helper method to define mock.On call
func (_e *UCOPEVInterface_Expecter) AddFeatures() *UCOPEVInterface_AddFeatures_Call {
	return &UCOPEVInterface_AddFeatures_Call{Call: _e.mock.On("AddFeatures")}
}

func (_c *UCOPEVInterface_AddFeatures_Call) Run(run func()) *UCOPEVInterface_AddFeatures_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCOPEVInterface_AddFeatures_Call) Return() *UCOPEVInterface_AddFeatures_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCOPEVInterface_AddFeatures_Call) RunAndReturn(run func()) *UCOPEVInterface_AddFeatures_Call {
	_c.Call.Return(run)
	return _c
}

// AddUseCase provides a mock function with given fields:
func (_m *UCOPEVInterface) AddUseCase() {
	_m.Called()
}

// UCOPEVInterface_AddUseCase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUseCase'
type UCOPEVInterface_AddUseCase_Call struct {
	*mock.Call
}

// AddUseCase is a helper method to define mock.On call
func (_e *UCOPEVInterface_Expecter) AddUseCase() *UCOPEVInterface_AddUseCase_Call {
	return &UCOPEVInterface_AddUseCase_Call{Call: _e.mock.On("AddUseCase")}
}

func (_c *UCOPEVInterface_AddUseCase_Call) Run(run func()) *UCOPEVInterface_AddUseCase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCOPEVInterface_AddUseCase_Call) Return() *UCOPEVInterface_AddUseCase_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCOPEVInterface_AddUseCase_Call) RunAndReturn(run func()) *UCOPEVInterface_AddUseCase_Call {
	_c.Call.Return(run)
	return _c
}

// IsUseCaseSupported provides a mock function with given fields: remoteEntity
func (_m *UCOPEVInterface) IsUseCaseSupported(remoteEntity api.EntityRemoteInterface) (bool, error) {
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

// UCOPEVInterface_IsUseCaseSupported_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsUseCaseSupported'
type UCOPEVInterface_IsUseCaseSupported_Call struct {
	*mock.Call
}

// IsUseCaseSupported is a helper method to define mock.On call
//   - remoteEntity api.EntityRemoteInterface
func (_e *UCOPEVInterface_Expecter) IsUseCaseSupported(remoteEntity interface{}) *UCOPEVInterface_IsUseCaseSupported_Call {
	return &UCOPEVInterface_IsUseCaseSupported_Call{Call: _e.mock.On("IsUseCaseSupported", remoteEntity)}
}

func (_c *UCOPEVInterface_IsUseCaseSupported_Call) Run(run func(remoteEntity api.EntityRemoteInterface)) *UCOPEVInterface_IsUseCaseSupported_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCOPEVInterface_IsUseCaseSupported_Call) Return(_a0 bool, _a1 error) *UCOPEVInterface_IsUseCaseSupported_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCOPEVInterface_IsUseCaseSupported_Call) RunAndReturn(run func(api.EntityRemoteInterface) (bool, error)) *UCOPEVInterface_IsUseCaseSupported_Call {
	_c.Call.Return(run)
	return _c
}

// LoadControlLimits provides a mock function with given fields: entity
func (_m *UCOPEVInterface) LoadControlLimits(entity api.EntityRemoteInterface) ([]cemdapi.LoadLimitsPhase, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for LoadControlLimits")
	}

	var r0 []cemdapi.LoadLimitsPhase
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) ([]cemdapi.LoadLimitsPhase, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) []cemdapi.LoadLimitsPhase); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cemdapi.LoadLimitsPhase)
		}
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCOPEVInterface_LoadControlLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadControlLimits'
type UCOPEVInterface_LoadControlLimits_Call struct {
	*mock.Call
}

// LoadControlLimits is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCOPEVInterface_Expecter) LoadControlLimits(entity interface{}) *UCOPEVInterface_LoadControlLimits_Call {
	return &UCOPEVInterface_LoadControlLimits_Call{Call: _e.mock.On("LoadControlLimits", entity)}
}

func (_c *UCOPEVInterface_LoadControlLimits_Call) Run(run func(entity api.EntityRemoteInterface)) *UCOPEVInterface_LoadControlLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCOPEVInterface_LoadControlLimits_Call) Return(limits []cemdapi.LoadLimitsPhase, resultErr error) *UCOPEVInterface_LoadControlLimits_Call {
	_c.Call.Return(limits, resultErr)
	return _c
}

func (_c *UCOPEVInterface_LoadControlLimits_Call) RunAndReturn(run func(api.EntityRemoteInterface) ([]cemdapi.LoadLimitsPhase, error)) *UCOPEVInterface_LoadControlLimits_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUseCaseAvailability provides a mock function with given fields: available
func (_m *UCOPEVInterface) UpdateUseCaseAvailability(available bool) {
	_m.Called(available)
}

// UCOPEVInterface_UpdateUseCaseAvailability_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUseCaseAvailability'
type UCOPEVInterface_UpdateUseCaseAvailability_Call struct {
	*mock.Call
}

// UpdateUseCaseAvailability is a helper method to define mock.On call
//   - available bool
func (_e *UCOPEVInterface_Expecter) UpdateUseCaseAvailability(available interface{}) *UCOPEVInterface_UpdateUseCaseAvailability_Call {
	return &UCOPEVInterface_UpdateUseCaseAvailability_Call{Call: _e.mock.On("UpdateUseCaseAvailability", available)}
}

func (_c *UCOPEVInterface_UpdateUseCaseAvailability_Call) Run(run func(available bool)) *UCOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *UCOPEVInterface_UpdateUseCaseAvailability_Call) Return() *UCOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCOPEVInterface_UpdateUseCaseAvailability_Call) RunAndReturn(run func(bool)) *UCOPEVInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return(run)
	return _c
}

// UseCaseName provides a mock function with given fields:
func (_m *UCOPEVInterface) UseCaseName() model.UseCaseNameType {
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

// UCOPEVInterface_UseCaseName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseCaseName'
type UCOPEVInterface_UseCaseName_Call struct {
	*mock.Call
}

// UseCaseName is a helper method to define mock.On call
func (_e *UCOPEVInterface_Expecter) UseCaseName() *UCOPEVInterface_UseCaseName_Call {
	return &UCOPEVInterface_UseCaseName_Call{Call: _e.mock.On("UseCaseName")}
}

func (_c *UCOPEVInterface_UseCaseName_Call) Run(run func()) *UCOPEVInterface_UseCaseName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCOPEVInterface_UseCaseName_Call) Return(_a0 model.UseCaseNameType) *UCOPEVInterface_UseCaseName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UCOPEVInterface_UseCaseName_Call) RunAndReturn(run func() model.UseCaseNameType) *UCOPEVInterface_UseCaseName_Call {
	_c.Call.Return(run)
	return _c
}

// WriteLoadControlLimits provides a mock function with given fields: entity, limits
func (_m *UCOPEVInterface) WriteLoadControlLimits(entity api.EntityRemoteInterface, limits []cemdapi.LoadLimitsPhase) (*model.MsgCounterType, error) {
	ret := _m.Called(entity, limits)

	if len(ret) == 0 {
		panic("no return value specified for WriteLoadControlLimits")
	}

	var r0 *model.MsgCounterType
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface, []cemdapi.LoadLimitsPhase) (*model.MsgCounterType, error)); ok {
		return rf(entity, limits)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface, []cemdapi.LoadLimitsPhase) *model.MsgCounterType); ok {
		r0 = rf(entity, limits)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.MsgCounterType)
		}
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface, []cemdapi.LoadLimitsPhase) error); ok {
		r1 = rf(entity, limits)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCOPEVInterface_WriteLoadControlLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteLoadControlLimits'
type UCOPEVInterface_WriteLoadControlLimits_Call struct {
	*mock.Call
}

// WriteLoadControlLimits is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
//   - limits []cemdapi.LoadLimitsPhase
func (_e *UCOPEVInterface_Expecter) WriteLoadControlLimits(entity interface{}, limits interface{}) *UCOPEVInterface_WriteLoadControlLimits_Call {
	return &UCOPEVInterface_WriteLoadControlLimits_Call{Call: _e.mock.On("WriteLoadControlLimits", entity, limits)}
}

func (_c *UCOPEVInterface_WriteLoadControlLimits_Call) Run(run func(entity api.EntityRemoteInterface, limits []cemdapi.LoadLimitsPhase)) *UCOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface), args[1].([]cemdapi.LoadLimitsPhase))
	})
	return _c
}

func (_c *UCOPEVInterface_WriteLoadControlLimits_Call) Return(_a0 *model.MsgCounterType, _a1 error) *UCOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCOPEVInterface_WriteLoadControlLimits_Call) RunAndReturn(run func(api.EntityRemoteInterface, []cemdapi.LoadLimitsPhase) (*model.MsgCounterType, error)) *UCOPEVInterface_WriteLoadControlLimits_Call {
	_c.Call.Return(run)
	return _c
}

// NewUCOPEVInterface creates a new instance of UCOPEVInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUCOPEVInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UCOPEVInterface {
	mock := &UCOPEVInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}