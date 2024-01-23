package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/cem"
	"github.com/enbility/cemd/emobility"
	"github.com/enbility/cemd/grid"
	"github.com/enbility/cemd/inverterbatteryvis"
	"github.com/enbility/cemd/inverterpvvis"
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/model"
)

type DemoCem struct {
	cem *cem.CemImpl

	emobilityScenario, gridScenario, inverterBatteryVisScenario, inverterPVVisScenario api.SolutionInterface
}

func NewDemoCem(configuration *eebusapi.Configuration) *DemoCem {
	demo := &DemoCem{}

	demo.cem = cem.NewCEM(configuration, demo, demo)

	return demo
}

func (d *DemoCem) Setup() error {
	if err := d.cem.Setup(); err != nil {
		return err
	}

	d.emobilityScenario = emobility.NewEMobilitySolution(d.cem.Service, d.cem.Currency, emobility.EmobilityConfiguration{
		CoordinatedChargingEnabled: true,
	})
	d.emobilityScenario.AddFeatures()
	d.emobilityScenario.AddUseCases()

	d.gridScenario = grid.NewGridScenario(d.cem.Service)
	d.gridScenario.AddFeatures()
	d.gridScenario.AddUseCases()

	d.inverterBatteryVisScenario = inverterbatteryvis.NewInverterVisScenario(d.cem.Service)
	d.inverterBatteryVisScenario.AddFeatures()
	d.inverterBatteryVisScenario.AddUseCases()

	d.inverterPVVisScenario = inverterpvvis.NewInverterVisScenario(d.cem.Service)
	d.inverterPVVisScenario.AddFeatures()
	d.inverterPVVisScenario.AddUseCases()

	d.cem.Service.Start()

	return nil
}

// report the Ship ID of a newly trusted connection
func (d *DemoCem) RemoteServiceShipIDReported(service eebusapi.ServiceInterface, ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	logging.Log().Info("SKI", ski, "has Ship ID:", shipID)
}

func (d *DemoCem) RemoteSKIConnected(service eebusapi.ServiceInterface, ski string) {}

func (d *DemoCem) RemoteSKIDisconnected(service eebusapi.ServiceInterface, ski string) {}

func (d *DemoCem) VisibleRemoteServicesUpdated(service eebusapi.ServiceInterface, entries []shipapi.RemoteService) {
}

func (h *DemoCem) ServiceShipIDUpdate(ski string, shipdID string) {}

func (h *DemoCem) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {}

func (h *DemoCem) AllowWaitingForTrust(ski string) bool { return true }

// Logging interface

func (d *DemoCem) log(level string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s", t.Format(time.RFC3339), level, fmt.Sprintln(args...))
}

func (d *DemoCem) logf(level, format string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s\n", t.Format(time.RFC3339), level, fmt.Sprintf(format, args...))
}

func (d *DemoCem) Trace(args ...interface{}) {
	d.log("TRACE", args...)
}

func (d *DemoCem) Tracef(format string, args ...interface{}) {
	d.logf("TRACE", format, args...)
}

func (d *DemoCem) Debug(args ...interface{}) {
	d.log("DEBUG", args...)
}

func (d *DemoCem) Debugf(format string, args ...interface{}) {
	d.logf("DEBUG", format, args...)
}

func (d *DemoCem) Info(args ...interface{}) {
	d.log("INFO", args...)
}

func (d *DemoCem) Infof(format string, args ...interface{}) {
	d.logf("INFO", format, args...)
}

func (d *DemoCem) Error(args ...interface{}) {
	d.log("ERROR", args...)
}

func (d *DemoCem) Errorf(format string, args ...interface{}) {
	d.logf("ERROR", format, args...)
}

// main app
func main() {
	remoteSki := flag.String("remoteski", "", "The remote device SKI")
	port := flag.Int("port", 4815, "Optional port for the EEBUS service")
	crt := flag.String("crt", "cert.crt", "Optional filepath for the cert file")
	key := flag.String("key", "cert.key", "Optional filepath for the key file")
	iface := flag.String("iface", "", "Optional network interface the EEBUS connection should be limited to")

	flag.Parse()

	if len(os.Args) == 1 || remoteSki == nil || *remoteSki == "" {
		flag.Usage()
		return
	}

	certificate, err := tls.LoadX509KeyPair(*crt, *key)
	if err != nil {
		certificate, err = cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-10")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Using certificate file", *crt, "and key file", *key)
	}

	configuration, err := eebusapi.NewConfiguration(
		"Demo",
		"Demo",
		"HEMS",
		"123456789",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		*port,
		certificate,
		230,
		time.Second*4)
	if err != nil {
		fmt.Println("Service data is invalid:", err)
		return
	}

	if iface != nil && *iface != "" {
		ifaces := []string{*iface}

		configuration.SetInterfaces(ifaces)
	}

	demo := NewDemoCem(configuration)
	if err := demo.Setup(); err != nil {
		fmt.Println("Error setting up cem: ", err)
		return
	}

	remoteService := shipapi.NewServiceDetails(*remoteSki)
	demo.emobilityScenario.RegisterRemoteDevice(remoteService, nil)

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}
