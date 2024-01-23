package grid

import (
	"sync"

	"github.com/enbility/cemd/api"
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type GridSolution struct {
	*api.Solution

	remoteDevices map[string]*Grid

	mux sync.Mutex
}

var _ api.SolutionInterface = (*GridSolution)(nil)

func NewGridScenario(service eebusapi.ServiceInterface) *GridSolution {
	return &GridSolution{
		Solution:      api.NewSolution(service),
		remoteDevices: make(map[string]*Grid),
	}
}

// adds all the supported features to the local entity
func (e *GridSolution) AddFeatures() {
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
func (e *GridSolution) AddUseCases() {
	localEntity := e.Service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	localEntity.AddUseCaseSupport(
		model.UseCaseActorTypeMonitoringAppliance,
		model.UseCaseNameTypeMonitoringOfGridConnectionPoint,
		model.SpecificationVersionType("1.0.0"),
		"RC5",
		true,
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7})
}

func (e *GridSolution) RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any {
	// TODO: grid should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	e.mux.Lock()
	defer e.mux.Unlock()

	if em, ok := e.remoteDevices[details.SKI()]; ok {
		return em
	}

	grid := NewGrid(e.Service, details)
	e.remoteDevices[details.SKI()] = grid
	return grid
}

func (e *GridSolution) UnRegisterRemoteDevice(remoteDeviceSki string) {
	e.mux.Lock()
	defer e.mux.Unlock()

	delete(e.remoteDevices, remoteDeviceSki)

	e.Service.RegisterRemoteSKI(remoteDeviceSki, false)
}

func (e *GridSolution) HandleResult(errorMsg spineapi.ResultMessage) {
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
