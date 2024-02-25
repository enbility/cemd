package democem

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

var _ api.EventReaderInterface = (*DemoCem)(nil)

// Handle incoming usecase specific events
func (h *DemoCem) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}
