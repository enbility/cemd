package inverterpvvis

import (
	"sync"

	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type InverterPVVisScenarioImpl struct {
	*api.Solution

	remoteDevices map[string]*InverterPVVis

	mux sync.Mutex
}

var _ api.SolutionInterface = (*InverterPVVisScenarioImpl)(nil)

func NewInverterVisScenario(service eebusapi.ServiceInterface) *InverterPVVisScenarioImpl {
	return &InverterPVVisScenarioImpl{
		Solution:      api.NewSolution(service),
		remoteDevices: make(map[string]*InverterPVVis),
	}
}

// adds all the supported features to the local entity
func (i *InverterPVVisScenarioImpl) AddFeatures() {
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
func (i *InverterPVVisScenarioImpl) AddUseCases() {
	localEntity := i.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeVisualizationAppliance,
		model.UseCaseNameTypeVisualizationOfAggregatedPhotovoltaicData,
		model.SpecificationVersionType("1.0.0"),
		"RC1",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3})
}

func (i *InverterPVVisScenarioImpl) RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any {
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

func (i *InverterPVVisScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	delete(i.remoteDevices, remoteDeviceSki)

	i.Service.RegisterRemoteSKI(remoteDeviceSki, false)
}

func (i *InverterPVVisScenarioImpl) HandleResult(errorMsg spineapi.ResultMessage) {
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
