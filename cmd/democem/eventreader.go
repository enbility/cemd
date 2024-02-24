package democem

import (
	"github.com/enbility/cemd/api"
	spineapi "github.com/enbility/spine-go/api"
)

var _ api.UseCaseEventReaderInterface = (*DemoCem)(nil)

// Handle incomfing usecase specific event
func (h *DemoCem) SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event api.UseCaseEventType) {
}
