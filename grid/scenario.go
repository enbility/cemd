package grid

import (
	"sync"

	"github.com/enbility/cemd/scenarios"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type GridScenarioImpl struct {
	*scenarios.ScenarioImpl

	remoteDevices map[string]*GridImpl

	mux sync.Mutex
}

var _ scenarios.ScenariosI = (*GridScenarioImpl)(nil)

func NewGridScenario(service *service.EEBUSService) *GridScenarioImpl {
	return &GridScenarioImpl{
		ScenarioImpl:  scenarios.NewScenarioImpl(service),
		remoteDevices: make(map[string]*GridImpl),
	}
}

// adds all the supported features to the local entity
func (e *GridScenarioImpl) AddFeatures() {
	localEntity := e.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	// client features
	var clientFeatures = []model.FeatureTypeType{
		model.FeatureTypeTypeDeviceConfiguration,
		model.FeatureTypeTypeElectricalConnection,
		model.FeatureTypeTypeMeasurement,
	}
	for _, feature := range clientFeatures {
		f := localEntity.GetOrAddFeature(feature, model.RoleTypeClient)
		f.AddResultHandler(e)
	}
}

// add supported grid usecases
func (e *GridScenarioImpl) AddUseCases() {
	localEntity := e.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	_ = spine.NewUseCaseWithActor(
		localEntity,
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfGridConnectionPoint,
		model.SpecificationVersionType("1.0.0 RC5"),
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7})
}

func (e *GridScenarioImpl) RegisterRemoteDevice(details *service.ServiceDetails, dataProvider any) any {
	// TODO: grid should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	e.mux.Lock()
	defer e.mux.Unlock()

	if em, ok := e.remoteDevices[details.SKI]; ok {
		return em
	}

	grid := NewGrid(e.Service, details)
	e.remoteDevices[details.SKI] = grid
	return grid
}

func (e *GridScenarioImpl) UnRegisterRemoteDevice(remoteDeviceSki string) {
	e.mux.Lock()
	defer e.mux.Unlock()

	delete(e.remoteDevices, remoteDeviceSki)

	e.Service.RegisterRemoteSKI(remoteDeviceSki, false)
}

func (e *GridScenarioImpl) HandleResult(errorMsg spine.ResultMessage) {
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
