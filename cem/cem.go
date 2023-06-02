package cem

import (
	"github.com/enbility/cemd/emobility"
	"github.com/enbility/cemd/grid"
	"github.com/enbility/cemd/inverterbatteryvis"
	"github.com/enbility/cemd/inverterpvvis"
	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Generic CEM implementation
type CemImpl struct {
	Service *service.EEBUSService

	emobilityScenario, gridScenario, inverterBatteryVisScenario, inverterPVVisScenario scenarios.ScenariosI

	Currency model.CurrencyType
}

func NewCEM(serviceDescription *service.Configuration, serviceHandler service.EEBUSServiceHandler, log logging.Logging) *CemImpl {
	cem := &CemImpl{
		Service:  service.NewEEBUSService(serviceDescription, serviceHandler),
		Currency: model.CurrencyTypeEur,
	}

	cem.Service.SetLogging(log)

	return cem
}

// Set up the supported usecases and features
func (h *CemImpl) Setup() error {
	if err := h.Service.Setup(); err != nil {
		return err
	}

	spine.Events.Subscribe(h)

	return nil
}

// Enable the supported usecases and features

func (h *CemImpl) EnableEmobility(configuration emobility.EmobilityConfiguration) {
	h.emobilityScenario = emobility.NewEMobilityScenario(h.Service, h.Currency, configuration)
	h.emobilityScenario.AddFeatures()
	h.emobilityScenario.AddUseCases()
}

func (h *CemImpl) EnableGrid() {
	h.gridScenario = grid.NewGridScenario(h.Service)
	h.gridScenario.AddFeatures()
	h.gridScenario.AddUseCases()
}

func (h *CemImpl) EnableBatteryVisualization() {
	h.inverterBatteryVisScenario = inverterbatteryvis.NewInverterVisScenario(h.Service)
	h.inverterBatteryVisScenario.AddFeatures()
	h.inverterBatteryVisScenario.AddUseCases()
}

func (h *CemImpl) EnablePVVisualization() {
	h.inverterPVVisScenario = inverterpvvis.NewInverterVisScenario(h.Service)
	h.inverterPVVisScenario.AddFeatures()
	h.inverterPVVisScenario.AddUseCases()
}

func (h *CemImpl) Start() {
	h.Service.Start()
}

func (h *CemImpl) Shutdown() {
	h.Service.Shutdown()
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

func (h *CemImpl) RegisterInverterBatteryVisRemoteDevice(details *service.ServiceDetails) *grid.GridImpl {
	impl := h.inverterBatteryVisScenario.RegisterRemoteDevice(details, nil)
	return impl.(*grid.GridImpl)
}

func (h *CemImpl) UnRegisterInverterBatteryVisRemoteDevice(remoteDeviceSki string) error {
	return h.inverterBatteryVisScenario.UnRegisterRemoteDevice(remoteDeviceSki)
}

func (h *CemImpl) RegisterInverterPVVisRemoteDevice(details *service.ServiceDetails) *grid.GridImpl {
	impl := h.inverterPVVisScenario.RegisterRemoteDevice(details, nil)
	return impl.(*grid.GridImpl)
}

func (h *CemImpl) UnRegisterInverterPVVisRemoteDevice(remoteDeviceSki string) error {
	return h.inverterPVVisScenario.UnRegisterRemoteDevice(remoteDeviceSki)
}
