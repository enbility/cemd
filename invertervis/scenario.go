package invertervis

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type InverterVisScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*InverterVisImpl

	mux sync.Mutex
}

var _ scenarios.ScenariosI = (*InverterVisScenarioImpl)(nil)

func NewInverterVisScenario(service *service.EEBUSService) *InverterVisScenarioImpl {
	return &InverterVisScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*InverterVisImpl),
	}
}

// adds all the supported features to the local entity
func (e *InverterVisScenarioImpl) AddFeatures() {
	localEntity := e.Service.LocalEntity()

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}

	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

// add supported inverter usecases
func (e *InverterVisScenarioImpl) AddUseCases() {
	localEntity := e.Service.LocalEntity()

	_ = spine.NewUseCaseWithActor(
		localEntity,
		model.UseCaseActorTypeVisualizationAppliance,
		model.UseCaseNameTypeVisualizationOfAggregatedBatteryData,
		model.SpecificationVersionType("1.0.0 RC1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (e *InverterVisScenarioImpl) RegisterRemoteDevice(details *service.ServiceDetails, dataProvider any) any {
	// TODO: invertervis should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	e.mux.Lock()
	defer e.mux.Unlock()

	if em, ok := e.remoteDevices[details.SKI()]; ok {
		return em
	}

	inverter := NewInverterVis(e.Service, details)
	e.remoteDevices[details.SKI()] = inverter
	return inverter
}

func (e *InverterVisScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) error {
	e.mux.Lock()
	defer e.mux.Unlock()

	delete(e.remoteDevices, remoteDeviceSki)

	return e.Service.UnpairRemoteService(remoteDeviceSki)
}

func (e *InverterVisScenarioImpl) HandleResult(errorMsg spine.ResultMessage) {
	e.mux.Lock()
	defer e.mux.Unlock()

	if errorMsg.DeviceRemote == nil {
		return
	}

	em, ok := e.remoteDevices[errorMsg.DeviceRemote.Ski()]
	if !ok {
		return
	}

	em.HandleResult(errorMsg)
}
