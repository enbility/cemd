package ucopev

import (
	"github.com/enbility/cemd/api"
	"github.com/stretchr/testify/assert"
)

func (s *UCOPEVSuite) Test_Public() {
	// The actual tests of the functionality is located in the util package

	_, err := s.sut.LoadControlLimits(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.LoadControlLimits(s.evEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.mockRemoteEntity, []api.LoadLimitsPhase{})
	assert.NotNil(s.T(), err)

	_, err = s.sut.WriteLoadControlLimits(s.evEntity, []api.LoadLimitsPhase{})
	assert.NotNil(s.T(), err)
}
