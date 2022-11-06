package cem

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/DerAndereAndi/eebus-go-cem/emobility"
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type CemImpl struct {
	brand        string
	model        string
	serialNumber string
	identifier   string
	myService    *service.EEBUSService
	emobility    *emobility.EMobilityImpl
}

func NewCEM(brand, model, serialNumber, identifier string) *CemImpl {
	return &CemImpl{
		brand:        brand,
		model:        model,
		serialNumber: serialNumber,
		identifier:   identifier,
	}
}

func (h *CemImpl) Setup(port, remoteSKI, certFile, keyFile string, ifaces []string) error {
	serviceDescription := &service.ServiceDescription{
		Brand:        h.brand,
		Model:        h.model,
		SerialNumber: h.serialNumber,
		Identifier:   h.identifier,
		DeviceType:   model.DeviceTypeTypeEnergyManagementSystem,
		Interfaces:   ifaces,
	}

	h.myService = service.NewEEBUSService(serviceDescription, h)
	h.myService.SetLogging(h)

	var err error
	var certificate tls.Certificate

	serviceDescription.Port, err = strconv.Atoi(port)
	if err != nil {
		return err
	}

	certificate, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	serviceDescription.Certificate = certificate

	if err = h.myService.Setup(); err != nil {
		return err
	}

	spine.Events.Subscribe(h)

	// Setup the supported usecases and features
	emobilityScenario := emobility.NewEMobilityScenario(h.myService)
	emobilityScenario.AddFeatures()
	emobilityScenario.AddUseCases()

	// TODO: emobility should be stored per remote SKI and
	// only be set for the SKI if the device supports it
	h.emobility = emobility.NewEMobility(h.myService, remoteSKI)

	h.myService.Start()
	// defer h.myService.Shutdown()

	remoteService := service.ServiceDetails{
		SKI: remoteSKI,
	}
	h.myService.RegisterRemoteService(remoteService)

	return nil
}

// EEBUSServiceDelegate

// handle a request to trust a remote service
func (h *CemImpl) RemoteServiceTrustRequested(ski string) {
	// we directly trust it in this example
	h.myService.UpdateRemoteServiceTrust(ski, true)
}

// report the Ship ID of a newly trusted connection
func (h *CemImpl) RemoteServiceShipIDReported(ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	logging.Log.Info("SKI", ski, "has Ship ID:", shipID)
}

func (h *CemImpl) RemoteSKIConnected(ski string) {}

func (h *CemImpl) RemoteSKIDisconnected(ski string) {}

// Logging interface

func (h *CemImpl) log(level string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s", t.Format(time.RFC3339), level, fmt.Sprintln(args...))
}

func (h *CemImpl) logf(level, format string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s\n", t.Format(time.RFC3339), level, fmt.Sprintf(format, args...))
}

func (h *CemImpl) Trace(args ...interface{}) {
	h.log("TRACE", args...)
}

func (h *CemImpl) Tracef(format string, args ...interface{}) {
	h.logf("TRACE", format, args...)
}

func (h *CemImpl) Debug(args ...interface{}) {
	h.log("DEBUG", args...)
}

func (h *CemImpl) Debugf(format string, args ...interface{}) {
	h.logf("DEBUG", format, args...)
}

func (h *CemImpl) Info(args ...interface{}) {
	h.log("INFO", args...)
}

func (h *CemImpl) Infof(format string, args ...interface{}) {
	h.logf("INFO", format, args...)
}

func (h *CemImpl) Error(args ...interface{}) {
	h.log("ERROR", args...)
}

func (h *CemImpl) Errorf(format string, args ...interface{}) {
	h.logf("ERROR", format, args...)
}
