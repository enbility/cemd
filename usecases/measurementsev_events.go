package usecases

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go-cem/features"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (m *MeasurementOfElectricityDuringEVCharging) HandleEvent(payload spine.EventPayload) {
	// we only care about events from an EV entity
	if payload.Entity == nil || payload.Entity.EntityType() != model.EntityTypeTypeEV {
		return
	}

	switch payload.EventType {
	case spine.EventTypeDataChange:
		if payload.ChangeType != spine.ElementChangeUpdate {
			return
		}

		switch payload.Data.(type) {
		case *model.ElectricalConnectionDescriptionListDataType:
			_, err := features.GetElectricalDescription(m.service, payload.Entity)
			if err != nil {
				fmt.Println("Error getting electrical description:", err)
				return
			}

			// TODO: provide the electrical description data
		case *model.ElectricalConnectionPermittedValueSetListDataType:
			_, err := features.GetElectricalLimitValues(m.service, payload.Entity)
			if err != nil {
				fmt.Println("Error getting electrical limit values:", err)
				return
			}

			// TODO: provide the electrical limit data
		case *model.MeasurementListDataType:
			_, err := features.GetMeasurementValues(m.service, payload.Entity)
			if err != nil {
				fmt.Println("Error getting measurement values:", err)
				return
			}

			// TODO: provide the measurement data
		}
	}
}
