package util

import (
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_MeasurementValueForScope() {
	value, err := MeasurementValueForScope(s.service, s.mockRemoteEntity, model.ScopeTypeTypeACPower)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, value)

	value, err = MeasurementValueForScope(s.service, s.evEntity, model.ScopeTypeTypeACPower)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, value)

	descData := &model.MeasurementDescriptionListDataType{
		MeasurementDescriptionData: []model.MeasurementDescriptionDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				ScopeType:     eebusutil.Ptr(model.ScopeTypeTypeACPower),
			},
		},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.evEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	fErr := rFeature.UpdateData(model.FunctionTypeMeasurementDescriptionListData, descData, nil, nil)
	assert.Nil(s.T(), fErr)

	value, err = MeasurementValueForScope(s.service, s.evEntity, model.ScopeTypeTypeACPower)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0.0, value)

	data := &model.MeasurementListDataType{
		MeasurementData: []model.MeasurementDataType{
			{
				MeasurementId: eebusutil.Ptr(model.MeasurementIdType(0)),
				Value:         model.NewScaledNumberType(80),
			},
		},
	}

	fErr = rFeature.UpdateData(model.FunctionTypeMeasurementListData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	value, err = MeasurementValueForScope(s.service, s.evEntity, model.ScopeTypeTypeACPower)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 80.0, value)
}
