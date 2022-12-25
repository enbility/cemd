package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/enbility/cemd/cem"
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine/model"
)

type DemoCem struct {
	cem *cem.CemImpl
}

func NewDemoCem(configuration *service.Configuration) *DemoCem {
	demo := &DemoCem{}

	demo.cem = cem.NewCEM(configuration, demo, demo)

	return demo
}

func (d *DemoCem) Setup() error {
	return d.cem.Setup(cem.CemConfiguration{
		EmobilityScenarioEnabled: true,
		GridScenarioEnabled:      true,
	})
}

// report the Ship ID of a newly trusted connection
func (d *DemoCem) RemoteServiceShipIDReported(service *service.EEBUSService, ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	logging.Log.Info("SKI", ski, "has Ship ID:", shipID)
}

func (d *DemoCem) RemoteSKIConnected(service *service.EEBUSService, ski string) {}

func (d *DemoCem) RemoteSKIDisconnected(service *service.EEBUSService, ski string) {}

func (h *DemoCem) ReportServiceShipID(ski string, shipdID string) {}

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
func usage() {
	fmt.Println("Usage: go run /cmd/main.go <serverport> <evse-ski> <crtfile> <keyfile> <iface>")
}

func main() {
	if len(os.Args) < 5 {
		usage()
		return
	}

	portValue, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Port is invalid:", err)
		return
	}

	certificate, err := tls.LoadX509KeyPair(os.Args[3], os.Args[4])
	if err != nil {
		fmt.Println("Certificate is invalid:", err)
		return
	}

	ifaces := []string{os.Args[5]}

	configuration, err := service.NewConfiguration(
		"Demo",
		"Demo",
		"HEMS",
		"123456789",
		model.DeviceTypeTypeEnergyManagementSystem,
		portValue,
		certificate,
		230)
	if err != nil {
		fmt.Println("Service data is invalid:", err)
		return
	}
	configuration.SetInterfaces(ifaces)

	demo := NewDemoCem(configuration)
	if err := demo.Setup(); err != nil {
		fmt.Println("Error setting up cem: ", err)
		return
	}

	remoteService := service.NewServiceDetails(os.Args[2])
	demo.cem.RegisterEmobilityRemoteDevice(remoteService, nil)

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}
