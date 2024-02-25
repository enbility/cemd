package cem

import (
	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	spineapi "github.com/enbility/spine-go/api"
)

// handle SPINE events
func (h *Cem) HandleEvent(payload spineapi.EventPayload) {

	if util.IsDeviceConnected(payload) {
		h.eventCB(payload.Ski, payload.Device, nil, api.DeviceConnected)
		return
	}

	if util.IsDeviceDisconnected(payload) {
		h.eventCB(payload.Ski, payload.Device, nil, api.DeviceDisconnected)
		return
	}
}
