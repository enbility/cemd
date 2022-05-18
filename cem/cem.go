package cem

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type Cem struct {
	brand        string
	model        string
	serialNumber string
	identifier   string
	myService    *service.EEBUSService
}

func NewCEM(brand, model, serialNumber, identifier string) *Cem {
	return &Cem{
		brand:        brand,
		model:        model,
		serialNumber: serialNumber,
		identifier:   identifier,
	}
}

func (h *Cem) Setup(port, remoteSKI, certFile, keyFile string) error {
	serviceDescription := &service.ServiceDescription{
		Brand:        h.brand,
		Model:        h.model,
		SerialNumber: h.serialNumber,
		Identifier:   h.identifier,
		DeviceType:   model.DeviceTypeTypeEnergyManagementSystem,
		Interfaces:   []string{"en0"},
	}

	h.myService = service.NewEEBUSService(serviceDescription, h)

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

	// Setup the supported UseCases and their features
	evseSupport := AddEVSESupport(h.myService)
	evseSupport.Delegate = h
	evSupport, _ := AddEVSupport(h.myService)
	evSupport.Delegate = h

	h.myService.Start()
	// defer h.myService.Shutdown()

	remoteService := service.ServiceDetails{
		SKI: remoteSKI,
	}
	_ = h.myService.RegisterRemoteService(remoteService)

	return nil
}

// EEBUSServiceDelegate

// handle a request to trust a remote service
func (h *Cem) RemoteServiceTrustRequested(ski string) {
	// we directly trust it in this example
	h.myService.UpdateRemoteServiceTrust(ski, true)
}

// report the Ship ID of a newly trusted connection
func (h *Cem) RemoteServiceShipIDReported(ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	fmt.Println("SKI", ski, "has Ship ID:", shipID)
}

// EVSEDelegate

// handle device state updates from the remote EVSE device
func (h *Cem) HandleEVSEDeviceState(ski string, failure bool, errorCode string) {
	fmt.Println("EVSE Error State:", failure, errorCode)
}

// EVDelegate

// handle device state updates from the remote EVSE device
func (h *Cem) HandleEVEntityState(ski string, failure bool, errorCode string) {
	fmt.Println("EV Error State:", failure, errorCode)
}
