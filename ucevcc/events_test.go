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

	payload.Data = eebusutil.Ptr(model.DeviceClassificationManufacturerDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.ElectricalConnectionParameterDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.ElectricalConnectionPermittedValueSetListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.IdentificationListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCEVCCSuite) Test_evConfigurationDataUpdate() {
	s.sut.evConfigurationDataUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evConfigurationDataUpdate(remoteSki, s.evEntity)

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

	s.sut.evConfigurationDataUpdate(remoteSki, s.evEntity)

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

	s.sut.evConfigurationDataUpdate(remoteSki, s.evEntity)
}

func (s *UCEVCCSuite) Test_evIdentificationDataUpdate() {
	s.sut.evIdentificationDataUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evIdentificationDataUpdate(remoteSki, s.evEntity)

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

	s.sut.evIdentificationDataUpdate(remoteSki, s.evEntity)
}

func (s *UCEVCCSuite) Test_evManufacturerDataUpdate() {
	s.sut.evManufacturerDataUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evManufacturerDataUpdate(remoteSki, s.evEntity)

	data := &model.DeviceClassificationManufacturerDataType{
		BrandName: eebusutil.Ptr(model.DeviceClassificationStringType("test")),
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evManufacturerDataUpdate(remoteSki, s.evEntity)
}

func (s *UCEVCCSuite) Test_evElectricalPermittedValuesUpdate() {
	s.sut.evElectricalPermittedValuesUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evElectricalPermittedValuesUpdate(remoteSki, s.evEntity)

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

	s.sut.evElectricalPermittedValuesUpdate(remoteSki, s.evEntity)

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

	s.sut.evElectricalPermittedValuesUpdate(remoteSki, s.evEntity)
}
