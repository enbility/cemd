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
func (e *EMobilityImpl) HandleEvent(payload api.EventPayload) {
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
		switch payload.ChangeType {
		case api.ElementChangeRemove:
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
				if e.evDeviceConfiguration == nil {
					break
				}

				// key value descriptions received, now get the data
				if _, err := e.evDeviceConfiguration.RequestKeyValues(); err != nil {
					logging.Log().Error("Error getting configuration key values:", err)
				}

			case *model.ElectricalConnectionParameterDescriptionListDataType:
				if e.evElectricalConnection == nil {
					break
				}
				if _, err := e.evElectricalConnection.RequestPermittedValueSets(); err != nil {
					logging.Log().Error("Error getting electrical permitted values:", err)
				}

			case *model.LoadControlLimitDescriptionListDataType:
				if e.evLoadControl == nil {
					break
				}
				if _, err := e.evLoadControl.RequestLimitValues(); err != nil {
					logging.Log().Error("Error getting loadcontrol limit values:", err)
				}

			case *model.MeasurementDescriptionListDataType:
				if e.evMeasurement == nil {
					break
				}
				if _, err := e.evMeasurement.RequestValues(); err != nil {
					logging.Log().Error("Error getting measurement list values:", err)
				}

			case *model.TimeSeriesDescriptionListDataType:
				if e.evTimeSeries == nil || payload.CmdClassifier == nil {
					break
				}

				switch *payload.CmdClassifier {
				case model.CmdClassifierTypeReply:
					if err := e.evTimeSeries.RequestConstraints(); err == nil {
						break
					}

					// if constraints do not exist, directly request values
					e.evRequestTimeSeriesValues()

				case model.CmdClassifierTypeNotify:
					// check if we are required to update the plan
					if !e.evCheckTimeSeriesDescriptionConstraintsUpdateRequired() {
						break
					}

					demand, err := e.EVEnergyDemand()
					if err != nil {
						logging.Log().Error("Error getting energy demand:", err)
						break
					}

					if e.dataProvider != nil {
						e.dataProvider.EVProvidedEnergyDemand(demand)
					}

					timeConstraints, err := e.EVTimeSlotConstraints()
					if err != nil {
						logging.Log().Error("Error getting timeseries constraints:", err)
						break
					}

					incentiveConstraints, err := e.EVIncentiveConstraints()
					if err != nil {
						logging.Log().Error("Error getting incentive constraints:", err)
						break
					}

					if e.dataProvider != nil {
						e.dataProvider.EVRequestPowerLimits(demand, timeConstraints)
						e.dataProvider.EVRequestIncentives(demand, incentiveConstraints)
						break
					}

					e.evWriteDefaultIncentives()
					e.evWriteDefaultPowerLimits()
				}

			case *model.TimeSeriesConstraintsListDataType:
				if e.evTimeSeries == nil || payload.CmdClassifier == nil {
					break
				}

				if *payload.CmdClassifier != model.CmdClassifierTypeReply {
					break
				}

				e.evRequestTimeSeriesValues()

			case *model.TimeSeriesListDataType:
				if e.evTimeSeries == nil || payload.CmdClassifier == nil {
					break
				}

				// check if we received a plan
				e.evForwardChargePlanIfProvided()

			case *model.IncentiveTableDescriptionDataType:
				if e.evIncentiveTable == nil || payload.CmdClassifier == nil {
					break
				}

				switch *payload.CmdClassifier {
				case model.CmdClassifierTypeReply:
					if err := e.evIncentiveTable.RequestConstraints(); err == nil {
						break
					}

					// if constraints do not exist, directly request values
					e.evRequestIncentiveValues()

				case model.CmdClassifierTypeNotify:
					// check if we are required to update the plan
					if !e.evCheckIncentiveTableDescriptionUpdateRequired() {
						break
					}

					e.evWriteIncentiveTableDescriptions()
				}

			case *model.IncentiveTableConstraintsDataType:
				if *payload.CmdClassifier == model.CmdClassifierTypeReply {
					e.evRequestIncentiveValues()
				}
			}
		}
	}

	if e.dataProvider == nil {
		return
	}

	// check if the charge strategy changed
	chargeStrategy := e.EVChargeStrategy()
	if chargeStrategy == e.evCurrentChargeStrategy {
		return
	}

	// update the current value and inform the dataProvider
	e.evCurrentChargeStrategy = chargeStrategy
	e.dataProvider.EVProvidedChargeStrategy(chargeStrategy)
}

func (e *EMobilityImpl) evWriteDefaultIncentives() {
	// send default incentives for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default incentives")
	data := []EVDurationSlotValue{
		{Duration: 7 * time.Hour * 24, Value: 0.30},
	}
	_ = e.EVWriteIncentives(data)
}

func (e *EMobilityImpl) evWriteDefaultPowerLimits() {
	// send default power limits for the maximum timeframe
	// to fullfill spec, as there is no data provided
	logging.Log().Info("Fallback sending default power limits")

	paramDesc, err := e.evElectricalConnection.GetParameterDescriptionForScopeType(model.ScopeTypeTypeACPower)
	if err != nil {
		logging.Log().Error("Error getting parameter descriptions:", err)
		return
	}

	permitted, err := e.evElectricalConnection.GetPermittedValueSetForParameterId(*paramDesc.ParameterId)
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
	_ = e.EVWritePowerLimits(data)
}

// request time series values
func (e *EMobilityImpl) evRequestTimeSeriesValues() {
	if e.evTimeSeries == nil {
		return
	}

	if _, err := e.evTimeSeries.RequestValues(); err != nil {
		logging.Log().Error("Error getting time series list values:", err)
	}
}

// send the ev provided charge plan to the CEM
func (e *EMobilityImpl) evForwardChargePlanIfProvided() {
	if e.dataProvider == nil {
		return
	}

	if plan, err := e.EVChargePlan(); err == nil {
		e.dataProvider.EVProvidedChargePlan(plan)
	}

	if constraints, err := e.EVChargePlanConstraints(); err == nil {
		e.dataProvider.EVProvidedChargePlanConstraints(constraints)
	}
}

// request incentive table values
func (e *EMobilityImpl) evRequestIncentiveValues() {
	if e.evIncentiveTable == nil {
		return
	}

	if _, err := e.evIncentiveTable.RequestValues(); err != nil {
		logging.Log().Error("Error getting time series list values:", err)
	}
}

// process required steps when an evse is connected
func (e *EMobilityImpl) evseConnected(ski string, entity api.EntityRemote) {
	e.evseEntity = entity
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f1, err := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.evseDeviceClassification = f1

	f2, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if err != nil {
		return
	}
	e.evseDeviceDiagnosis = f2

	_, _ = e.evseDeviceClassification.RequestManufacturerDetails()
	_, _ = e.evseDeviceDiagnosis.RequestState()
}

// an EV was disconnected
func (e *EMobilityImpl) evseDisconnected() {
	e.evseEntity = nil

	e.evseDeviceClassification = nil
	e.evseDeviceDiagnosis = nil

	e.evDisconnected()
}

// an EV was disconnected, trigger required cleanup
func (e *EMobilityImpl) evDisconnected() {
	if e.evEntity == nil {
		return
	}

	e.evEntity = nil

	e.evDeviceClassification = nil
	e.evDeviceDiagnosis = nil
	e.evDeviceConfiguration = nil
	e.evElectricalConnection = nil
	e.evMeasurement = nil
	e.evIdentification = nil
	e.evLoadControl = nil
	e.evTimeSeries = nil
	e.evIncentiveTable = nil

	logging.Log().Debug("ev disconnected")

	// TODO: add error handling
}

// an EV was connected, trigger required communication
func (e *EMobilityImpl) evConnected(entity api.EntityRemote) {
	e.evEntity = entity
	localDevice := e.service.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	logging.Log().Debug("ev connected")

	// setup features
	e.evDeviceClassification, _ = features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evDeviceDiagnosis, _ = features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evDeviceConfiguration, _ = features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evElectricalConnection, _ = features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evMeasurement, _ = features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evIdentification, _ = features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	e.evLoadControl, _ = features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	if e.configuration.CoordinatedChargingEnabled {
		e.evTimeSeries, _ = features.NewTimeSeries(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
		e.evIncentiveTable, _ = features.NewIncentiveTable(model.RoleTypeClient, model.RoleTypeServer, localEntity, entity)
	}

	// optional requests are only logged as debug

	// subscribe
	if e.evDeviceClassification != nil {
		if err := e.evDeviceClassification.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evDeviceConfiguration != nil {
		if err := e.evDeviceConfiguration.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evDeviceDiagnosis != nil {
		if err := e.evDeviceDiagnosis.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evElectricalConnection != nil {
		if err := e.evElectricalConnection.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evMeasurement != nil {
		if err := e.evMeasurement.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evLoadControl != nil {
		if err := e.evLoadControl.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}
	if e.evIdentification != nil {
		if err := e.evIdentification.SubscribeForEntity(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if e.configuration.CoordinatedChargingEnabled {
		if e.evTimeSeries != nil {
			if err := e.evTimeSeries.SubscribeForEntity(); err != nil {
				logging.Log().Debug(err)
			}
		}
		// this is optional
		if e.evIncentiveTable != nil {
			if err := e.evIncentiveTable.SubscribeForEntity(); err != nil {
				logging.Log().Debug(err)
			}
		}
	}

	// bindings
	if err := e.evLoadControl.Bind(); err != nil {
		logging.Log().Debug(err)
	}

	if e.configuration.CoordinatedChargingEnabled {
		if e.evTimeSeries != nil {
			// this is optional
			if err := e.evTimeSeries.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if e.evIncentiveTable != nil {
			// this is optional
			if err := e.evIncentiveTable.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}
	}

	// get ev configuration data
	if e.evDeviceConfiguration != nil {
		if err := e.evDeviceConfiguration.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get manufacturer details
	if e.evDeviceClassification != nil {
		if _, err := e.evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get device diagnosis state
	if e.evDeviceDiagnosis != nil {
		if _, err := e.evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get electrical connection parameter
	if e.evElectricalConnection != nil {
		if err := e.evElectricalConnection.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if e.evElectricalConnection != nil {
		if err := e.evElectricalConnection.RequestParameterDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get measurement parameters
	if e.evMeasurement != nil {
		if err := e.evMeasurement.RequestDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get loadlimit parameter
	if e.evLoadControl != nil {
		if err := e.evLoadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	// get identification
	if e.evIdentification != nil {
		if _, err := e.evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if e.configuration.CoordinatedChargingEnabled {
		if e.evTimeSeries != nil {
			// get time series parameter
			if err := e.evTimeSeries.RequestDescriptions(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if e.evIncentiveTable != nil {
			// get incentive table parameter
			if err := e.evIncentiveTable.RequestDescriptions(); err != nil {
				logging.Log().Debug(err)
			}
		}
	}
}

// inform the EVSE about used currency and boundary units
//
// # SPINE UC CoordinatedEVCharging 2.4.3
func (e *EMobilityImpl) evWriteIncentiveTableDescriptions() {
	if e.evIncentiveTable == nil {
		return
	}

	descriptions, err := e.evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
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

	_, err = e.evIncentiveTable.WriteDescriptions(data)
	if err != nil {
		logging.Log().Error(err)
	}
}

// check timeSeries descriptions if constraints element has updateRequired set to true
// as this triggers the CEM to send power tables within 20s
func (e *EMobilityImpl) evCheckTimeSeriesDescriptionConstraintsUpdateRequired() bool {
	if e.evTimeSeries == nil {
		return false
	}

	data, err := e.evTimeSeries.GetDescriptionForType(model.TimeSeriesTypeTypeConstraints)
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
func (e *EMobilityImpl) evCheckIncentiveTableDescriptionUpdateRequired() bool {
	if e.evIncentiveTable == nil {
		return false
	}

	data, err := e.evIncentiveTable.GetDescriptionsForScope(model.ScopeTypeTypeSimpleIncentiveTable)
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
