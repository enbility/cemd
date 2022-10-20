package features

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// request identification data to properly interpret the corresponding data messages
func RequestIdentification(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeIdentification, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request FunctionTypeIdentificationListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeIdentificationListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	return nil
}
