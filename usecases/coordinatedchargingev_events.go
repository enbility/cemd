package usecases

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (c *CoordinatedChargingEV) HandleEvent(payload spine.EventPayload) {
	// we only care about events from an EV entity
	if payload.Entity == nil || payload.Entity.EntityType() != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {

			}
		}
	}
}
