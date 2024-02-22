package api

import (
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

	// returns if the entity supports the usecase
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	IsUseCaseSupported(remoteEntity spineapi.EntityRemoteInterface) (bool, error)
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

// Implemented by *Solutions, used by Cem
type SolutionInterface interface {
	RegisterRemoteDevice(details *shipapi.ServiceDetails, dataProvider any) any
	UnRegisterRemoteDevice(remoteDeviceSki string)
	AddFeatures()
	AddUseCases()
}
