package cem

import (
	"github.com/DerAndereAndi/cemd/emobility"
	"github.com/DerAndereAndi/cemd/scenarios"
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
)

// Generic CEM implementation
type CemImpl struct {
	siteConfig        *scenarios.SiteConfig
	service           *service.EEBUSService
	emobilityScenario *emobility.EmobilityScenarioImpl
}

func NewCEM(siteConfig *scenarios.SiteConfig, serviceDescription *service.ServiceDescription, serviceHandler service.EEBUSServiceHandler, log logging.Logging) *CemImpl {
	cem := &CemImpl{
		siteConfig: siteConfig,
		service:    service.NewEEBUSService(serviceDescription, serviceHandler),
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
		h.emobilityScenario = emobility.NewEMobilityScenario(h.siteConfig, h.service)
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

func (h *CemImpl) RegisterEmobilityRemoteDevice(details service.ServiceDetails) *emobility.EMobilityImpl {
	return h.emobilityScenario.RegisterEmobilityRemoteDevice(details)
}

func (h *CemImpl) UnRegisterEmobilityRemoteDevice(remoteDeviceSki string) error {
	return h.emobilityScenario.UnRegisterEmobilityRemoteDevice(remoteDeviceSki)
}
