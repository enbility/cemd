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
func RequestMeasurementDescriptionList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// request MeasurementDescriptionListData from a remote entity
	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func RequestMeasurementConstraintsList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementConstraintsListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}

// request FunctionTypeMeasurementListData from a remote entity
func RequestMeasurementList(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*model.MsgCounterType, error) {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msgCounter, fErr := featureLocal.RequestData(model.FunctionTypeMeasurementListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return nil, errors.New(fErr.String())
	}

	return msgCounter, nil
}
