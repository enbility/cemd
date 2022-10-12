package cem

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Measurements of OverloadProtection Use Case implementation
type OverloadProtection struct {
	*spine.UseCaseImpl
	service *service.EEBUSService

	data map[string]*MeasurementType
}

// Register the use case and features for measurements
// CEM will call this on startup
func AddOverloadProtectionSupport(service *service.EEBUSService) (*Measurement, error) {
	if service.ServiceDescription.DeviceType != model.DeviceTypeTypeEnergyManagementSystem {
		return nil, errors.New("device type not supported")
	}

	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &Measurement{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
			model.SpecificationVersionType("1.0.1b"),
			[]model.UseCaseScenarioSupportType{1, 2, 3}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	// add the features
	{
		f := entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
		entity.AddFeature(f)
	}
	{
		f := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)

		// Set the initial state
		state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
		deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
			OperatingState: &state,
		}
		f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)

		entity.AddFeature(f)
	}
	{
		f := entity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient, "LoadControl Client")
		entity.AddFeature(f)
	}
	{
		f := entity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
		entity.AddFeature(f)
	}

	return useCase, nil
}

// get the remote device specific data element
func (o *OverloadProtection) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *MeasurementType {
	if evdata, ok := o.data[remoteDevice.Ski()]; ok {
		return evdata
	}

	return &MeasurementType{}
}

// Internal EventHandler Interface for the CEM
func (o *OverloadProtection) HandleEvent(payload spine.EventPayload) {
}

// request ElectricalConnectionParameterDescriptionListDataType from a remote entity
func (o *OverloadProtection) requestElectricalConnectionParameterDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := o.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	o.requestElectricalConnectionDescriptionListData(entity)
}

// request ElectricalConnectionDescriptionListDataType from a remote entity
func (o *OverloadProtection) requestElectricalConnectionDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := o.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	o.requestElectricalConnectionPermittedValueSetListData(entity)
}

// request ElectricalConnectionPermittedValueSetListDataType from a remote entity
func (o *OverloadProtection) requestElectricalConnectionPermittedValueSetListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := o.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	o.updateElectricalConnectionData(entity)
}

// set the new electrical connection data
func (o *OverloadProtection) updateElectricalConnectionData(entity *spine.EntityRemoteImpl) {
	_, featureRemote, err := o.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	paramDescriptionData := featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData).(*model.ElectricalConnectionParameterDescriptionListDataType)
	descriptionData := featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData).(*model.ElectricalConnectionDescriptionListDataType)
	data := featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData).(*model.ElectricalConnectionPermittedValueSetListDataType)
	if descriptionData == nil || data == nil {
		return
	}

	evData := o.dataForRemoteDevice(entity.Device())

	var phases = map[string]uint{
		"a": 1,
		"b": 2,
		"c": 3,
	}

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
				evData.LimitsPower.Min = minValue
				evData.LimitsPower.Max = maxValue
			case descriptionItem.AcMeasuredPhases != nil && hasRange && hasValue:
				// AC Phase Current Limits
				phase, ok := phases[string(*descriptionItem.AcMeasuredPhases)]
				if !ok {
					continue
				}
				limits := EVCurrentLimitType{
					Min:     minValue,
					Max:     maxValue,
					Default: value,
				}

				evData.LimitsPhase[phase] = limits
			}
		}
	}

	// Validate Limits

	// Min current data should be derived from min power data
	// but as this is only properly provided via VAS the currrent min values can not be trusted.
	// Min current for 3-phase should be at least 2.2A, for 1-phase 6.6A

	/*
				if c.clientData.EVData.ConnectedPhases == 1 {
					minCurrent := 6.6
					if c.clientData.EVData.LimitsL1.Min < minCurrent {
						c.clientData.EVData.LimitsL1.Min = minCurrent
					}
				} else if c.clientData.EVData.ConnectedPhases == 3 {
					minCurrent := 2.2
					if c.clientData.EVData.LimitsL1.Min < minCurrent {
						c.clientData.EVData.LimitsL1.Min = minCurrent
					}
					if c.clientData.EVData.LimitsL2.Min < minCurrent {
						c.clientData.EVData.LimitsL2.Min = minCurrent
					}
					if c.clientData.EVData.LimitsL3.Min < minCurrent {
						c.clientData.EVData.LimitsL3.Min = minCurrent
					}
				}
				c.callDataUpdateHandler(EVDataElementUpdateAmperageLimits)
			}
			if powerLimitsUpdated {
				// Min power data is only properly provided via VAS in ISO15118-2!
				// So use the known min limits and calculate a more likely min power
				if c.clientData.EVData.ConnectedPhases == 1 {
					minPower := c.clientData.EVData.LimitsL1.Min * 230
					if c.clientData.EVData.LimitsPower.Min < minPower {
						c.clientData.EVData.LimitsPower.Min = minPower
					}
				} else if c.clientData.EVData.ConnectedPhases == 3 {
					minPower := c.clientData.EVData.LimitsL1.Min*230 + c.clientData.EVData.LimitsL2.Min*230 + c.clientData.EVData.LimitsL3.Min*230
					if c.clientData.EVData.LimitsPower.Min < minPower {
						c.clientData.EVData.LimitsPower.Min = minPower
					}
				}
				c.callDataUpdateHandler(EVDataElementUpdatePowerLimits)
			}
		}
	*/
}
