package democem

import (
	"github.com/enbility/cemd/cem"
	"github.com/enbility/cemd/ucevsecc"
	eebusapi "github.com/enbility/eebus-go/api"
)

type DemoCem struct {
	cem *cem.Cem
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

	evsecc := ucevsecc.NewUCEVSECC(d.cem.Service, d.cem.Service.LocalService(), d)
	d.cem.AddUseCase(evsecc)

	d.cem.Start()

	return nil
}
