package cem

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/model"
)

// Generic CEM implementation
type CemImpl struct {
	Service api.ServiceInterface

	Currency model.CurrencyType
}

func NewCEM(serviceDescription *api.Configuration, serviceHandler api.ServiceReaderInterface, log logging.LoggingInterface) *CemImpl {
	cem := &CemImpl{
		Service:  service.NewService(serviceDescription, serviceHandler),
		Currency: model.CurrencyTypeEur,
	}

	cem.Service.SetLogging(log)

	return cem
}

// Set up the supported usecases and features
func (h *CemImpl) Setup() error {
	if err := h.Service.Setup(); err != nil {
		return err
	}

	return nil
}

func (h *CemImpl) Start() {
	h.Service.Start()
}

func (h *CemImpl) Shutdown() {
	h.Service.Shutdown()
}
