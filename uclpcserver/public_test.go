package uclpcserver

import (
	"time"

	"github.com/enbility/cemd/api"
	"github.com/stretchr/testify/assert"
)

func (s *UCLPCServerSuite) Test_LoadControlLimit() {
	limit, err := s.sut.LoadControlLimit()
	assert.Equal(s.T(), 0.0, limit.Value)
	assert.NotNil(s.T(), err)

	newLimit := api.LoadLimit{
		Duration:     time.Duration(time.Hour * 2),
		IsActive:     true,
		IsChangeable: true,
		Value:        16,
	}
	err = s.sut.SetLoadControlLimit(newLimit)
	assert.Nil(s.T(), err)

	limit, err = s.sut.LoadControlLimit()
	assert.Equal(s.T(), 16.0, limit.Value)
	assert.Nil(s.T(), err)
}

func (s *UCLPCServerSuite) Test_FailsafeConsumptionActivePowerLimit() {
	limit, changeable, err := s.sut.FailsafeConsumptionActivePowerLimit()
	assert.Equal(s.T(), 0.0, limit)
	assert.Equal(s.T(), false, changeable)
	assert.NotNil(s.T(), err)

	err = s.sut.SetFailsafeConsumptionActivePowerLimit(10, true)
	assert.Nil(s.T(), err)

	limit, changeable, err = s.sut.FailsafeConsumptionActivePowerLimit()
	assert.Equal(s.T(), 10.0, limit)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)
}

func (s *UCLPCServerSuite) Test_FailsafeDurationMinimum() {
	// The actual tests of the functionality is located in the util package
	duration, changeable, err := s.sut.FailsafeDurationMinimum()
	assert.Equal(s.T(), time.Duration(0), duration)
	assert.Equal(s.T(), false, changeable)
	assert.NotNil(s.T(), err)

	err = s.sut.SetFailsafeDurationMinimum(time.Duration(time.Hour*1), true)
	assert.NotNil(s.T(), err)

	err = s.sut.SetFailsafeDurationMinimum(time.Duration(time.Hour*2), true)
	assert.Nil(s.T(), err)

	duration, changeable, err = s.sut.FailsafeDurationMinimum()
	assert.Equal(s.T(), time.Duration(time.Hour*2), duration)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)
}

func (s *UCLPCServerSuite) Test_PowerConsumptionNominalMax() {
	value, err := s.sut.PowerConsumptionNominalMax()
	assert.Equal(s.T(), 0.0, value)
	assert.NotNil(s.T(), err)

	err = s.sut.SetPowerConsumptionNominalMax(10)
	assert.Nil(s.T(), err)

	value, err = s.sut.PowerConsumptionNominalMax()
	assert.Equal(s.T(), 10.0, value)
	assert.Nil(s.T(), err)
}

func (s *UCLPCServerSuite) Test_ContractualConsumptionNominalMax() {
	value, err := s.sut.ContractualConsumptionNominalMax()
	assert.Equal(s.T(), 0.0, value)
	assert.NotNil(s.T(), err)

	err = s.sut.SetContractualConsumptionNominalMax(10)
	assert.Nil(s.T(), err)

	value, err = s.sut.ContractualConsumptionNominalMax()
	assert.Equal(s.T(), 10.0, value)
	assert.Nil(s.T(), err)
}
