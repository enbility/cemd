package democem

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Handle incoming usecase specific events
func (h *DemoCem) deviceEventCB(ski string, device spineapi.DeviceRemoteInterface, event api.EventType) {
}

func (h *DemoCem) entityEventCB(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}
