package inverterbatteryvis

type InverterBatteryVisInterface interface {
	CurrentDisChargePower() (float64, error)
	TotalChargeEnergy() (float64, error)
	TotalDischargeEnergy() (float64, error)
	CurrentStateOfCharge() (float64, error)
}
