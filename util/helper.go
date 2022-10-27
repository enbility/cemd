package util

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// helper method which adds checking if the feature is available and the operation is allowed
func requestData(featureLocal spine.FeatureLocal, featureRemote *spine.FeatureRemoteImpl, function model.FunctionType) (*model.MsgCounterType, error) {
	fTypes := featureRemote.Operations()
	if _, exists := fTypes[function]; !exists {
		return nil, ErrFunctionNotSupported
	}

	if !fTypes[function].Read {
		return nil, ErrOperationOnFunctionNotSupported
	}

	msgCounter, fErr := featureLocal.RequestData(function, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// check if the given usecase, actor is supported by the remote device
func IsUsecaseSupported(usecase model.UseCaseNameType, actor model.UseCaseActorType, remoteDevice *spine.DeviceRemoteImpl) bool {
	uci := remoteDevice.UseCaseManager().UseCaseInformation()
	for _, element := range uci {
		if *element.Actor != actor {
			continue
		}
		for _, uc := range element.UseCaseSupport {
			if *uc.UseCaseName == usecase {
				return true
			}
		}
	}

	return false
}

// return the remote entity of a given type and device ski
func EntityOfTypeForSki(service *service.EEBUSService, entityType model.EntityTypeType, ski string) (*spine.EntityRemoteImpl, error) {
	rDevice := service.RemoteDeviceForSki(ski)

	entities := rDevice.Entities()
	for _, entity := range entities {
		if entity.EntityType() == entityType {
			return entity, nil
		}
	}

	return nil, ErrEntityNotFound
}
