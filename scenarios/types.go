package scenarios

import (
	"github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
)

// Implemented by *ScenarioImpl, used by CemImpl
type ScenariosI interface {
	RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any
	UnRegisterRemoteDevice(remoteDeviceSki string)
	AddFeatures()
	AddUseCases()
}

type ScenarioImpl struct {
	Service api.EEBUSService
}

func NewScenarioImpl(service api.EEBUSService) *ScenarioImpl {
	return &ScenarioImpl{
		Service: service,
	}
}
