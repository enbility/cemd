package democem

import (
	"fmt"
	"time"

	"github.com/enbility/cemd/api"
	"github.com/enbility/cemd/cem"
	"github.com/enbility/cemd/ucevsecc"
	"github.com/enbility/cemd/uclpcserver"
	"github.com/enbility/cemd/uclppserver"
	eebusapi "github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/logging"
)

type DemoCem struct {
	cem *cem.Cem

	remoteSki string
}

func NewDemoCem(configuration *eebusapi.Configuration, remoteSki string) *DemoCem {
	demo := &DemoCem{
		remoteSki: remoteSki,
	}

	demo.cem = cem.NewCEM(configuration, demo, demo.deviceEventCB, demo)

	return demo
}

func (d *DemoCem) Setup() error {
	if err := d.cem.Setup(); err != nil {
		return err
	}

	lpcs := uclpcserver.NewUCLPC(d.cem.Service, d.entityEventCB)
	d.cem.AddUseCase(lpcs)

	if err := lpcs.SetConsumptionLimit(api.LoadLimit{
		IsChangeable: true,
		IsActive:     false,
		Value:        0,
	}); err != nil {
		logging.Log().Error(err)
	}
	if err := lpcs.SetContractualConsumptionNominalMax(22000); err != nil {
		logging.Log().Error(err)
	}
	if err := lpcs.SetFailsafeConsumptionActivePowerLimit(4300, true); err != nil {
		logging.Log().Error(err)
	}
	if err := lpcs.SetFailsafeDurationMinimum(time.Hour*2, true); err != nil {
		logging.Log().Error(err)
	}

	lpps := uclppserver.NewUCLPP(d.cem.Service, d.entityEventCB)
	d.cem.AddUseCase(lpps)

	if err := lpps.SetProductionLimit(api.LoadLimit{
		IsChangeable: true,
		IsActive:     false,
		Value:        0,
	}); err != nil {
		logging.Log().Error(err)
	}
	if err := lpps.SetContractualProductionNominalMax(-7000); err != nil {
		logging.Log().Error(err)
	}
	if err := lpps.SetFailsafeProductionActivePowerLimit(0, true); err != nil {
		logging.Log().Error(err)
	}
	if err := lpps.SetFailsafeDurationMinimum(time.Hour*2, true); err != nil {
		logging.Log().Error(err)
	}

	evsecc := ucevsecc.NewUCEVSECC(d.cem.Service, d.entityEventCB)
	d.cem.AddUseCase(evsecc)

	d.cem.Service.RegisterRemoteSKI(d.remoteSki)

	d.cem.Start()

	return nil
}

// Logging interface

func (d *DemoCem) Trace(args ...interface{}) {
	d.print("TRACE", args...)
}

func (d *DemoCem) Tracef(format string, args ...interface{}) {
	d.printFormat("TRACE", format, args...)
}

func (d *DemoCem) Debug(args ...interface{}) {
	d.print("DEBUG", args...)
}

func (d *DemoCem) Debugf(format string, args ...interface{}) {
	d.printFormat("DEBUG", format, args...)
}

func (d *DemoCem) Info(args ...interface{}) {
	d.print("INFO ", args...)
}

func (d *DemoCem) Infof(format string, args ...interface{}) {
	d.printFormat("INFO ", format, args...)
}

func (d *DemoCem) Error(args ...interface{}) {
	d.print("ERROR", args...)
}

func (d *DemoCem) Errorf(format string, args ...interface{}) {
	d.printFormat("ERROR", format, args...)
}

func (d *DemoCem) currentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (d *DemoCem) print(msgType string, args ...interface{}) {
	value := fmt.Sprintln(args...)
	fmt.Printf("%s %s %s", d.currentTimestamp(), msgType, value)
}

func (d *DemoCem) printFormat(msgType, format string, args ...interface{}) {
	value := fmt.Sprintf(format, args...)
	fmt.Println(d.currentTimestamp(), msgType, value)
}
