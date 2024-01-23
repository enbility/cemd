package inverterpvvis

type InverterPVVisInterface interface {
	CurrentProductionPower() (float64, error)
	NominalPeakPower() (float64, error)
	TotalPVYield() (float64, error)
}
