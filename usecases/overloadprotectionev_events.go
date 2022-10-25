package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (o *OverloadProtectionEV) HandleEvent(payload spine.EventPayload) {
	// we only care about events from an EV entity
	if payload.Entity == nil || payload.Entity.EntityType() != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.LoadControlLimitDescriptionListDataType:
				_, err := features.RequestLoadControlLimitList(o.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
				}

			case *model.LoadControlLimitListDataType:
				data, err := features.GetLoadControlLimitValues(o.service, payload.Entity)
				if err != nil {
					fmt.Println("Error getting loadcontrol limit values:", err)
					return
				}

				// TODO: provide the loadcontrol limit data
				fmt.Printf("Loadcontrol Limits: %#v\n", data)
			}
		}
	}
}
