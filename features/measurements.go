package features

import (
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type MeasurementType struct {
	Timestamp                       time.Time
	CurrentL1, CurrentL2, CurrentL3 float64
	PowerL1, PowerL2, PowerL3       float64
	ChargedEnergy                   float64
	SoC                             float64
}

// request measurement data to properly interpret the corresponding data messages
func RequestMeasurement(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeMeasurement, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request MeasurementDescriptionListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	// request FunctionTypeMeasurementConstraintsListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementConstraintsListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	// request FunctionTypeMeasurementListData from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return errors.New(fErr.String())
	}

	return nil
}
