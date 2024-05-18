package cem

import (
	spineapi "github.com/enbility/spine-go/api"
)

func (s *CemSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Device: s.mockRemoteDevice,
	}
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDeviceChange
	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)
}
