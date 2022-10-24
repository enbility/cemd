package usecases

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (o *OverloadProtection) HandleEvent(payload spine.EventPayload) {
	// we only care about events from an EV entity
	if payload.Entity == nil || payload.Entity.EntityType() != model.EntityTypeTypeEV {
		return
	}
}
