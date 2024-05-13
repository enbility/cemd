package ucvabd

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCVABDSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.batteryEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
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

	payload.Data = eebusutil.Ptr(model.MeasurementDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCVABDSuite) Test_Failures() {
	s.sut.inverterConnected(s.mockRemoteEntity)

	s.sut.inverterMeasurementDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *UCVABDSuite) Test_inverterMeasurementDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.batteryEntity,
	}
	s.sut.inverterMeasurementDataUpdate(payload)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeACPowerTotal),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeCharge),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeDischarge),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(3)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.batteryEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.inverterMeasurementDataUpdate(payload)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(3)),
				Value:         model.NewScaledNumberType(10),
			},
		},
	}

	payload.Data = data

	s.sut.inverterMeasurementDataUpdate(payload)
}
