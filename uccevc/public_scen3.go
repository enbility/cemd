package uccevc

import (
	"errors"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/util"
	"github.com/enbility/eebus-go/features"
	eebusutil "github.com/enbility/eebus-go/util"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// returns the minimum and maximum number of incentive slots allowed
func (e *UCCEVC) IncentiveConstraints(entity spineapi.EntityRemoteInterface) (api.IncentiveSlotConstraints, error) {
	result := api.IncentiveSlotConstraints{}

	if !e.isCompatibleEntity(entity) {
		return result, api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := util.IncentiveTable(e.service, entity)
	if err != nil {
		return result, features.ErrDataNotAvailable
	}

	constraints, err := evIncentiveTable.GetConstraints()
	if err != nil {
		return result, err
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.IncentiveSlotConstraints.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMin)
	}
	if constraint.IncentiveSlotConstraints.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMax)
	}

	return result, nil
}

// inform the EVSE about used currency and boundary units
//
// SPINE UC CoordinatedEVCharging 2.4.3
func (e *UCCEVC) WriteIncentiveTableDescriptions(entity spineapi.EntityRemoteInterface, data []api.IncentiveTariffDescription) error {
	if !e.isCompatibleEntity(entity) {
		return api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := util.IncentiveTable(e.service, entity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return err
	}

	descriptions, err := evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
	if err != nil {
		logging.Log().Error(err)
		return err
	}

	// default tariff
	//
	// - tariff, min 1
	//   each tariff has
	//   - tiers: min 1, max 3
	//     each tier has:
	//     - boundaries: min 1, used for different power limits, e.g. 0-1kW x€, 1-3kW y€, ...
	//     - incentives: min 1, max 3
	//       - price/costs (absolute or relative)
	//       - renewable energy percentage
	//       - CO2 emissions
	//
	// limit this to
	// - 1 tariff
	//   - 1 tier
	//     - 1 boundary
	//     - 1 incentive (price)
	//       incentive type has to be the same for all sent power limits!
	descData := []model.IncentiveTableDescriptionType{
		{
			TariffDescription: descriptions[0].TariffDescription,
			Tier: []model.IncentiveTableDescriptionTierType{
				{
					TierDescription: &model.TierDescriptionDataType{
						TierId:   eebusutil.Ptr(model.TierIdType(0)),
						TierType: eebusutil.Ptr(model.TierTypeTypeDynamicCost),
					},
					BoundaryDescription: []model.TierBoundaryDescriptionDataType{
						{
							BoundaryId:   eebusutil.Ptr(model.TierBoundaryIdType(0)),
							BoundaryType: eebusutil.Ptr(model.TierBoundaryTypeTypePowerBoundary),
							BoundaryUnit: eebusutil.Ptr(model.UnitOfMeasurementTypeW),
						},
					},
					IncentiveDescription: []model.IncentiveDescriptionDataType{
						{
							IncentiveId:   eebusutil.Ptr(model.IncentiveIdType(0)),
							IncentiveType: eebusutil.Ptr(model.IncentiveTypeTypeAbsoluteCost),
							Currency:      eebusutil.Ptr(model.CurrencyTypeEur),
						},
					},
				},
			},
		},
	}

	if len(data) > 0 && len(data[0].Tiers) > 0 {
		newDescData := []model.IncentiveTableDescriptionType{}
		allDataPresent := false

		for index, tariff := range data {
			tariffDesc := descriptions[0].TariffDescription
			if len(descriptions) > index {
				tariffDesc = descriptions[index].TariffDescription
			}

			newTariff := model.IncentiveTableDescriptionType{
				TariffDescription: tariffDesc,
			}

			tierData := []model.IncentiveTableDescriptionTierType{}
			for _, tier := range tariff.Tiers {
				newTier := model.IncentiveTableDescriptionTierType{}

				newTier.TierDescription = &model.TierDescriptionDataType{
					TierId:   eebusutil.Ptr(model.TierIdType(tier.Id)),
					TierType: eebusutil.Ptr(tier.Type),
				}

				boundaryDescription := []model.TierBoundaryDescriptionDataType{}
				for _, boundary := range tier.Boundaries {
					newBoundary := model.TierBoundaryDescriptionDataType{
						BoundaryId:   eebusutil.Ptr(model.TierBoundaryIdType(boundary.Id)),
						BoundaryType: eebusutil.Ptr(boundary.Type),
						BoundaryUnit: eebusutil.Ptr(boundary.Unit),
					}
					boundaryDescription = append(boundaryDescription, newBoundary)
				}
				newTier.BoundaryDescription = boundaryDescription

				incentiveDescription := []model.IncentiveDescriptionDataType{}
				for _, incentive := range tier.Incentives {
					newIncentive := model.IncentiveDescriptionDataType{
						IncentiveId:   eebusutil.Ptr(model.IncentiveIdType(incentive.Id)),
						IncentiveType: eebusutil.Ptr(incentive.Type),
					}
					if incentive.Currency != "" {
						newIncentive.Currency = eebusutil.Ptr(incentive.Currency)
					}
					incentiveDescription = append(incentiveDescription, newIncentive)
				}
				newTier.IncentiveDescription = incentiveDescription

				if len(newTier.BoundaryDescription) > 0 &&
					len(newTier.IncentiveDescription) > 0 {
					allDataPresent = true
				}
				tierData = append(tierData, newTier)
			}

			newTariff.Tier = tierData

			newDescData = append(newDescData, newTariff)
		}

		if allDataPresent {
			descData = newDescData
		}
	}

	_, err = evIncentiveTable.WriteDescriptions(descData)
	if err != nil {
		logging.Log().Error(err)
		return err
	}

	return nil
}

// send incentives to the EV
// if no data is provided, default incentives with the same price for 7 days will be sent
func (e *UCCEVC) WriteIncentives(entity spineapi.EntityRemoteInterface, data []api.DurationSlotValue) error {
	if !e.isCompatibleEntity(entity) {
		return api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := util.IncentiveTable(e.service, entity)
	if err != nil {
		return features.ErrDataNotAvailable
	}

	if len(data) == 0 {
		// send default incentives for the maximum timeframe
		// to fullfill spec, as there is no data provided
		logging.Log().Info("Fallback sending default incentives")
		data = []api.DurationSlotValue{
			{Duration: 7 * time.Hour * 24, Value: 0.30},
		}
	}

	constraints, err := e.IncentiveConstraints(entity)
	if err != nil {
		return err
	}

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	incentiveSlots := []model.IncentiveTableIncentiveSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		timeInterval := &model.TimeTableDataType{
			StartTime: &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeStart),
			},
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeInterval.EndTime = &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeEndTime),
			}
		}

		incentiveSlot := model.IncentiveTableIncentiveSlotType{
			TimeInterval: timeInterval,
			Tier: []model.IncentiveTableTierType{
				{
					Tier: &model.TierDataType{
						TierId: eebusutil.Ptr(model.TierIdType(0)),
					},
					Boundary: []model.TierBoundaryDataType{
						{
							BoundaryId:         eebusutil.Ptr(model.TierBoundaryIdType(0)), // only 1 boundary exists
							LowerBoundaryValue: model.NewScaledNumberType(0),
						},
					},
					Incentive: []model.IncentiveDataType{
						{
							IncentiveId: eebusutil.Ptr(model.IncentiveIdType(0)), // always use price
							Value:       model.NewScaledNumberType(slot.Value),
						},
					},
				},
			},
		}
		incentiveSlots = append(incentiveSlots, incentiveSlot)

		totalDuration += slot.Duration
	}

	incentiveData := model.IncentiveTableType{
		Tariff: &model.TariffDataType{
			TariffId: eebusutil.Ptr(model.TariffIdType(0)),
		},
		IncentiveSlot: incentiveSlots,
	}

	_, err = evIncentiveTable.WriteValues([]model.IncentiveTableType{incentiveData})

	return err
}
