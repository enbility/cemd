package scenarios

import (
	"github.com/enbility/eebus-go/service"
)

// Implemented by EmobilityScenarioImpl, used by CemImpl
type ScenariosI interface {
	RegisterRemoteDevice(details *service.ServiceDetails) any
	UnRegisterRemoteDevice(remoteDeviceSki string) error
	AddFeatures()
	AddUseCases()
}

type ScenarioImpl struct {
	Service *service.EEBUSService
}

func NewScenarioImpl(service *service.EEBUSService) *ScenarioImpl {
	return &ScenarioImpl{
		Service: service,
	}
}
