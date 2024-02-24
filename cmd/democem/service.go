package democem

import (
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
)

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
