package cem

import (
	"github.com/enbility/cemd/emobility"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
)

// Generic CEM implementation
type CemImpl struct {
	service           *service.EEBUSService
	emobilityScenario *emobility.EmobilityScenarioImpl
}

func NewCEM(serviceDescription *service.Configuration, serviceHandler service.EEBUSServiceHandler, log logging.Logging) *CemImpl {
	cem := &CemImpl{
		service: service.NewEEBUSService(serviceDescription, serviceHandler),
	}

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
		h.emobilityScenario = emobility.NewEMobilityScenario(h.service)
		h.emobilityScenario.AddFeatures()
		h.emobilityScenario.AddUseCases()
	}

	return nil
}

func (h *CemImpl) Start() {
	h.service.Start()
}

func (h *CemImpl) Shutdown() {
	h.service.Shutdown()
}

func (h *CemImpl) RegisterEmobilityRemoteDevice(details *service.ServiceDetails) *emobility.EMobilityImpl {
	return h.emobilityScenario.RegisterEmobilityRemoteDevice(details)
}

func (h *CemImpl) UnRegisterEmobilityRemoteDevice(remoteDeviceSki string) error {
	return h.emobilityScenario.UnRegisterEmobilityRemoteDevice(remoteDeviceSki)
}
