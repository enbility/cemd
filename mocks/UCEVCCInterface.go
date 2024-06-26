// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	cemdapi "github.com/enbility/cemd/api"
	api "github.com/enbility/spine-go/api"

	mock "github.com/stretchr/testify/mock"

	model "github.com/enbility/spine-go/model"
)

// UCEVCCInterface is an autogenerated mock type for the UCEVCCInterface type
type UCEVCCInterface struct {
	mock.Mock
}

type UCEVCCInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *UCEVCCInterface) EXPECT() *UCEVCCInterface_Expecter {
	return &UCEVCCInterface_Expecter{mock: &_m.Mock}
}

// AddFeatures provides a mock function with given fields:
func (_m *UCEVCCInterface) AddFeatures() {
	_m.Called()
}

// UCEVCCInterface_AddFeatures_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFeatures'
type UCEVCCInterface_AddFeatures_Call struct {
	*mock.Call
}

// AddFeatures is a helper method to define mock.On call
func (_e *UCEVCCInterface_Expecter) AddFeatures() *UCEVCCInterface_AddFeatures_Call {
	return &UCEVCCInterface_AddFeatures_Call{Call: _e.mock.On("AddFeatures")}
}

func (_c *UCEVCCInterface_AddFeatures_Call) Run(run func()) *UCEVCCInterface_AddFeatures_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCEVCCInterface_AddFeatures_Call) Return() *UCEVCCInterface_AddFeatures_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCEVCCInterface_AddFeatures_Call) RunAndReturn(run func()) *UCEVCCInterface_AddFeatures_Call {
	_c.Call.Return(run)
	return _c
}

// AddUseCase provides a mock function with given fields:
func (_m *UCEVCCInterface) AddUseCase() {
	_m.Called()
}

// UCEVCCInterface_AddUseCase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddUseCase'
type UCEVCCInterface_AddUseCase_Call struct {
	*mock.Call
}

// AddUseCase is a helper method to define mock.On call
func (_e *UCEVCCInterface_Expecter) AddUseCase() *UCEVCCInterface_AddUseCase_Call {
	return &UCEVCCInterface_AddUseCase_Call{Call: _e.mock.On("AddUseCase")}
}

func (_c *UCEVCCInterface_AddUseCase_Call) Run(run func()) *UCEVCCInterface_AddUseCase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCEVCCInterface_AddUseCase_Call) Return() *UCEVCCInterface_AddUseCase_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCEVCCInterface_AddUseCase_Call) RunAndReturn(run func()) *UCEVCCInterface_AddUseCase_Call {
	_c.Call.Return(run)
	return _c
}

// AsymmetricChargingSupport provides a mock function with given fields: entity
func (_m *UCEVCCInterface) AsymmetricChargingSupport(entity api.EntityRemoteInterface) (bool, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for AsymmetricChargingSupport")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (bool, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) bool); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_AsymmetricChargingSupport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AsymmetricChargingSupport'
type UCEVCCInterface_AsymmetricChargingSupport_Call struct {
	*mock.Call
}

// AsymmetricChargingSupport is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) AsymmetricChargingSupport(entity interface{}) *UCEVCCInterface_AsymmetricChargingSupport_Call {
	return &UCEVCCInterface_AsymmetricChargingSupport_Call{Call: _e.mock.On("AsymmetricChargingSupport", entity)}
}

func (_c *UCEVCCInterface_AsymmetricChargingSupport_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_AsymmetricChargingSupport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_AsymmetricChargingSupport_Call) Return(_a0 bool, _a1 error) *UCEVCCInterface_AsymmetricChargingSupport_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_AsymmetricChargingSupport_Call) RunAndReturn(run func(api.EntityRemoteInterface) (bool, error)) *UCEVCCInterface_AsymmetricChargingSupport_Call {
	_c.Call.Return(run)
	return _c
}

// ChargeState provides a mock function with given fields: entity
func (_m *UCEVCCInterface) ChargeState(entity api.EntityRemoteInterface) (cemdapi.EVChargeStateType, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for ChargeState")
	}

	var r0 cemdapi.EVChargeStateType
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (cemdapi.EVChargeStateType, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) cemdapi.EVChargeStateType); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(cemdapi.EVChargeStateType)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_ChargeState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChargeState'
type UCEVCCInterface_ChargeState_Call struct {
	*mock.Call
}

// ChargeState is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) ChargeState(entity interface{}) *UCEVCCInterface_ChargeState_Call {
	return &UCEVCCInterface_ChargeState_Call{Call: _e.mock.On("ChargeState", entity)}
}

func (_c *UCEVCCInterface_ChargeState_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_ChargeState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_ChargeState_Call) Return(_a0 cemdapi.EVChargeStateType, _a1 error) *UCEVCCInterface_ChargeState_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_ChargeState_Call) RunAndReturn(run func(api.EntityRemoteInterface) (cemdapi.EVChargeStateType, error)) *UCEVCCInterface_ChargeState_Call {
	_c.Call.Return(run)
	return _c
}

// ChargingPowerLimits provides a mock function with given fields: entity
func (_m *UCEVCCInterface) ChargingPowerLimits(entity api.EntityRemoteInterface) (float64, float64, float64, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for ChargingPowerLimits")
	}

	var r0 float64
	var r1 float64
	var r2 float64
	var r3 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (float64, float64, float64, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) float64); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) float64); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Get(1).(float64)
	}

	if rf, ok := ret.Get(2).(func(api.EntityRemoteInterface) float64); ok {
		r2 = rf(entity)
	} else {
		r2 = ret.Get(2).(float64)
	}

	if rf, ok := ret.Get(3).(func(api.EntityRemoteInterface) error); ok {
		r3 = rf(entity)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// UCEVCCInterface_ChargingPowerLimits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChargingPowerLimits'
type UCEVCCInterface_ChargingPowerLimits_Call struct {
	*mock.Call
}

// ChargingPowerLimits is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) ChargingPowerLimits(entity interface{}) *UCEVCCInterface_ChargingPowerLimits_Call {
	return &UCEVCCInterface_ChargingPowerLimits_Call{Call: _e.mock.On("ChargingPowerLimits", entity)}
}

func (_c *UCEVCCInterface_ChargingPowerLimits_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_ChargingPowerLimits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_ChargingPowerLimits_Call) Return(_a0 float64, _a1 float64, _a2 float64, _a3 error) *UCEVCCInterface_ChargingPowerLimits_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3)
	return _c
}

func (_c *UCEVCCInterface_ChargingPowerLimits_Call) RunAndReturn(run func(api.EntityRemoteInterface) (float64, float64, float64, error)) *UCEVCCInterface_ChargingPowerLimits_Call {
	_c.Call.Return(run)
	return _c
}

// CommunicationStandard provides a mock function with given fields: entity
func (_m *UCEVCCInterface) CommunicationStandard(entity api.EntityRemoteInterface) (model.DeviceConfigurationKeyValueStringType, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for CommunicationStandard")
	}

	var r0 model.DeviceConfigurationKeyValueStringType
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (model.DeviceConfigurationKeyValueStringType, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) model.DeviceConfigurationKeyValueStringType); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(model.DeviceConfigurationKeyValueStringType)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_CommunicationStandard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommunicationStandard'
type UCEVCCInterface_CommunicationStandard_Call struct {
	*mock.Call
}

// CommunicationStandard is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) CommunicationStandard(entity interface{}) *UCEVCCInterface_CommunicationStandard_Call {
	return &UCEVCCInterface_CommunicationStandard_Call{Call: _e.mock.On("CommunicationStandard", entity)}
}

func (_c *UCEVCCInterface_CommunicationStandard_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_CommunicationStandard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_CommunicationStandard_Call) Return(_a0 model.DeviceConfigurationKeyValueStringType, _a1 error) *UCEVCCInterface_CommunicationStandard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_CommunicationStandard_Call) RunAndReturn(run func(api.EntityRemoteInterface) (model.DeviceConfigurationKeyValueStringType, error)) *UCEVCCInterface_CommunicationStandard_Call {
	_c.Call.Return(run)
	return _c
}

// EVConnected provides a mock function with given fields: entity
func (_m *UCEVCCInterface) EVConnected(entity api.EntityRemoteInterface) bool {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for EVConnected")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) bool); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// UCEVCCInterface_EVConnected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EVConnected'
type UCEVCCInterface_EVConnected_Call struct {
	*mock.Call
}

// EVConnected is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) EVConnected(entity interface{}) *UCEVCCInterface_EVConnected_Call {
	return &UCEVCCInterface_EVConnected_Call{Call: _e.mock.On("EVConnected", entity)}
}

func (_c *UCEVCCInterface_EVConnected_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_EVConnected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_EVConnected_Call) Return(_a0 bool) *UCEVCCInterface_EVConnected_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UCEVCCInterface_EVConnected_Call) RunAndReturn(run func(api.EntityRemoteInterface) bool) *UCEVCCInterface_EVConnected_Call {
	_c.Call.Return(run)
	return _c
}

// Identifications provides a mock function with given fields: entity
func (_m *UCEVCCInterface) Identifications(entity api.EntityRemoteInterface) ([]cemdapi.IdentificationItem, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for Identifications")
	}

	var r0 []cemdapi.IdentificationItem
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) ([]cemdapi.IdentificationItem, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) []cemdapi.IdentificationItem); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]cemdapi.IdentificationItem)
		}
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_Identifications_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Identifications'
type UCEVCCInterface_Identifications_Call struct {
	*mock.Call
}

// Identifications is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) Identifications(entity interface{}) *UCEVCCInterface_Identifications_Call {
	return &UCEVCCInterface_Identifications_Call{Call: _e.mock.On("Identifications", entity)}
}

func (_c *UCEVCCInterface_Identifications_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_Identifications_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_Identifications_Call) Return(_a0 []cemdapi.IdentificationItem, _a1 error) *UCEVCCInterface_Identifications_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_Identifications_Call) RunAndReturn(run func(api.EntityRemoteInterface) ([]cemdapi.IdentificationItem, error)) *UCEVCCInterface_Identifications_Call {
	_c.Call.Return(run)
	return _c
}

// IsInSleepMode provides a mock function with given fields: entity
func (_m *UCEVCCInterface) IsInSleepMode(entity api.EntityRemoteInterface) (bool, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for IsInSleepMode")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (bool, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) bool); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_IsInSleepMode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsInSleepMode'
type UCEVCCInterface_IsInSleepMode_Call struct {
	*mock.Call
}

// IsInSleepMode is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) IsInSleepMode(entity interface{}) *UCEVCCInterface_IsInSleepMode_Call {
	return &UCEVCCInterface_IsInSleepMode_Call{Call: _e.mock.On("IsInSleepMode", entity)}
}

func (_c *UCEVCCInterface_IsInSleepMode_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_IsInSleepMode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_IsInSleepMode_Call) Return(_a0 bool, _a1 error) *UCEVCCInterface_IsInSleepMode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_IsInSleepMode_Call) RunAndReturn(run func(api.EntityRemoteInterface) (bool, error)) *UCEVCCInterface_IsInSleepMode_Call {
	_c.Call.Return(run)
	return _c
}

// IsUseCaseSupported provides a mock function with given fields: remoteEntity
func (_m *UCEVCCInterface) IsUseCaseSupported(remoteEntity api.EntityRemoteInterface) (bool, error) {
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

// UCEVCCInterface_IsUseCaseSupported_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsUseCaseSupported'
type UCEVCCInterface_IsUseCaseSupported_Call struct {
	*mock.Call
}

// IsUseCaseSupported is a helper method to define mock.On call
//   - remoteEntity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) IsUseCaseSupported(remoteEntity interface{}) *UCEVCCInterface_IsUseCaseSupported_Call {
	return &UCEVCCInterface_IsUseCaseSupported_Call{Call: _e.mock.On("IsUseCaseSupported", remoteEntity)}
}

func (_c *UCEVCCInterface_IsUseCaseSupported_Call) Run(run func(remoteEntity api.EntityRemoteInterface)) *UCEVCCInterface_IsUseCaseSupported_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_IsUseCaseSupported_Call) Return(_a0 bool, _a1 error) *UCEVCCInterface_IsUseCaseSupported_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_IsUseCaseSupported_Call) RunAndReturn(run func(api.EntityRemoteInterface) (bool, error)) *UCEVCCInterface_IsUseCaseSupported_Call {
	_c.Call.Return(run)
	return _c
}

// ManufacturerData provides a mock function with given fields: entity
func (_m *UCEVCCInterface) ManufacturerData(entity api.EntityRemoteInterface) (cemdapi.ManufacturerData, error) {
	ret := _m.Called(entity)

	if len(ret) == 0 {
		panic("no return value specified for ManufacturerData")
	}

	var r0 cemdapi.ManufacturerData
	var r1 error
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) (cemdapi.ManufacturerData, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(api.EntityRemoteInterface) cemdapi.ManufacturerData); ok {
		r0 = rf(entity)
	} else {
		r0 = ret.Get(0).(cemdapi.ManufacturerData)
	}

	if rf, ok := ret.Get(1).(func(api.EntityRemoteInterface) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UCEVCCInterface_ManufacturerData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ManufacturerData'
type UCEVCCInterface_ManufacturerData_Call struct {
	*mock.Call
}

// ManufacturerData is a helper method to define mock.On call
//   - entity api.EntityRemoteInterface
func (_e *UCEVCCInterface_Expecter) ManufacturerData(entity interface{}) *UCEVCCInterface_ManufacturerData_Call {
	return &UCEVCCInterface_ManufacturerData_Call{Call: _e.mock.On("ManufacturerData", entity)}
}

func (_c *UCEVCCInterface_ManufacturerData_Call) Run(run func(entity api.EntityRemoteInterface)) *UCEVCCInterface_ManufacturerData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.EntityRemoteInterface))
	})
	return _c
}

func (_c *UCEVCCInterface_ManufacturerData_Call) Return(_a0 cemdapi.ManufacturerData, _a1 error) *UCEVCCInterface_ManufacturerData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UCEVCCInterface_ManufacturerData_Call) RunAndReturn(run func(api.EntityRemoteInterface) (cemdapi.ManufacturerData, error)) *UCEVCCInterface_ManufacturerData_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUseCaseAvailability provides a mock function with given fields: available
func (_m *UCEVCCInterface) UpdateUseCaseAvailability(available bool) {
	_m.Called(available)
}

// UCEVCCInterface_UpdateUseCaseAvailability_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUseCaseAvailability'
type UCEVCCInterface_UpdateUseCaseAvailability_Call struct {
	*mock.Call
}

// UpdateUseCaseAvailability is a helper method to define mock.On call
//   - available bool
func (_e *UCEVCCInterface_Expecter) UpdateUseCaseAvailability(available interface{}) *UCEVCCInterface_UpdateUseCaseAvailability_Call {
	return &UCEVCCInterface_UpdateUseCaseAvailability_Call{Call: _e.mock.On("UpdateUseCaseAvailability", available)}
}

func (_c *UCEVCCInterface_UpdateUseCaseAvailability_Call) Run(run func(available bool)) *UCEVCCInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *UCEVCCInterface_UpdateUseCaseAvailability_Call) Return() *UCEVCCInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return()
	return _c
}

func (_c *UCEVCCInterface_UpdateUseCaseAvailability_Call) RunAndReturn(run func(bool)) *UCEVCCInterface_UpdateUseCaseAvailability_Call {
	_c.Call.Return(run)
	return _c
}

// UseCaseName provides a mock function with given fields:
func (_m *UCEVCCInterface) UseCaseName() model.UseCaseNameType {
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

// UCEVCCInterface_UseCaseName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseCaseName'
type UCEVCCInterface_UseCaseName_Call struct {
	*mock.Call
}

// UseCaseName is a helper method to define mock.On call
func (_e *UCEVCCInterface_Expecter) UseCaseName() *UCEVCCInterface_UseCaseName_Call {
	return &UCEVCCInterface_UseCaseName_Call{Call: _e.mock.On("UseCaseName")}
}

func (_c *UCEVCCInterface_UseCaseName_Call) Run(run func()) *UCEVCCInterface_UseCaseName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *UCEVCCInterface_UseCaseName_Call) Return(_a0 model.UseCaseNameType) *UCEVCCInterface_UseCaseName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UCEVCCInterface_UseCaseName_Call) RunAndReturn(run func() model.UseCaseNameType) *UCEVCCInterface_UseCaseName_Call {
	_c.Call.Return(run)
	return _c
}

// NewUCEVCCInterface creates a new instance of UCEVCCInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUCEVCCInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *UCEVCCInterface {
	mock := &UCEVCCInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
