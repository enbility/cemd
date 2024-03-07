package cem

import (
	"testing"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/ucevsecc"
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCemSuite(t *testing.T) {
	suite.Run(t, new(CemSuite))
}

type CemSuite struct {
	suite.Suite

	sut              *Cem
	mockRemoteDevice *mocks.DeviceRemoteInterface
}

func (s *CemSuite) BeforeTest(suiteName, testName string) {
	s.mockRemoteDevice = mocks.NewDeviceRemoteInterface(s.T())

	certificate, err := cert.CreateCertificate("Demo", "Demo", "DE", "Demo-Unit-10")
	assert.Nil(s.T(), err)

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
	assert.Nil(s.T(), err)

	noLogging := &logging.NoLogging{}
	s.sut = NewCEM(configuration, s, s.deviceEventCB, noLogging)
	assert.NotNil(s.T(), s.sut)
}
func (s *CemSuite) Test_CEM() {
	err := s.sut.Setup()
	assert.Nil(s.T(), err)

	ucEvseCC := ucevsecc.NewUCEVSECC(s.sut.Service, s.entityEventCB)
	s.sut.AddUseCase(ucEvseCC)

	s.sut.Start()
	s.sut.Shutdown()
}

// Callbacks
func (d *CemSuite) deviceEventCB(ski string, device spineapi.DeviceRemoteInterface, event api.EventType) {
}

func (d *CemSuite) entityEventCB(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
}

// eebusapi.ServiceReaderInterface

// report the Ship ID of a newly trusted connection
func (d *CemSuite) RemoteServiceShipIDReported(service eebusapi.ServiceInterface, ski string, shipID string) {
	// we should associated the Ship ID with the SKI and store it
	// so the next connection can start trusted
	logging.Log().Info("SKI", ski, "has Ship ID:", shipID)
}

func (d *CemSuite) RemoteSKIConnected(service eebusapi.ServiceInterface, ski string) {}

func (d *CemSuite) RemoteSKIDisconnected(service eebusapi.ServiceInterface, ski string) {}

func (d *CemSuite) VisibleRemoteServicesUpdated(service eebusapi.ServiceInterface, entries []shipapi.RemoteService) {
}

func (h *CemSuite) ServiceShipIDUpdate(ski string, shipdID string) {}

func (h *CemSuite) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {}

func (h *CemSuite) AllowWaitingForTrust(ski string) bool { return true }
