package cem

import (
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type EVCurrentLimitType struct {
	Min, Max, Default float64
}

type EVPowerLimitType struct {
	Min, Max float64
}

type EVMeasurementsType struct {
	Timestamp                       time.Time
	CurrentL1, CurrentL2, CurrentL3 float64
	PowerL1, PowerL2, PowerL3       float64
	ChargedEnergy                   float64
	SoC                             float64
}

type MeasurementType struct {
	ConnectedPhases uint
	LimitsPhase     map[uint]EVCurrentLimitType
	LimitsPower     EVPowerLimitType
	Measurements    EVMeasurementsType
}

// Measurements of Electricity during EV Charging Use Case implementation
type Measurement struct {
	*spine.UseCaseImpl
	service *service.EEBUSService

	data map[string]*MeasurementType
}

// Register the use case and features for measurements
// CEM will call this on startup
func AddMeasurementSupport(service *service.EEBUSService) (*Measurement, error) {
	if service.ServiceDescription.DeviceType != model.DeviceTypeTypeEnergyManagementSystem {
		return nil, errors.New("device type not supported")
	}

	// A CEM has all the features implemented in the main entity
	entity := service.LocalEntity()

	// add the use case
	useCase := &Measurement{
		UseCaseImpl: spine.NewUseCase(
			entity,
			model.UseCaseNameTypeEVCommissioningAndConfiguration,
			model.SpecificationVersionType("1.0.1"),
			[]model.UseCaseScenarioSupportType{1, 2, 3}),
		service: service,
	}

	// subscribe to get incoming Measurement events
	spine.Events.Subscribe(useCase)

	// add the features
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient, "Measurement Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
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
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeClient, "LoadControl Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient, "Electrical Connection Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeTimeSeries, model.RoleTypeClient, "TimeSeries Client")
	}
	{
		_ = entity.GetOrAddFeature(model.FeatureTypeTypeIncentiveTable, model.RoleTypeClient, "IncentiveTable Client")
	}

	return useCase, nil
}

// get the remote device specific data element
func (m *Measurement) dataForRemoteDevice(remoteDevice *spine.DeviceRemoteImpl) *MeasurementType {
	if evdata, ok := m.data[remoteDevice.Ski()]; ok {
		return evdata
	}

	return &MeasurementType{}
}

// Internal EventHandler Interface for the CEM
func (m *Measurement) HandleEvent(payload spine.EventPayload) {
}

// request ElectricalConnectionParameterDescriptionListDataType from a remote entity
func (m *Measurement) requestElectricalConnectionParameterDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := m.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	m.requestElectricalConnectionDescriptionListData(entity)
}

// request ElectricalConnectionDescriptionListDataType from a remote entity
func (m *Measurement) requestElectricalConnectionDescriptionListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := m.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionDescriptionListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	m.requestElectricalConnectionPermittedValueSetListData(entity)
}

// request ElectricalConnectionPermittedValueSetListDataType from a remote entity
func (m *Measurement) requestElectricalConnectionPermittedValueSetListData(entity *spine.EntityRemoteImpl) {
	featureLocal, featureRemote, err := m.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, fErr := featureLocal.RequestAndFetchData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, featureRemote)
	if fErr != nil {
		fmt.Println(fErr.String())
		return
	}

	m.updateElectricalConnectionData(entity)
}

// set the new electrical connection data
func (m *Measurement) updateElectricalConnectionData(entity *spine.EntityRemoteImpl) {
	_, featureRemote, err := m.service.GetLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeElectricalConnection, entity)
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

	evData := m.dataForRemoteDevice(entity.Device())

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
