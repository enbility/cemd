package grid

type GridInterface interface {
	PowerLimitationFactor() (float64, error)
	MomentaryPowerConsumptionOrProduction() (float64, error)
	TotalFeedInEnergy() (float64, error)
	TotalConsumedEnergy() (float64, error)
	MomentaryCurrentConsumptionOrProduction() ([]float64, error)
	Voltage() ([]float64, error)
	Frequency() (float64, error)
}
