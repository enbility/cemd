package inverterbatteryvis

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type InverterBatteryVisScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*InverterBatteryVisImpl

	mux sync.Mutex
}

var _ scenarios.ScenariosI = (*InverterBatteryVisScenarioImpl)(nil)

func NewInverterVisScenario(service api.EEBUSService) *InverterBatteryVisScenarioImpl {
	return &InverterBatteryVisScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*InverterBatteryVisImpl),
	}
}

// adds all the supported features to the local entity
func (i *InverterBatteryVisScenarioImpl) AddFeatures() {
	localEntity := i.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

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
func (i *InverterBatteryVisScenarioImpl) AddUseCases() {
	localEntity := i.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeVisualizationAppliance,
		model.UseCaseNameTypeVisualizationOfAggregatedBatteryData,
		model.SpecificationVersionType("1.0.0"),
		"RC1",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4})
}

func (i *InverterBatteryVisScenarioImpl) RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any {
	// TODO: invertervis should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	i.mux.Lock()
	defer i.mux.Unlock()

	if em, ok := i.remoteDevices[details.SKI()]; ok {
		return em
	}

	inverter := NewInverterBatteryVis(i.Service, details)
	i.remoteDevices[details.SKI()] = inverter
	return inverter
}

func (i *InverterBatteryVisScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	delete(i.remoteDevices, remoteDeviceSki)

	i.Service.RegisterRemoteSKI(remoteDeviceSki, false)
}

func (i *InverterBatteryVisScenarioImpl) HandleResult(errorMsg spineapi.ResultMessage) {
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
