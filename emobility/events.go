package emobility

import (
	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Internal EventHandler Interface for the CEM
func (e *EMobilityImpl) HandleEvent(payload spine.EventPayload) {
	// we only care about the registered SKI
	if payload.Ski != e.ski {
		return
	}

	// we care only about events for this remote device
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
	case spine.EventTypeDeviceChange:
		switch payload.ChangeType {
		case spine.ElementChangeRemove:
			e.evseDisconnected()
			e.evDisconnected()
		}

	case spine.EventTypeEntityChange:
		if payload.Entity == nil {
			return
		}

		switch payload.ChangeType {
		case spine.ElementChangeAdd:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseConnected(payload.Ski, payload.Entity)
			case model.EntityTypeTypeEV:
				e.evConnected(payload.Entity)
			}
		case spine.ElementChangeRemove:
			switch entityType {
			case model.EntityTypeTypeEVSE:
				e.evseDisconnected()
			case model.EntityTypeTypeEV:
				e.evDisconnected()
			}
		}

	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			/*
				case *model.DeviceClassificationManufacturerDataType:
					var feature *features.DeviceClassification
					if entityType == model.EntityTypeTypeEVSE {
						feature = e.evseDeviceClassification
					} else {
						feature = e.evDeviceClassification
					}
					_, err := feature.GetManufacturerDetails()
					if err != nil {
						logging.Log.Error("Error getting manufacturer data:", err)
						return
					}

					// TODO: provide the current data to the CEM
			*/
			case *model.DeviceConfigurationKeyValueDescriptionListDataType:
				// key value descriptions received, now get the data
				_, err := e.evDeviceConfiguration.RequestKeyValueList()
				if err != nil {
					logging.Log.Error("Error getting configuration key values:", err)
				}

				/*
					case *model.DeviceConfigurationKeyValueListDataType:
						data, err := e.evDeviceConfiguration.GetValues()
						if err != nil {
							logging.Log.Error("Error getting device configuration values:", err)
							return
						}

						// TODO: provide the device configuration data
						logging.Log.Debugf("Device Configuration Values: %#v\n", data)
				*/

				/*
					case *model.DeviceDiagnosisStateDataType:
						var feature *features.DeviceDiagnosis
						if entityType == model.EntityTypeTypeEVSE {
							feature = e.evseDeviceDiagnosis
						} else {
							feature = e.evDeviceDiagnosis
						}
						_, err := feature.GetState()
						if err != nil {
							logging.Log.Error("Error getting device diagnosis state:", err)
						}
				*/

				/*
					case *model.ElectricalConnectionDescriptionListDataType:
						data, err := e.evElectricalConnection.GetDescription()
						if err != nil {
							logging.Log.Error("Error getting electrical description:", err)
							return
						}

						// TODO: provide the electrical description data
						logging.Log.Debugf("Electrical Description: %#v\n", data)
				*/

			case *model.ElectricalConnectionParameterDescriptionListDataType:
				_, err := e.evElectricalConnection.RequestPermittedValueSet()
				if err != nil {
					logging.Log.Debug("Error getting electrical permitted values:", err)
				}

				/*
					case *model.ElectricalConnectionPermittedValueSetListDataType:
						data, err := e.evElectricalConnection.GetEVLimitValues()
						if err != nil {
							logging.Log.Error("Error getting electrical limit values:", err)
							return
						}

						// TODO: provide the electrical limit data
						logging.Log.Debugf("Electrical Permitted Values: %#v\n", data)
				*/

			case *model.LoadControlLimitDescriptionListDataType:
				_, err := e.evLoadControl.RequestLimits()
				if err != nil {
					logging.Log.Debug("Error getting loadcontrol limit values:", err)
				}

				/*
					case *model.LoadControlLimitListDataType:
						data, err := e.evLoadControl.GetLimitValues()
						if err != nil {
							logging.Log.Error("Error getting loadcontrol limit values:", err)
							return
						}

						// TODO: provide the loadcontrol limit data
						logging.Log.Debugf("Loadcontrol Limits: %#v\n", data)
				*/

			case *model.MeasurementDescriptionListDataType:
				_, err := e.evMeasurement.Request()
				if err != nil {
					logging.Log.Debug("Error getting measurement list values:", err)
				}

				/*
					case *model.MeasurementListDataType:
						data, err := e.evMeasurement.GetValues()
						if err != nil {
							logging.Log.Error("Error getting measurement values:", err)
							return
						}

						// TODO: provide the measurement data
						logging.Log.Debugf("Measurements: %#v\n", data)
				*/

				/*
					case *model.IdentificationListDataType:
						data, err := e.evIdentification.GetValues()
						if err != nil {
							logging.Log.Error("Error getting identification values:", err)
							return
						}

						// TODO: provide the device configuration data
						logging.Log.Debugf("Identification Values: %#v\n", data)
				*/
			}
		}
	}
}

// process required steps when an evse is connected
func (e *EMobilityImpl) evseConnected(ski string, entity *spine.EntityRemoteImpl) {
	e.evseEntity = entity
	localDevice := e.service.LocalDevice()

	f1, err := features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	if err != nil {
		return
	}
	e.evseDeviceClassification = f1

	f2, err := features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	if err != nil {
		return
	}
	e.evseDeviceDiagnosis = f2

	_, _ = e.evseDeviceClassification.RequestManufacturerDetailsForEntity()
	_, _ = e.evseDeviceDiagnosis.RequestStateForEntity()
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

	logging.Log.Debug("ev disconnected")

	// TODO: add error handling

}

// an EV was connected, trigger required communication
func (e *EMobilityImpl) evConnected(entity *spine.EntityRemoteImpl) {
	e.evEntity = entity
	localDevice := e.service.LocalDevice()

	logging.Log.Debug("ev connected")

	// TODO: add error handling

	// setup features
	e.evDeviceClassification, _ = features.NewDeviceClassification(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evDeviceDiagnosis, _ = features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evDeviceConfiguration, _ = features.NewDeviceConfiguration(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evElectricalConnection, _ = features.NewElectricalConnection(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evMeasurement, _ = features.NewMeasurement(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evIdentification, _ = features.NewIdentification(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)
	e.evLoadControl, _ = features.NewLoadControl(model.RoleTypeClient, model.RoleTypeServer, localDevice, entity)

	// subscribe
	if err := e.evDeviceClassification.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evDeviceConfiguration.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evDeviceDiagnosis.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evElectricalConnection.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evMeasurement.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evIdentification.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	if err := e.evLoadControl.SubscribeForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}
	// if err := util.SubscribeTimeSeriesForEntity(e.service, entity); err != nil {
	// 	logging.Log.Error(err)
	// 	return
	// }
	// if err := util.SubscribeIncentiveTableForEntity(e.service, entity); err != nil {
	// 	logging.Log.Error(err)
	// 	return
	// }

	// bindings
	if err := e.evLoadControl.Bind(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get ev configuration data
	if err := e.evDeviceConfiguration.Request(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get manufacturer details
	if _, err := e.evDeviceClassification.RequestManufacturerDetailsForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get device diagnosis state
	if _, err := e.evDeviceDiagnosis.RequestStateForEntity(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get electrical connection parameter
	if err := e.evElectricalConnection.RequestDescription(); err != nil {
		logging.Log.Debug(err)
		return
	}

	if err := e.evElectricalConnection.RequestParameterDescription(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get measurement parameters
	if err := e.evMeasurement.RequestDescription(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get identification
	if _, err := e.evIdentification.Request(); err != nil {
		logging.Log.Debug(err)
		return
	}

	// get loadlimit parameter
	if err := e.evLoadControl.RequestLimitDescription(); err != nil {
		logging.Log.Debug(err)
		return
	}

}
