package usecases

import (
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EVSECommissioningAndConfiguration struct {
	entity *spine.EntityLocalImpl

	service *service.EEBUSService

	// Delegate EVSEDelegate

	// map connected remote entity to the remote SKI
	remoteEntity map[string]*spine.EntityRemoteImpl

	// map of device SKIs to EVSEData
	data map[string]*EVSEData
}

// Add EVSE support
func NewEVSECommissioningAndConfiguration(service *service.EEBUSService) *EVSECommissioningAndConfiguration {
	// add the use case
	evse := &EVSECommissioningAndConfiguration{
		service: service,
		entity:  service.LocalEntity(),
	}
	spine.Events.Subscribe(evse)

	_ = spine.NewUseCase(
		evse.entity,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2})

	return evse
}

// get the remote device specific data element
func (e *EVSECommissioningAndConfiguration) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *EVSEData {
	if evsedata, ok := e.data[remoteDevice.Ski()]; ok {
		return evsedata
	}

	return &EVSEData{
		OperatingState: model.DeviceDiagnosisOperatingStateTypeNormalOperation,
	}
}
