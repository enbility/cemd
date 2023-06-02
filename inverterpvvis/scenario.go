package inverterpvvis

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type InverterPVVisScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*InverterPVVisImpl

	mux sync.Mutex
}

var _ scenarios.ScenariosI = (*InverterPVVisScenarioImpl)(nil)

func NewInverterVisScenario(service *service.EEBUSService) *InverterPVVisScenarioImpl {
	return &InverterPVVisScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*InverterPVVisImpl),
	}
}

// adds all the supported features to the local entity
func (i *InverterPVVisScenarioImpl) AddFeatures() {
	localEntity := i.Service.LocalEntity()

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}

	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(i)
	}
}

// add supported inverter usecases
func (i *InverterPVVisScenarioImpl) AddUseCases() {
	localEntity := i.Service.LocalEntity()

	_ = spine.NewUseCaseWithActor(
		localEntity,
		model.UseCaseActorTypeVisualizationAppliance,
		model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData,
		model.SpecificationVersionType("1.0.0 RC1"),
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}

func (i *InverterPVVisScenarioImpl) RegisterRemoteDevice(details *service.ServiceDetails, dataProvider any) any {
	// TODO: invertervis should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	i.mux.Lock()
	defer i.mux.Unlock()

	if em, ok := i.remoteDevices[details.SKI()]; ok {
		return em
	}

	inverter := NewInverterPVVis(i.Service, details)
	i.remoteDevices[details.SKI()] = inverter
	return inverter
}

func (i *InverterPVVisScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	delete(i.remoteDevices, remoteDeviceSki)

	return i.Service.UnpairRemoteService(remoteDeviceSki)
}

func (i *InverterPVVisScenarioImpl) HandleResult(errorMsg spine.ResultMessage) {
	i.mux.Lock()
	defer i.mux.Unlock()

	if errorMsg.DeviceRemote == nil {
		return
	}

	em, ok := i.remoteDevices[errorMsg.DeviceRemote.Ski()]
	if !ok {
		return
	}

	em.HandleResult(errorMsg)
}
