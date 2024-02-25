package democem

import (
	"github.com/enbility/cemd/cem"
	"github.com/enbility/cemd/ucevsecc"
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/logging"
)

type DemoCem struct {
	cem *cem.Cem
}

func NewDemoCem(configuration *eebusapi.Configuration) *DemoCem {
	demo := &DemoCem{}

	noLogging := &logging.NoLogging{}
	demo.cem = cem.NewCEM(configuration, demo, demo, noLogging)

	return demo
}

func (d *DemoCem) Setup() error {
	if err := d.cem.Setup(); err != nil {
		return err
	}

	evsecc := ucevsecc.NewUCEVSECC(d.cem.Service, d.cem.Service.LocalService(), d)
	d.cem.AddUseCase(evsecc)

	d.cem.Start()

	return nil
}
