package ucevsoc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the last known SoC of the connected EV
//
// only works with a current ISO15118-2 with VAS or ISO15118-20
// communication between EVSE and EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCEVSOC) StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return 0, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := util.Measurement(e.service, entity)
	if err != nil || evMeasurement == nil {
		return 0, err
	}

	data, err := evMeasurement.GetValuesForTypeCommodityScope(model.MeasurementTypeTypePercentage, model.CommodityTypeTypeElectricity, model.ScopeTypeTypeStateOfCharge)
	if err != nil {
		return 0, err
	}

	// we assume there is only one value, nil is already checked
	value := data[0].Value
	if value == nil {
		return 0, features.ErrDataNotAvailable
	}

	return value.GetValue(), nil
}
