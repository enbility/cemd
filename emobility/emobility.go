package emobility

import (
	"github.com/enbility/eebus-go/api"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type EMobility struct {
	entity spineapi.EntityLocalInterface

	service api.ServiceInterface

	evseEntity spineapi.EntityRemoteInterface
	evEntity   spineapi.EntityRemoteInterface

	evCurrentChargeStrategy EVChargeStrategyType

	ski      string
	currency model.CurrencyType

	configuration EmobilityConfiguration
	dataProvider  EmobilityDataProvider
}

var _ EMobilityInterface = (*EMobility)(nil)

// Add E-Mobility support
func NewEMobility(service api.ServiceInterface, details *shipapi.ServiceDetails, currency model.CurrencyType, configuration EmobilityConfiguration, dataProvider EmobilityDataProvider) *EMobility {
	ski := util.NormalizeSKI(details.SKI())

	localEntity := service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

	emobility := &EMobility{
		service:                 service,
		entity:                  localEntity,
		ski:                     ski,
		currency:                currency,
		dataProvider:            dataProvider,
		evCurrentChargeStrategy: EVChargeStrategyTypeUnknown,
		configuration:           configuration,
	}
	_ = spine.Events.Subscribe(emobility)

	service.RegisterRemoteSKI(ski, true)

	return emobility
}
