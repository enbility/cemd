package util

import "github.com/stretchr/testify/assert"

func (s *UtilSuite) Test_Features() {
	feature1, err := DeviceClassification(s.service, s.evseEntity)
	assert.Nil(s.T(), feature1)
	assert.NotNil(s.T(), err)

	feature2, err := DeviceConfiguration(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature2)
	assert.NotNil(s.T(), err)

	feature3, err := DeviceDiagnosis(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature3)
	assert.NotNil(s.T(), err)

	feature4, err := DeviceDiagnosisServer(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature4)
	assert.NotNil(s.T(), err)

	feature5, err := ElectricalConnection(s.service, s.evseEntity)
	assert.Nil(s.T(), feature5)
	assert.NotNil(s.T(), err)

	feature6, err := Identification(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature6)
	assert.NotNil(s.T(), err)

	feature7, err := Measurement(s.service, s.evseEntity)
	assert.Nil(s.T(), feature7)
	assert.NotNil(s.T(), err)

	feature8, err := LoadControl(s.service, s.evseEntity)
	assert.Nil(s.T(), feature8)
	assert.NotNil(s.T(), err)

	feature9, err := TimeSeries(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature9)
	assert.NotNil(s.T(), err)

	feature10, err := IncentiveTable(s.service, s.monitoredEntity)
	assert.Nil(s.T(), feature10)
	assert.NotNil(s.T(), err)
}
