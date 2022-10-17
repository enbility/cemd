package cem

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/DerAndereAndi/eebus-go-cem/usecases"
	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
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

func (h *Cem) Setup(port, remoteSKI, certFile, keyFile string, ifaces []string) error {
	serviceDescription := &service.ServiceDescription{
		Brand:        h.brand,
		Model:        h.model,
		SerialNumber: h.serialNumber,
		Identifier:   h.identifier,
		DeviceType:   model.DeviceTypeTypeEnergyManagementSystem,
		Interfaces:   ifaces,
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

	spine.Events.Subscribe(h)

	// Setup the supported UseCases and their features
	localEntity := h.myService.LocalEntity()

	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient, "Device Configuration Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient, "Device Classification Client")
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient, "Device Diagnosis Client")
	}
	{
		f := localEntity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer, "Device Diagnosis Server")
		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)

		// Set the initial state
		state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
		deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
			OperatingState: &state,
		}
		f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

		f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, false)
	}
	{
		_ = localEntity.GetOrAddFeature(model.FeatureTypeTypeIdentification, model.RoleTypeClient, "Identification Client")
	}

	// e-mobilty specific use cases
	_ = usecases.NewEVSECommissioningAndConfiguration(h.myService)
	_ = usecases.NewEVCommissioningAndConfiguration(h.myService)

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
