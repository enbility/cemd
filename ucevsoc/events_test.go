package ucevsoc

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCEVSOCSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.Entity = s.evEntity
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = eebusutil.Ptr(model.MeasurementListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCEVSOCSuite) Test_evMeasurementDataUpdate() {
	s.sut.evMeasurementDataUpdate(remoteSki, s.mockRemoteEntity)

	s.sut.evMeasurementDataUpdate(remoteSki, s.evEntity)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeStateOfCharge),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeStateOfHealth),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeTravelRange),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evMeasurementDataUpdate(remoteSki, s.evEntity)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(10),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(1)),
				Value:         model.NewScaledNumberType(200),
			},
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(2)),
				Value:         model.NewScaledNumberType(3000),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evMeasurementDataUpdate(remoteSki, s.evEntity)
}
