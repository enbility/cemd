package scenarios

import "github.com/enbility/eebus-go/service"

// Implemented by EmobilityScenarioImpl, used by CemImpl
type ScenariosI interface {
	AddFeatures()
	AddUseCases()
}

type ScenarioImpl struct {
	SiteConfig *SiteConfig
	Service    *service.EEBUSService
}

func NewScenarioImpl(siteConfig *SiteConfig, service *service.EEBUSService) *ScenarioImpl {
	return &ScenarioImpl{
		SiteConfig: siteConfig,
		Service:    service,
	}
}

// Generic site specific data
type SiteConfig struct {
	// This is useful when e.g. power values are not available and therefor
	// need to be calculated using the current values
	voltage float64
}

// Create a new site config
// voltage of the electrical installation
func NewSiteConfig(voltage float64) *SiteConfig {
	return &SiteConfig{
		voltage: voltage,
	}
}

func (s *SiteConfig) Voltage() float64 {
	return s.voltage
}
