package features

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

/*
var phasesMapping = map[string]uint{
	"a": 1,
	"b": 2,
	"c": 3,
}
*/

type CurrentLimitType struct {
	Min, Max, Default float64
}

type PowerLimitType struct {
	Min, Max float64
}

type ElectricalConnectionType struct {
	ConnectedPhases uint
	LimitsPhase     map[uint]CurrentLimitType
	LimitsPower     PowerLimitType
}

// request electrical connection data to properly interpret the corresponding data messages
func RequestElectricalConnection(service *service.EEBUSService, entity *spine.EntityRemoteImpl) error {
	featureLocal, featureRemote, err := service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// request ElectricalConnectionParameterDescriptionListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	// request ElectricalConnectionDescriptionListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionDescriptionListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	// request ElectricalConnectionPermittedValueSetListDataType from a remote entity
	if _, fErr := featureLocal.RequestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, featureRemote); fErr != nil {
		fmt.Println(fErr.String())
		return err
	}

	return nil
}

/*
// set the new electrical connection data
func updateData(featureRemote *spine.FeatureRemoteImpl, entity *spine.EntityRemoteImpl) {
	paramDescriptionData := featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData).(*model.ElectricalConnectionParameterDescriptionListDataType)
	descriptionData := featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData).(*model.ElectricalConnectionDescriptionListDataType)
	data := featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if descriptionData == nil || data == nil {
		return
	}

	deviceData := e.dataForRemoteDevice(entity.Device())

	for _, descriptionItem := range paramDescriptionData.ElectricalConnectionParameterDescriptionData {

		for _, dataItem := range data.ElectricalConnectionPermittedValueSetData {
			if descriptionItem.ParameterId != dataItem.ParameterId {
				continue
			}

			if len(dataItem.PermittedValueSet) == 0 {
				continue
			}

			var value, minValue, maxValue float64
			hasValue := false
			hasRange := false

			for _, item := range dataItem.PermittedValueSet {
				// is a value set
				if item.Value != nil && len(item.Value) > 0 {
					value = item.Value[0].GetValue()
					hasValue = true
				}
				// is a range set
				if item.Range != nil && len(item.Range) > 0 {
					minValue = item.Range[0].Min.GetValue()
					maxValue = item.Range[0].Max.GetValue()
					hasRange = true
				}
			}

			switch {
			// AC Total Power Limits
			case descriptionItem.ScopeType != nil && *descriptionItem.ScopeType == model.ScopeTypeTypeACPowerTotal && hasRange:
				deviceData.LimitsPower.Min = minValue
				deviceData.LimitsPower.Max = maxValue
			case descriptionItem.AcMeasuredPhases != nil && hasRange && hasValue:
				// AC Phase Current Limits
				phase, ok := phasesMapping[string(*descriptionItem.AcMeasuredPhases)]
				if !ok {
					continue
				}
				limits := CurrentLimitType{
					Min:     minValue,
					Max:     maxValue,
					Default: value,
				}

				deviceData.LimitsPhase[phase] = limits
			}
		}
	}

	e.setDataForRemoteDevice(deviceData, entity.Device())
}
*/
