package util

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_IsCompatibleEntity() {
	payload := spineapi.EventPayload{}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}
	result := IsCompatibleEntity(payload.Entity, validEntityTypes)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	result = IsCompatibleEntity(payload.Entity, validEntityTypes)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.evEntity,
	}
	result = IsCompatibleEntity(payload.Entity, validEntityTypes)
	assert.Equal(s.T(), true, result)
}

func (s *UtilSuite) Test_IsDeviceDisconnected() {
	payload := spineapi.EventPayload{}
	result := IsDeviceDisconnected(payload)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		EventType:  spineapi.EventTypeDeviceChange,
		ChangeType: spineapi.ElementChangeRemove,
	}
	result = IsDeviceDisconnected(payload)
	assert.Equal(s.T(), true, result)
}

func (s *UtilSuite) Test_IsEntityConnected() {
	payload := spineapi.EventPayload{}
	result := IsEntityConnected(payload)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeAdd,
	}
	result = IsEntityConnected(payload)
	assert.Equal(s.T(), true, result)
}

func (s *UtilSuite) Test_IsEntityDisconnected() {
	payload := spineapi.EventPayload{}
	result := IsEntityDisconnected(payload)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeRemove,
	}
	result = IsEntityDisconnected(payload)
	assert.Equal(s.T(), true, result)
}
