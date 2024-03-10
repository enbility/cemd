package ucevcc

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCEVCCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)

	var value model.DeviceDiagnosisOperatingStateType
	payload.Data = &value
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.DeviceClassificationManufacturerDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.ElectricalConnectionParameterDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.ElectricalConnectionPermittedValueSetListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.IdentificationListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCEVCCSuite) Test_Failures() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.evConnected(payload)

	s.sut.evConfigurationDescriptionDataUpdate(s.mockRemoteEntity)

	s.sut.evElectricalParamerDescriptionUpdate(s.mockRemoteEntity)
}

func (s *UCEVCCSuite) Test_evConfigurationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evConfigurationDataUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evConfigurationDataUpdate(payload)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
			},
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(1)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evConfigurationDataUpdate(payload)

	data := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(0)),
				Value: eebusutil.Ptr(model.DeviceConfigurationKeyValueValueType{
					String: eebusutil.Ptr(model.DeviceConfigurationKeyValueStringTypeISO151182ED2),
				}),
			},
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: eebusutil.Ptr(model.DeviceConfigurationKeyValueValueType{
					Boolean: eebusutil.Ptr(false),
				}),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evConfigurationDataUpdate(payload)
}

func (s *UCEVCCSuite) Test_evOperatingStateDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evOperatingStateDataUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evOperatingStateDataUpdate(payload)

	data := &model.DeviceDiagnosisStateDataType{
		OperatingState: eebusutil.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evOperatingStateDataUpdate(payload)
}

func (s *UCEVCCSuite) Test_evIdentificationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evIdentificationDataUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evIdentificationDataUpdate(payload)

	data := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:   eebusutil.Ptr(model.IdentificationIdType(0)),
				IdentificationType: eebusutil.Ptr(model.IdentificationTypeTypeEui48),
			},
			{
				IdentificationId:    eebusutil.Ptr(model.IdentificationIdType(1)),
				IdentificationType:  eebusutil.Ptr(model.IdentificationTypeTypeEui48),
				IdentificationValue: eebusutil.Ptr(model.IdentificationValueType("test")),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeIdentification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeIdentificationListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evIdentificationDataUpdate(payload)
}

func (s *UCEVCCSuite) Test_evManufacturerDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evManufacturerDataUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evManufacturerDataUpdate(payload)

	data := &model.DeviceClassificationManufacturerDataType{
		BrandName: eebusutil.Ptr(model.DeviceClassificationStringType("test")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evManufacturerDataUpdate(payload)
}

func (s *UCEVCCSuite) Test_evElectricalPermittedValuesUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evElectricalPermittedValuesUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evElectricalPermittedValuesUpdate(payload)

	paramData := &model.ElectricalConnectionParameterDescriptionListDataType{
		ElectricalConnectionParameterDescriptionData: []model.ElectricalConnectionParameterDescriptionDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				ScopeType:              eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evElectricalPermittedValuesUpdate(payload)

	permData := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, permData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evElectricalPermittedValuesUpdate(payload)
}
