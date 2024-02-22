package util

import (
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the measurement value for a scope or an error
func MeasurementValueForScope(service eebusapi.ServiceInterface, entity spineapi.EntityRemoteInterface, scope model.ScopeTypeType) (float64, error) {
	evMeasurement, err := Measurement(service, entity)
	if err != nil {
		return 0, features.ErrFunctionNotSupported
	}

	if data, err := evMeasurement.GetDescriptionsForScope(scope); err == nil {
		for _, item := range data {
			if item.MeasurementId == nil {
				continue
			}

			if _, err := evMeasurement.GetValueForMeasurementId(*item.MeasurementId); err != nil {
				continue
			}

		}
	}

	return 0, features.ErrDataNotAvailable
}
