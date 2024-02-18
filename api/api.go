package api

import (
	"errors"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

//go:generate mockery

// Implemented by CEM
type CemInterface interface {
	// Setup the EEBUS service
	Setup() error

	// Start the EEBUS service
	Start()

	// Shutdown the EEBUS service
	Shutdown()

	// Add a use case implementation
	AddUseCase(usecase UseCaseInterface)
}

// Implemented by each UseCase
type UseCaseInterface interface {
	// provide the usecase name
	UseCaseName() model.UseCaseNameType

	// add the features
	AddFeatures()

	// add the usecase
	AddUseCase()
}

// interface for informing the cem about specific events
// for each supported usecase
//
// UseCaseEventType values can be found in the api definition of each
// supported usecase
//
// implemented by the actual CEM, used by UCEvseCCInterface implementation
type UseCaseEventReaderInterface interface {
	// Inform about a new usecase specific event
	SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event UseCaseEventType)
}

// type for usecase specfic event names
type UseCaseEventType string

var ErrNoEvseEntity = errors.New("entity is not an EVSE")
var ErrNoEvEntity = errors.New("entity is not an EV")

// Implemented by *Solutions, used by Cem
type SolutionInterface interface {
	RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any
	UnRegisterRemoteDevice(remoteDeviceSki string)
	AddFeatures()
	AddUseCases()
}
