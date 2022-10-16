package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EV) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeEntityChange:
		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			// EV connected
			if !e.checkEntityBeingEV(payload.Entity) {
				return
			}
			e.evConnected(payload.Entity)
		case spine.ElementChangeRemove:
			// EV disconnected
			if !e.checkEntityBeingEV(payload.Entity) {
				return
			}
			fmt.Println("EV DISCONNECTED")
		}
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				// TODO: received diagnosis state

			}
		}
	}
}

// check if the provided entity is an EV
func (e *EV) checkEntityBeingEV(entity *spine.EntityRemoteImpl) bool {
	if entity == nil || entity.EntityType() != model.EntityTypeTypeEV {
		return false
	}
	return true
}