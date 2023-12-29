package cem

import (
	"github.com/enbility/eebus-go/logging"
	"github.com/enbility/eebus-go/service"
	"github.com/enbility/eebus-go/spine/model"
)

// Generic CEM implementation
type CemImpl struct {
	Service *service.EEBUSService

	Currency model.CurrencyType
}

func NewCEM(serviceDescription *service.Configuration, serviceHandler service.EEBUSServiceHandler, log logging.Logging) *CemImpl {
	cem := &CemImpl{
		Service:  service.NewEEBUSService(serviceDescription, serviceHandler),
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
