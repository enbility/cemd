package cem

import (
	"github.com/DerAndereAndi/eebus-go-cem/emobility"
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
)

// Generic CEM implementation
type CemImpl struct {
	service *service.EEBUSService
}

func NewCEM(serviceDescription *service.ServiceDescription, serviceDelegate service.EEBUSServiceDelegate, log logging.Logging) *CemImpl {
	cem := &CemImpl{}

	cem.service = service.NewEEBUSService(serviceDescription, serviceDelegate)
	cem.service.SetLogging(log)

	return cem
}

// Set up the supported usecases and features
func (h *CemImpl) Setup(enableEmobility bool) error {
	if err := h.service.Setup(); err != nil {
		return err
	}

	spine.Events.Subscribe(h)

	// Setup the supported usecases and features
	if enableEmobility {
		emobilityScenario := emobility.NewEMobilityScenario(h.service)
		emobilityScenario.AddFeatures()
		emobilityScenario.AddUseCases()
	}

	h.service.Start()
	// defer h.myService.Shutdown()

	return nil
}

func (h *CemImpl) Start() {
	h.service.Start()
}

func (h *CemImpl) Shutdown() {
	h.service.Shutdown()
}

func (h *CemImpl) RegisterEmobilityRemoteDevice(details service.ServiceDetails) *emobility.EMobilityImpl {
	// TODO: emobility should be stored per remote SKI and
	// only be set for the SKI if the device supports it

	return emobility.NewEMobility(h.service, details)
}

func (h *CemImpl) UnRegisterEmobilityRemoteDevice(remoteDeviceSki string) error {
	return h.service.UnregisterRemoteService(remoteDeviceSki)
}
