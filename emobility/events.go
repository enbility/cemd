package emobility

import (
	"time"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Internal EventHandler Interface for the CEM
func (e *EMobility) HandleEvent(payload api.EventPayload) {
	// only care about the registered SKI
	if payload.Ski != e.ski {
		return
	}

	// only care about events for this remote device
	if payload.Device != nil && payload.Device.Ski() != e.ski {
		return
	}

	// we care only about events from an EVSE or EV entity or device changes for this remote device
	var entityType model.EntityTypeType
	if payload.Entity != nil {
		entityType = payload.Entity.EntityType()
		if entityType != model.EntityTypeTypeEVSE && entityType != model.EntityTypeTypeEV {
			return
		}
	}

	switch payload.EventType {
	case api.EventTypeDeviceChange:
		if payload.ChangeType == api.ElementChangeRemove {
			e.evseDisconnected()
			e.evDisconnected()
		}

	case api.EventTypeEntityChange:
		if payload.Entity == nil {
			return
		}

		switch payload.ChangeType {
		case api.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseConnected(payload.Ski, payload.Entity)
			case model.EntityTypeTypeEV:
				e.evConnected(payload.Entity)
			}
		case api.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseDisconnected()
			case model.EntityTypeTypeEV:
				e.evDisconnected()
			}
		}

	case api.EventTypeDataChange:
		if payload.ChangeType == api.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceConfigurationKeyValueDescriptionListDataType:
				evDeviceConfiguration, err := e.deviceConfiguration(payload.Entity)
				if err != nil {
					break
				}

				// key value descriptions received, now get the data
				if _, err := evDeviceConfiguration.RequestKeyValues(); err != nil {
					logging.Log().Error("Error getting configuration key values:", err)
				}

			case *model.ElectricalConnectionParameterDescriptionListDataType:
				evElectricalConnection, err := e.electricalConnection(payload.Entity)
				if err != nil {
					break
				}

				if _, err := evElectricalConnection.RequestPermittedValueSets(); err != nil {
					logging.Log().Error("Error getting electrical permitted values:", err)
				}

			case *model.LoadControlLimitDescriptionListDataType:
				evLoadControl, err := e.loadControl(payload.Entity)
				if err != nil {
					break
				}

				if _, err := evLoadControl.RequestLimitValues(); err != nil {
					logging.Log().Error("Error getting loadcontrol limit values:", err)
				}

			case *model.MeasurementDescriptionListDataType:
				evMeasurement, err := e.measurement(payload.Entity)
				if err != nil {
					break
				}

				if _, err := evMeasurement.RequestValues(); err != nil {
					logging.Log().Error("Error getting measurement list values:", err)
				}

			case *model.TimeSeriesDescriptionListDataType:
				evTimeSeries, err := e.timeSeries(payload.Entity)
				if err != nil || payload.CmdClassifier == nil {
					break
				}

				switch *payload.CmdClassifier {
				case model.CmdClassifierTypeReply:
					if err := evTimeSeries.RequestConstraints(); err == nil {
						break
					}

					// if constraints do not exist, directly request values
					e.evRequestTimeSeriesValues(payload.Entity)

				case model.CmdClassifierTypeNotify:
					// check if we are required to update the plan
					if !e.evCheckTimeSeriesDescriptionConstraintsUpdateRequired(payload.Entity) {
						break
					}

					demand, err := e.EVEnergyDemand(payload.Entity)
					if err != nil {
						logging.Log().Error("Error getting energy demand:", err)
						break
					}

					if e.dataProvider != nil {
						e.dataProvider.EVProvidedEnergyDemand(demand)
					}

					timeConstraints, err := e.EVTimeSlotConstraints(payload.Entity)
					if err != nil {
						logging.Log().Error("Error getting timeseries constraints:", err)
						break
					}

					incentiveConstraints, err := e.EVIncentiveConstraints(payload.Entity)
					if err != nil {
						logging.Log().Error("Error getting incentive constraints:", err)
						break
					}

					if e.dataProvider != nil {
						e.dataProvider.EVRequestPowerLimits(demand, timeConstraints)
						e.dataProvider.EVRequestIncentives(demand, incentiveConstraints)
						break
					}

					e.evWriteDefaultIncentives(payload.Entity)
					e.evWriteDefaultPowerLimits(payload.Entity)
				}

			case *model.TimeSeriesConstraintsListDataType:
				if _, err := e.timeSeries(payload.Entity); err != nil || payload.CmdClassifier == nil {
					break
				}

				if *payload.CmdClassifier != model.CmdClassifierTypeReply {
					break
				}

				e.evRequestTimeSeriesValues(payload.Entity)

			case *model.TimeSeriesListDataType:
				if _, err := e.timeSeries(payload.Entity); err != nil || payload.CmdClassifier == nil {
					break
				}

				// check if we received a plan
				e.evForwardChargePlanIfProvided(payload.Entity)

			case *model.IncentiveTableDescriptionDataType:
				evIncentiveTable, err := e.incentiveTable(payload.Entity)
				if err != nil || payload.CmdClassifier == nil {
					break
				}

				switch *payload.CmdClassifier {
				case model.CmdClassifierTypeReply:
					if err := evIncentiveTable.RequestConstraints(); err == nil {
						break
					}

					// if constraints do not exist, directly request values
					e.evRequestIncentiveValues(payload.Entity)

				case model.CmdClassifierTypeNotify:
					// check if we are required to update the plan
					if !e.evCheckIncentiveTableDescriptionUpdateRequired(payload.Entity) {
						break
					}

					e.evWriteIncentiveTableDescriptions(payload.Entity)
				}

			case *model.IncentiveTableConstraintsDataType:
				if *payload.CmdClassifier == model.CmdClassifierTypeReply {
					e.evRequestIncentiveValues(payload.Entity)
				}
			}
		}
	}

	if e.dataProvider == nil || payload.Entity == nil {
		return
	}

	// check if the charge strategy changed
	chargeStrategy := e.EVChargeStrategy(payload.Entity)
	if chargeStrategy == e.evCurrentChargeStrategy {
		return
	}

	// update the current value and inform the dataProvider
	e.evCurrentChargeStrategy = chargeStrategy
	e.dataProvider.EVProvidedChargeStrategy(chargeStrategy)
}

func (e *EMobility) localCemEntity() api.EntityLocalInterface {
	localDevice := e.service.LocalDevice()

	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	return localEntity
}

func (e *EMobility) evWriteDefaultIncentives(remoteEntity api.EntityRemoteInterface) {
	// send default incentives for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default incentives")
	data := []EVDurationSlotValue{
		{Duration: 7 * time.Hour * 24, Value: 0.30},
	}
	_ = e.EVWriteIncentives(remoteEntity, data)
}

func (e *EMobility) evWriteDefaultPowerLimits(remoteEntity api.EntityRemoteInterface) {
	// send default power limits for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default power limits")

	evElectricalConnection, err := e.electricalConnection(remoteEntity)
	if err != nil {
		logging.Log().Error("electrical connection feature not found")
		return
	}

	paramDesc, err := evElectricalConnection.GetParameterDescriptionForScopeType(model.ScopeTypeTypeACPower)
	if err != nil {
		logging.Log().Error("Error getting parameter descriptions:", err)
		return
	}

	permitted, err := evElectricalConnection.GetPermittedValueSetForParameterId(*paramDesc.ParameterId)
	if err != nil {
		logging.Log().Error("Error getting permitted values:", err)
		return
	}

	if len(permitted.PermittedValueSet) < 1 || len(permitted.PermittedValueSet[0].Range) < 1 {
		logging.Log().Error("No permitted value set available")
		return
	}

	data := []EVDurationSlotValue{
		{Duration: 7 * time.Hour * 24, Value: permitted.PermittedValueSet[0].Range[0].Max.GetValue()},
	}
	_ = e.EVWritePowerLimits(remoteEntity, data)
}

// request time series values
func (e *EMobility) evRequestTimeSeriesValues(remoteEntity api.EntityRemoteInterface) {
	localEntity := e.localCemEntity()

	evTimeSeries, err := features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
	if err != nil {
		return
	}

	if _, err := evTimeSeries.RequestValues(); err != nil {
		logging.Log().Error("Error getting time series list values:", err)
	}
}

// send the ev provided charge plan to the CEM
func (e *EMobility) evForwardChargePlanIfProvided(remoteEntity api.EntityRemoteInterface) {
	if e.dataProvider == nil {
		return
	}

	if plan, err := e.EVChargePlan(remoteEntity); err == nil {
		e.dataProvider.EVProvidedChargePlan(plan)
	}

	if constraints, err := e.EVChargePlanConstraints(remoteEntity); err == nil {
		e.dataProvider.EVProvidedChargePlanConstraints(constraints)
	}
}

// request incentive table values
func (e *EMobility) evRequestIncentiveValues(remoteEntity api.EntityRemoteInterface) {
	localEntity := e.localCemEntity()

	evIncentiveTable, err := features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, remoteEntity)
	if err != nil {
		return
	}

	if _, err := evIncentiveTable.RequestValues(); err != nil {
		logging.Log().Error("Error getting time series list values:", err)
	}
}

// process required steps when an evse is connected
func (e *EMobility) evseConnected(ski string, entity api.EntityRemoteInterface) {
	e.evseEntity = entity
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	if evseDeviceClassification, err := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceClassification.RequestManufacturerDetails()
	}

	if evseDeviceDiagnosis, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity); err == nil {
		_, _ = evseDeviceDiagnosis.RequestState()
	}
}

// an EV was disconnected
func (e *EMobility) evseDisconnected() {
	e.evseEntity = nil

	e.evDisconnected()
}

// an EV was disconnected, trigger required cleanup
func (e *EMobility) evDisconnected() {
	if e.evEntity == nil {
		return
	}

	e.evEntity = nil

	logging.Log().Debug("ev disconnected")

	// TODO: add error handling
}

// an EV was connected, trigger required communication
func (e *EMobility) evConnected(entity api.EntityRemoteInterface) {
	e.evEntity = entity

	logging.Log().Debug("ev connected")

	// initialise features, e.g. subscriptions, bindings
	if evDeviceClassification, err := e.deviceClassification(entity); err == nil {
		if err := evDeviceClassification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get manufacturer details
		if _, err := evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceConfiguration, err := e.deviceConfiguration(entity); err == nil {
		if err := evDeviceConfiguration.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
		// get ev configuration data
		if err := evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceDiagnosis, err := e.deviceDiagnosis(entity); err == nil {
		if err := evDeviceDiagnosis.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get device diagnosis state
		if _, err := evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := e.electricalConnection(entity); err == nil {
		if err := evElectricalConnection.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical connection parameter
		if err := evElectricalConnection.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

		if err := evElectricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

	}

	if evMeasurement, err := e.measurement(entity); err == nil {
		if err := evMeasurement.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get measurement parameters
		if err := evMeasurement.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

	}

	if evLoadControl, err := e.loadControl(entity); err == nil {
		if err := evLoadControl.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if err := evLoadControl.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get loadlimit parameter
		if err := evLoadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}

	}

	if evIdentification, err := e.identification(entity); err == nil {
		if err := evIdentification.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		// get identification
		if _, err := evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if e.configuration.CoordinatedChargingEnabled {
		if evTimeSeries, err := e.timeSeries(entity); err == nil {
			if err := evTimeSeries.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}

			if err := evTimeSeries.Bind(); err != nil {
				logging.Log().Debug(err)
			}

			// get time series parameter
			if err := evTimeSeries.RequestDescriptions(); err != nil {
				logging.Log().Debug(err)
			}

		}

		if evIncentiveTable, err := e.incentiveTable(entity); err == nil {
			if err := evIncentiveTable.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}

			if err := evIncentiveTable.Bind(); err != nil {
				logging.Log().Debug(err)
			}

			// get incentive table parameter
			if err := evIncentiveTable.RequestDescriptions(); err != nil {
				logging.Log().Debug(err)
			}

		}

	}
}

// inform the EVSE about used currency and boundary units
//
// # SPINE UC CoordinatedEVCharging 2.4.3
func (e *EMobility) evWriteIncentiveTableDescriptions(remoteEntity api.EntityRemoteInterface) {
	evIncentiveTable, err := e.incentiveTable(remoteEntity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return
	}

	descriptions, err := evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
	if err != nil {
		logging.Log().Error(err)
		return
	}

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
	data := []model.IncentiveTableDescriptionType{
		{
			TariffDescription: descriptions[0].TariffDescription,
			Tier: []model.IncentiveTableDescriptionTierType{
				{
					TierDescription: &model.TierDescriptionDataType{
						TierId:   util.Ptr(model.TierIdType(1)),
						TierType: util.Ptr(model.TierTypeTypeDynamicCost),
					},
					BoundaryDescription: []model.TierBoundaryDescriptionDataType{
						{
							BoundaryId:   util.Ptr(model.TierBoundaryIdType(1)),
							BoundaryType: util.Ptr(model.TierBoundaryTypeTypePowerBoundary),
							BoundaryUnit: util.Ptr(model.UnitOfMeasurementTypeW),
						},
					},
					IncentiveDescription: []model.IncentiveDescriptionDataType{
						{
							IncentiveId:   util.Ptr(model.IncentiveIdType(1)),
							IncentiveType: util.Ptr(model.IncentiveTypeTypeAbsoluteCost),
							Currency:      util.Ptr(e.currency),
						},
					},
				},
			},
		},
	}

	_, err = evIncentiveTable.WriteDescriptions(data)
	if err != nil {
		logging.Log().Error(err)
	}
}

// check timeSeries descriptions if constraints element has updateRequired set to true
// as this triggers the CEM to send power tables within 20s
func (e *EMobility) evCheckTimeSeriesDescriptionConstraintsUpdateRequired(remoteEntity api.EntityRemoteInterface) bool {
	evTimeSeries, err := e.timeSeries(remoteEntity)
	if err != nil {
		logging.Log().Error("timeseries feature not found")
		return false
	}

	data, err := evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints)
	if err != nil {
		return false
	}

	if data.UpdateRequired != nil {
		return *data.UpdateRequired
	}

	return false
}

// check incentibeTable descriptions if the tariff description has updateRequired set to true
// as this triggers the CEM to send incentive tables within 20s
func (e *EMobility) evCheckIncentiveTableDescriptionUpdateRequired(remoteEntity api.EntityRemoteInterface) bool {
	evIncentiveTable, err := e.incentiveTable(remoteEntity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return false
	}

	data, err := evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
	if err != nil {
		return false
	}

	// only use the first description and therein the first tariff
	item := data[0].TariffDescription
	if item.UpdateRequired != nil {
		return *item.UpdateRequired
	}

	return false
}
