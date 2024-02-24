package cem

import (
	"fmt"
	"testing"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/ucevsecc"
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_CEM(t *testing.T) {
	certificate, err := cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-10")
	assert.Nil(t, err)

	configuration, err := eebusapi.NewConfiguration(
		"Demo",
		"Demo",
		"HEMS",
		"123456789",
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		7654,
		certificate,
		230,
		time.Second*4)
	assert.Nil(t, err)

	demo := &DemoCem{}
	cem := NewCEM(configuration, demo, demo)
	assert.NotNil(t, cem)

	err = cem.Setup()
	assert.Nil(t, err)

	ucEvseCC := ucevsecc.NewUCEVSECC(cem.Service, cem.Service.LocalService(), demo)
	cem.AddUseCase(ucEvseCC)

	cem.Start()
	cem.Shutdown()

}

type DemoCem struct {
}

// UseCaseEventReaderInterface

func (d *DemoCem) SpineEvent(ski string, entity spineapi.EntityRemoteInterface, event api.UseCaseEventType) {

}

// eebusapi.ServiceReaderInterface

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
