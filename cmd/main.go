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

	"github.com/enbility/cemd/cmd/democem"
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/spine-go/model"
)

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

	demo := democem.NewDemoCem(configuration)
	if err := demo.Setup(); err != nil {
		fmt.Println("Error setting up cem: ", err)
		return
	}

	// remoteService := shipapi.NewServiceDetails(*remoteSki)
	// demo.emobilityScenario.RegisterRemoteDevice(remoteService, nil)

	// Clean exit to make sure mdns shutdown is invoked
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	// User exit
}
