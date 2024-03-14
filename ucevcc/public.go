package ucevcc

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	eebusapi "github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the current charge state of the EV
func (e *UCEVCC) ChargeState(entity spineapi.EntityRemoteInterface) (api.EVChargeStateType, error) {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return api.EVChargeStateTypeUnplugged, nil
	}

	evDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity)
	if err != nil {
		return api.EVChargeStateTypeUnplugged, nil
	}

	diagnosisState, err := evDeviceDiagnosis.GetState()
	if err != nil {
		return api.EVChargeStateTypeUnknown, err
	}

	operatingState := diagnosisState.OperatingState
	if operatingState == nil {
		return api.EVChargeStateTypeUnknown, eebusapi.ErrDataNotAvailable
	}

	switch *operatingState {
	case model.DeviceDiagnosisOperatingStateTypeNormalOperation:
		return api.EVChargeStateTypeActive, nil
	case model.DeviceDiagnosisOperatingStateTypeStandby:
		return api.EVChargeStateTypePaused, nil
	case model.DeviceDiagnosisOperatingStateTypeFailure:
		return api.EVChargeStateTypeError, nil
	case model.DeviceDiagnosisOperatingStateTypeFinished:
		return api.EVChargeStateTypeFinished, nil
	}

	return api.EVChargeStateTypeUnknown, nil
}

// return if an EV is connected
//
// this includes all required features and
// minimal data being available
func (e *UCEVCC) EVConnected(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil || entity.Device() == nil {
		return false
	}

	// getting current charge state should work
	if _, err := e.ChargeState(entity); err != nil {
		return false
	}

	remoteDevice := e.service.LocalDevice().RemoteDeviceForSki(entity.Device().Ski())
	if remoteDevice == nil {
		return false
	}

	// check if the device still has an entity assigned with the provided entities address
	return remoteDevice.Entity(entity.Address().Entity) == entity
}

func (e *UCEVCC) deviceConfigurationValueForKeyName(
	entity spineapi.EntityRemoteInterface,
	keyname model.DeviceConfigurationKeyNameType,
	valueType model.DeviceConfigurationKeyValueTypeType) (any, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evDeviceConfiguration, err := util.DeviceConfiguration(e.service, entity)
	if err != nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	// check if device configuration descriptions has an communication standard key name
	_, err = evDeviceConfiguration.GetDescriptionForKeyName(keyname)
	if err != nil {
		return nil, err
	}

	data, err := evDeviceConfiguration.GetKeyValueForKeyName(keyname, valueType)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	return data, nil
}

// return the current communication standard type used to communicate between EVSE and EV
//
// if an EV is connected via IEC61851, no ISO15118 specific data can be provided!
// sometimes the connection starts with IEC61851 before it switches
// to ISO15118, and sometimes it falls back again. so the error return is
// never absolut for the whole connection time, except if the use case
// is not supported
//
// the values are not constant and can change due to communication problems, bugs, and
// sometimes communication starts with IEC61851 before it switches to ISO
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVCC) CommunicationStandard(entity spineapi.EntityRemoteInterface) (model.DeviceConfigurationKeyValueStringType, error) {
	unknown := UCEVCCCommunicationStandardUnknown

	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return unknown, api.ErrNoCompatibleEntity
	}

	data, err := e.deviceConfigurationValueForKeyName(entity, model.DeviceConfigurationKeyNameTypeCommunicationsStandard, model.DeviceConfigurationKeyValueTypeTypeString)
	if err != nil || data == nil {
		return unknown, eebusapi.ErrDataNotAvailable
	}

	value, ok := data.(*model.DeviceConfigurationKeyValueStringType)
	if !ok || value == nil {
		return unknown, eebusapi.ErrDataNotAvailable
	}

	return *value, nil
}

// return if the EV supports asymmetric charging
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
func (e *UCEVCC) AsymmetricChargingSupport(entity spineapi.EntityRemoteInterface) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	data, err := e.deviceConfigurationValueForKeyName(entity, model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported, model.DeviceConfigurationKeyValueTypeTypeBoolean)
	if err != nil || data == nil {
		return false, eebusapi.ErrDataNotAvailable
	}

	value, ok := data.(*bool)
	if !ok || value == nil {
		return false, eebusapi.ErrDataNotAvailable
	}

	return bool(*value), nil
}

// return the identifications of the currently connected EV or nil if not available
//
// possible errors:
//   - ErrDataNotAvailable if that information is not (yet) available
//   - and others
func (e *UCEVCC) Identifications(entity spineapi.EntityRemoteInterface) ([]api.IdentificationItem, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evIdentification, err := util.Identification(e.service, entity)
	if err != nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	identifications, err := evIdentification.GetValues()
	if err != nil {
		return nil, err
	}

	var ids []api.IdentificationItem
	for _, identification := range identifications {
		item := api.IdentificationItem{}

		typ := identification.IdentificationType
		if typ != nil {
			item.ValueType = *typ
		}

		value := identification.IdentificationValue
		if value != nil {
			item.Value = string(*value)
		}

		ids = append(ids, item)
	}

	return ids, nil
}

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *UCEVCC) ManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	*api.ManufacturerData,
	error,
) {

	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, api.ErrNoCompatibleEntity
	}

	evDeviceClassification, err := util.DeviceClassification(e.service, entity)
	if err != nil {
		return nil, eebusapi.ErrDataNotAvailable
	}

	data, err := evDeviceClassification.GetManufacturerDetails()
	if err != nil {
		return nil, err
	}

	ret := &api.ManufacturerData{
		DeviceName:                     util.Deref((*string)(data.DeviceName)),
		DeviceCode:                     util.Deref((*string)(data.DeviceCode)),
		SerialNumber:                   util.Deref((*string)(data.SerialNumber)),
		SoftwareRevision:               util.Deref((*string)(data.SoftwareRevision)),
		HardwareRevision:               util.Deref((*string)(data.HardwareRevision)),
		VendorName:                     util.Deref((*string)(data.VendorName)),
		VendorCode:                     util.Deref((*string)(data.VendorCode)),
		BrandName:                      util.Deref((*string)(data.BrandName)),
		PowerSource:                    util.Deref((*string)(data.PowerSource)),
		ManufacturerNodeIdentification: util.Deref((*string)(data.ManufacturerNodeIdentification)),
		ManufacturerLabel:              util.Deref((*string)(data.ManufacturerLabel)),
		ManufacturerDescription:        util.Deref((*string)(data.ManufacturerDescription)),
	}

	return ret, nil
}

// return the min, max, default limits for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *UCEVCC) CurrentLimits(entity spineapi.EntityRemoteInterface) ([]float64, []float64, []float64, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return nil, nil, nil, api.ErrNoCompatibleEntity
	}

	evElectricalConnection, err := util.ElectricalConnection(e.service, entity)
	if err != nil {
		return nil, nil, nil, eebusapi.ErrDataNotAvailable
	}

	var resultMin, resultMax, resultDefault []float64

	for _, phaseName := range util.PhaseNameMapping {
		// electricalParameterDescription contains the measured phase for each measurementId
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionForMeasuredPhase(phaseName)
		if err != nil || elParamDesc.ParameterId == nil {
			continue
		}

		dataMin, dataMax, dataDefault, err := evElectricalConnection.GetLimitsForParameterId(*elParamDesc.ParameterId)
		if err != nil {
			continue
		}

		// Min current data should be derived from min power data
		// but as this value is only properly provided via VAS the
		// currrent min values can not be trusted.

		resultMin = append(resultMin, dataMin)
		resultMax = append(resultMax, dataMax)
		resultDefault = append(resultDefault, dataDefault)
	}

	if len(resultMin) == 0 {
		return nil, nil, nil, eebusapi.ErrDataNotAvailable
	}

	return resultMin, resultMax, resultDefault, nil
}

// is the EV in sleep mode
// returns operatingState, lastErrorCode, error
func (e *UCEVCC) IsInSleepMode(
	entity spineapi.EntityRemoteInterface,
) (bool, error) {
	if !util.IsCompatibleEntity(entity, e.validEntityTypes) {
		return false, api.ErrNoCompatibleEntity
	}

	evseDeviceDiagnosis, err := util.DeviceDiagnosis(e.service, entity)
	if err != nil {
		return false, err
	}

	data, err := evseDeviceDiagnosis.GetState()
	if err != nil {
		return false, err
	}

	if data.OperatingState != nil &&
		*data.OperatingState == model.DeviceDiagnosisOperatingStateTypeStandby {
		return true, nil
	}

	return false, nil
}
