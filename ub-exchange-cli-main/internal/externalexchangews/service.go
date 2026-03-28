package externalexchangews

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchangews/binance"
	"exchange-go/internal/externalexchangews/types"
	"exchange-go/internal/platform"
	"exchange-go/internal/processor"
)

// Service provides access to the active external exchange WebSocket implementation.
type Service interface {
	// GetActiveExternalExchangeWs returns the WebSocket client for the currently
	// configured external exchange (e.g. Binance).
	GetActiveExternalExchangeWs() types.ExternalWs
}

type service struct {
	wsClient                   platform.WsClient
	processor                  processor.Processor
	logger                     platform.Logger
	activeExternalExchangeName string
	currencyService            currency.Service
}

var Ws = map[string]func(wsClient platform.WsClient, processor processor.Processor, logger platform.Logger, activePairs []currency.Pair) types.ExternalWs{
	"binance": binance.NewWs,
}

func (s *service) GetActiveExternalExchangeWs() types.ExternalWs {
	ws, exists := Ws[s.activeExternalExchangeName]
	if !exists {
		//we would never reach here just handling the error
		panic("no active external exchange")
	}
	activePairs := s.currencyService.GetActivePairCurrenciesList()
	return ws(s.wsClient, s.processor, s.logger, activePairs)
}

func NewService(wsClient platform.WsClient, processor processor.Processor, logger platform.Logger,
	configs platform.Configs, currencyService currency.Service) Service {
	activeExternalExchangeName := configs.GetActiveExternalExchange()
	return &service{
		wsClient:                   wsClient,
		processor:                  processor,
		logger:                     logger,
		activeExternalExchangeName: activeExternalExchangeName,
		currencyService:            currencyService,
	}
}
