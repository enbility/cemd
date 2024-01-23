package api

import (
	shipapi "github.com/enbility/ship-go/api"
)

// Implemented by *Solutions, used by Cem
type SolutionInterface interface {
	RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any
	UnRegisterRemoteDevice(remoteDeviceSki string)
	AddFeatures()
	AddUseCases()
}
