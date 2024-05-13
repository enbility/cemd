package ucopev

import (
	eebusutil "github.com/enbility/eebus-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UCOPEVSuite) Test_Events() {
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
	payload.Data = eebusutil.Ptr(model.ElectricalConnectionPermittedValueSetListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.LoadControlLimitDescriptionListDataType{})
	s.sut.HandleEvent(payload)

	payload.Data = eebusutil.Ptr(model.LoadControlLimitListDataType{})
	s.sut.HandleEvent(payload)
}

func (s *UCOPEVSuite) Test_Failures() {
	s.sut.evConnected(s.mockRemoteEntity)

	s.sut.evLoadControlLimitDescriptionDataUpdate(s.mockRemoteEntity)
}

func (s *UCOPEVSuite) Test_evElectricalPermittedValuesUpdate() {
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
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(0)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeA),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(1)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(1)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeB),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(2)),
				MeasurementId:          eebusutil.Ptr(model.MeasurementIdType(2)),
				AcMeasuredPhases:       eebusutil.Ptr(model.ElectricalConnectionPhaseNameTypeC),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, paramData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evElectricalPermittedValuesUpdate(payload)

	data := &model.ElectricalConnectionPermittedValueSetListDataType{
		ElectricalConnectionPermittedValueSetData: []model.ElectricalConnectionPermittedValueSetDataType{
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(0)),
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(0.1),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(1400),
								Max: model.NewScaledNumberType(11000),
							},
						},
					},
				},
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(1)),
			},
			{
				ElectricalConnectionId: eebusutil.Ptr(model.ElectricalConnectionIdType(0)),
				ParameterId:            eebusutil.Ptr(model.ElectricalConnectionParameterIdType(2)),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evElectricalPermittedValuesUpdate(payload)
}

func (s *UCOPEVSuite) Test_evLoadControlLimitDataUpdate() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}
	s.sut.evLoadControlLimitDataUpdate(payload)

	payload.Entity = s.evEntity
	s.sut.evLoadControlLimitDataUpdate(payload)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:     eebusutil.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
				LimitCategory: eebusutil.Ptr(model.LoadControlCategoryTypeObligation),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeOverloadProtection),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	s.sut.evLoadControlLimitDataUpdate(payload)

	data := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	payload.Data = data

	s.sut.evLoadControlLimitDataUpdate(payload)

	data = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: eebusutil.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
		},
	}

	payload.Data = data

	s.sut.evLoadControlLimitDataUpdate(payload)
}
