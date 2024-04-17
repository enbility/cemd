package democem

import (
	eebusapi "github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
)

// report the Ship ID of a newly trusted connection
func (d *DemoCem) RemoteSKIConnected(service eebusapi.ServiceInterface, ski string) {}

func (d *DemoCem) RemoteSKIDisconnected(service eebusapi.ServiceInterface, ski string) {}

func (d *DemoCem) VisibleRemoteServicesUpdated(service eebusapi.ServiceInterface, entries []shipapi.RemoteService) {
}

func (d *DemoCem) ServiceShipIDUpdate(ski string, shipdID string) {}

func (d *DemoCem) ServicePairingDetailUpdate(ski string, detail *shipapi.ConnectionStateDetail) {}
