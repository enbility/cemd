package util

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *UtilSuite) Test_IsPayloadForEntityType() {
	payload := spineapi.EventPayload{}
	result := IsPayloadForEntityType(payload, model.EntityTypeTypeEV)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	result = IsPayloadForEntityType(payload, model.EntityTypeTypeEV)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity: s.evEntity,
	}
	result = IsPayloadForEntityType(payload, model.EntityTypeTypeEV)
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

func (s *UtilSuite) Test_IsEntityTypeConnected() {
	payload := spineapi.EventPayload{}
	result := IsEntityTypeConnected(payload, model.EntityTypeTypeEVSE)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeAdd,
	}
	result = IsEntityTypeConnected(payload, model.EntityTypeTypeEVSE)
	assert.Equal(s.T(), true, result)
}

func (s *UtilSuite) Test_IsEntityTypeDisconnected() {
	payload := spineapi.EventPayload{}
	result := IsEntityTypeDisconnected(payload, model.EntityTypeTypeEVSE)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeRemove,
	}
	result = IsEntityTypeDisconnected(payload, model.EntityTypeTypeEVSE)
	assert.Equal(s.T(), true, result)
}
