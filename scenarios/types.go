package scenarios

import (
	"github.com/enbility/eebus-go/api"
)

// Implemented by *ScenarioImpl, used by CemImpl
type ScenariosI interface {
	RegisterRemoteDevice(details *api.ServiceDetails, dataProvider any) any
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
