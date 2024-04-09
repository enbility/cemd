package democem

import (
	"github.com/enbility/cemd/cem"
	"github.com/enbility/cemd/ucevsecc"
	"github.com/enbility/cemd/uclpcserver"
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/logging"
)

type DemoCem struct {
	cem *cem.Cem
}

func NewDemoCem(configuration *eebusapi.Configuration) *DemoCem {
	demo := &DemoCem{}

	noLogging := &logging.NoLogging{}
	demo.cem = cem.NewCEM(configuration, demo, demo.deviceEventCB, noLogging)

	return demo
}

func (d *DemoCem) Setup() error {
	if err := d.cem.Setup(); err != nil {
		return err
	}

	lpcs := uclpcserver.NewUCLPC(d.cem.Service, d.entityEventCB)
	d.cem.AddUseCase(lpcs)

	evsecc := ucevsecc.NewUCEVSECC(d.cem.Service, d.entityEventCB)
	d.cem.AddUseCase(evsecc)

	d.cem.Start()

	return nil
}
