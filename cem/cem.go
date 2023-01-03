package cem

import (
	"github.com/enbility/cemd/emobility"
	"github.com/enbility/cemd/grid"
	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type CemConfiguration struct {
	EmobilityScenarioDisabled bool
	EmobilityConfiguration    emobility.EmobilityConfiguration
	GridScenarioDisabled      bool
	Currency                  model.CurrencyType
}

// Generic CEM implementation
type CemImpl struct {
	service *service.EEBUSService

	emobilityScenario, gridScenario scenarios.ScenariosI
}

func NewCEM(serviceDescription *service.Configuration, serviceHandler service.EEBUSServiceHandler, log logging.Logging) *CemImpl {
	cem := &CemImpl{
		service: service.NewEEBUSService(serviceDescription, serviceHandler),
	}

	cem.service.SetLogging(log)

	return cem
}

// Set up the supported usecases and features
func (h *CemImpl) Setup(configuration CemConfiguration) error {
	if err := h.service.Setup(); err != nil {
		return err
	}

	spine.Events.Subscribe(h)

	// Setup the supported usecases and features
	if !configuration.EmobilityScenarioDisabled {
		h.emobilityScenario = emobility.NewEMobilityScenario(h.service, configuration.Currency, configuration.EmobilityConfiguration)
		h.emobilityScenario.AddFeatures()
		h.emobilityScenario.AddUseCases()
	}

	// Setup the supported usecases and features
	if !configuration.GridScenarioDisabled {
		h.gridScenario = grid.NewGridScenario(h.service)
		h.gridScenario.AddFeatures()
		h.gridScenario.AddUseCases()
	}

	return nil
}

func (h *CemImpl) Start() {
	h.service.Start()
}

func (h *CemImpl) Shutdown() {
	h.service.Shutdown()
}

func (h *CemImpl) RegisterEmobilityRemoteDevice(details *service.ServiceDetails, dataProvider emobility.EmobilityDataProvider) *emobility.EMobilityImpl {
	var impl any

	if dataProvider != nil {
		impl = h.emobilityScenario.RegisterRemoteDevice(details, dataProvider)
	} else {
		impl = h.emobilityScenario.RegisterRemoteDevice(details, nil)
	}

	return impl.(*emobility.EMobilityImpl)
}

func (h *CemImpl) UnRegisterEmobilityRemoteDevice(remoteDeviceSki string) error {
	return h.emobilityScenario.UnRegisterRemoteDevice(remoteDeviceSki)
}

func (h *CemImpl) RegisterGridRemoteDevice(details *service.ServiceDetails) *grid.GridImpl {
	impl := h.gridScenario.RegisterRemoteDevice(details, nil)
	return impl.(*grid.GridImpl)
}

func (h *CemImpl) UnRegisterGridRemoteDevice(remoteDeviceSki string) error {
	return h.gridScenario.UnRegisterRemoteDevice(remoteDeviceSki)
}
