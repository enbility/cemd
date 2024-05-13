package uclpp

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCLPPSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.monitoredEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = eebusutil.Ptr(model.LoadControlLimitDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.LoadControlLimitListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.DeviceConfigurationKeyValueListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCLPPSuite) Test_Failures() {
	s.sut.connected(s.mockRemoteEntity)

	s.sut.configurationDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *UCLPPSuite) Test_loadControlLimitDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.loadControlLimitDataUpdate(payload)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:      eebusutil.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
				LimitCategory:  eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				LimitDirection: eebusutil.Ptr(model.EnergyDirectionTypeProduce),
				ScopeType:      eebusutil.Ptr(model.ScopeTypeTypeActivePowerLimit),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.loadControlLimitDataUpdate(payload)

	data := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	payload.Data = data

	s.sut.loadControlLimitDataUpdate(payload)

	data = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
		},
	}

	payload.Data = data

	s.sut.loadControlLimitDataUpdate(payload)
}

func (s *UCLPPSuite) Test_configurationDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.monitoredEntity,
	}
	s.sut.configurationDataUpdate(payload)

	descData := &model.DeviceConfigurationKeyValueDescriptionListDataType{
		DeviceConfigurationKeyValueDescriptionData: []model.DeviceConfigurationKeyValueDescriptionDataType{
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(1)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
			},
			{
				KeyId:   eebusutil.Ptr(model.DeviceConfigurationKeyIdType(2)),
				KeyName: eebusutil.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.configurationDataUpdate(payload)

	data := &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{},
	}

	payload.Data = data

	s.sut.configurationDataUpdate(payload)

	data = &model.DeviceConfigurationKeyValueListDataType{
		DeviceConfigurationKeyValueData: []model.DeviceConfigurationKeyValueDataType{
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(1)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
			{
				KeyId: eebusutil.Ptr(model.DeviceConfigurationKeyIdType(2)),
				Value: &model.DeviceConfigurationKeyValueValueType{},
			},
		},
	}

	payload.Data = data

	s.sut.configurationDataUpdate(payload)
}
