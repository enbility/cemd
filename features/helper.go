package features

import (
	"errors"
	"fmt"

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
