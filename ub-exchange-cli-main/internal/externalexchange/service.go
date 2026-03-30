package externalexchange

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange/handler"
	"exchange-go/internal/externalexchange/handler/binance"
	"exchange-go/internal/platform"
	"math/rand"
	"time"
)

const (
	TypePrivate    = "private"
	TypePublic     = "public"
	StatusEnabled  = "enabled"
	StatusDisabled = "disabled"

	ExchangeBinance = "binance"
)

type ExternalOrderData struct {
	Pair         string
	Type         string
	Amount       string
	Price        string
	ExchangeType string
}

type OrderPlacementResult struct {
	IsOrderPlaced           bool
	ExternalExchangeOrderID string
	ExternalExchangeID      int64
	Data                    string
	ErrorType               string
	ErrorMessage            string
}

type FetchKlinesParams struct {
	Pair      string
	TimeFrame string
	From      time.Time
	To        time.Time
}

type FetchKlinesResult struct {
	StartTime           time.Time
	EndTime             time.Time
	OpenPrice           string
	HighPrice           string
	LowPrice            string
	ClosePrice          string
	BaseVolume          string
	QuoteVolume         string
	TakerBuyBaseVolume  string
	TakerBuyQuoteVolume string
}

type FetchWithdrawalsParams struct {
	From   time.Time
	Coin   string
	Status string //this is just for test purposes and would not be used in production or dev
}

type FetchWithdrawalsResult struct {
	TxID               string
	Status             string
	ExternalWithdrawID string
	Data               string
}

type FetchOrdersParams struct {
	Pair string
	From int64
}

type FetchOrdersResult struct {
	OrderID       int64
	ClientOrderID string
	ExchangeType  string
	Type          string
	Price         string
	Amount        string
	Status        string
	DateTime      time.Time
	Timestamp     int64
	Data          string
}

type FetchTradesParams struct {
	Pair string
	From int64
}

type FetchTradesResult struct {
	ID         int64
	Price      string
	Amount     string
	Commission string
	Coin       string
	DateTime   time.Time
	Timestamp  int64
	Data       string
	OrderID    int64
}

type WithdrawParams struct {
	Coin      string
	Amount    string
	ToAddress string
	Network   string
}

type WithdrawResult struct {
	ID                 string
	ErrorMessage       string
	ExternalExchangeID int64
}

// Service provides a unified interface to the external exchange (e.g. Binance)
// for order placement, kline fetching, withdrawal operations, and trade history.
type Service interface {
	// OrderPlacement places an order on the external exchange using the given order data.
	OrderPlacement(od ExternalOrderData) (OrderPlacementResult, error)
	// FetchKlines retrieves historical OHLC candlestick data from the external exchange.
	FetchKlines(params FetchKlinesParams) ([]FetchKlinesResult, error)
	// FetchWithdrawals retrieves withdrawal records from the external exchange.
	FetchWithdrawals(params FetchWithdrawalsParams) ([]FetchWithdrawalsResult, error)
	// FetchOrders retrieves order records from the external exchange for a given pair.
	FetchOrders(params FetchOrdersParams) ([]FetchOrdersResult, error)
	// FetchTrades retrieves trade records from the external exchange for a given pair.
	FetchTrades(params FetchTradesParams) ([]FetchTradesResult, error)
	// Withdraw initiates a cryptocurrency withdrawal on the external exchange.
	Withdraw(params WithdrawParams) (WithdrawResult, error)
}

type service struct {
	repo            Repository
	rc              platform.RedisClient
	httpClient      platform.HTTPClient
	pg              currency.PriceGenerator
	configs         platform.Configs
	enabledExchange ExternalExchange
	exchangeHandler handler.ExchangeHandler
	logger          platform.Logger
}

func (s *service) OrderPlacement(od ExternalOrderData) (OrderPlacementResult, error) {
	if s.configs.GetEnv() == platform.EnvTest {
		res := OrderPlacementResult{
			IsOrderPlaced:           true,
			ExternalExchangeOrderID: "id",
			ExternalExchangeID:      1,
			Data:                    "",
			ErrorMessage:            "",
		}
		return res, nil
	}

	currentMarketPrice, err := s.pg.GetPrice(context.Background(), od.Pair)
	if err != nil {
		return OrderPlacementResult{}, err
	}

	params := handler.NewOrderParams{
		Pair:               od.Pair,
		Type:               od.Type,
		Amount:             od.Amount,
		Price:              od.Price,
		CurrentMarketPrice: currentMarketPrice,
		ExchangeType:       od.ExchangeType,
	}

	handlerResult, err := s.exchangeHandler.NewOrder(params)

	or := OrderPlacementResult{
		IsOrderPlaced:           handlerResult.IsOrderPlaced,
		ExternalExchangeID:      s.enabledExchange.ID,
		ExternalExchangeOrderID: handlerResult.ExternalExchangeOrderID,
		Data:                    handlerResult.Data,
		ErrorType:               handlerResult.CustomError.Type,
		ErrorMessage:            handlerResult.CustomError.Message,
	}
	return or, err
}

func (s *service) FetchKlines(params FetchKlinesParams) ([]FetchKlinesResult, error) {
	result := make([]FetchKlinesResult, 0)
	if s.configs.GetEnv() == platform.EnvTest {
		result = []FetchKlinesResult{
			{
				StartTime:           params.From, // for test purposes
				EndTime:             params.To,   // for test purposes
				OpenPrice:           "0.001",
				HighPrice:           "0.001",
				LowPrice:            "0.001",
				ClosePrice:          "0.001",
				BaseVolume:          "0.001",
				QuoteVolume:         "0.001",
				TakerBuyBaseVolume:  "0.001",
				TakerBuyQuoteVolume: "0.001",
			},
		}
		return result, nil
	}

	exchangeParams := handler.FetchKlinesParams{
		Pair:      params.Pair,
		TimeFrame: params.TimeFrame,
		From:      params.From,
		To:        params.To,
	}
	exchangeResult, err := s.exchangeHandler.FetchKlines(exchangeParams)
	if err != nil {
		return result, err
	}

	for _, singleKline := range exchangeResult {
		r := FetchKlinesResult{
			StartTime:           singleKline.StartTime,
			EndTime:             singleKline.EndTime,
			OpenPrice:           singleKline.OpenPrice,
			HighPrice:           singleKline.HighPrice,
			LowPrice:            singleKline.LowPrice,
			ClosePrice:          singleKline.ClosePrice,
			BaseVolume:          singleKline.BaseVolume,
			QuoteVolume:         singleKline.QuoteVolume,
			TakerBuyBaseVolume:  singleKline.TakerBuyBaseVolume,
			TakerBuyQuoteVolume: singleKline.TakerBuyQuoteVolume,
		}
		result = append(result, r)
	}

	return result, nil

}

func (s *service) FetchWithdrawals(params FetchWithdrawalsParams) ([]FetchWithdrawalsResult, error) {
	result := make([]FetchWithdrawalsResult, 0)
	if s.configs.GetEnv() == platform.EnvTest {
		result = []FetchWithdrawalsResult{
			{
				TxID:               "txId",
				Status:             params.Status,
				ExternalWithdrawID: "1",
				Data:               "test",
			},
		}
		return result, nil

	}
	fetchParams := handler.FetchWithdrawalsParams{
		From: params.From,
	}

	exchangeResult, err := s.exchangeHandler.FetchWithdrawals(fetchParams)
	if err != nil {
		return result, err
	}
	for _, item := range exchangeResult {
		r := FetchWithdrawalsResult{
			TxID:               item.TxID,
			Status:             item.Status,
			ExternalWithdrawID: item.ExternalWithdrawID,
			Data:               item.Data,
		}
		result = append(result, r)
	}

	return result, nil
}

func (s *service) FetchOrders(params FetchOrdersParams) ([]FetchOrdersResult, error) {
	result := make([]FetchOrdersResult, 0)
	if s.configs.GetEnv() == platform.EnvTest {
		result = []FetchOrdersResult{
			{
				OrderID:       1,
				ClientOrderID: "1",
				ExchangeType:  "MARKET",
				Type:          "BUY",
				Price:         "30000.00000000",
				Amount:        "1.00000000",
				Status:        "COMPLETED",
				DateTime:      time.Now(),
				Data:          "test",
				Timestamp:     time.Now().Unix() * 1000,
			},
		}
		return result, nil

	}

	fetchParams := handler.FetchOrdersParams{
		Pair: params.Pair,
		From: params.From,
	}

	exchangeResult, err := s.exchangeHandler.FetchOrders(fetchParams)
	if err != nil {
		return result, err
	}
	for _, item := range exchangeResult {
		r := FetchOrdersResult{
			OrderID:       item.OrderID,
			ClientOrderID: item.ClientOrderID,
			ExchangeType:  item.ExchangeType,
			Type:          item.Type,
			Price:         item.Price,
			Amount:        item.Amount,
			Status:        item.Status,
			DateTime:      item.DateTime,
			Timestamp:     item.Timestamp,
			Data:          item.Data,
		}
		result = append(result, r)
	}

	return result, nil
}

func (s *service) FetchTrades(params FetchTradesParams) ([]FetchTradesResult, error) {
	result := make([]FetchTradesResult, 0)
	if s.configs.GetEnv() == platform.EnvTest {
		min := 111111
		max := 999999
		id := rand.Int63n(int64(max-min+1)) + int64(min)
		result = []FetchTradesResult{
			{
				ID:         id,
				Price:      "30000.00000000",
				Amount:     "1.00000000",
				Commission: "0.00010000",
				Coin:       "BTC",
				DateTime:   time.Time{},
				Data:       "test",
			},
		}
		return result, nil

	}

	fetchParams := handler.FetchTradesParams{
		Pair: params.Pair,
		From: params.From,
	}

	exchangeResult, err := s.exchangeHandler.FetchTrades(fetchParams)
	if err != nil {
		return result, err
	}
	for _, item := range exchangeResult {
		r := FetchTradesResult{
			ID:         item.ID,
			Price:      item.Price,
			Amount:     item.Amount,
			Commission: item.Commission,
			Coin:       item.Coin,
			DateTime:   item.DateTime,
			Timestamp:  item.Timestamp,
			Data:       item.Data,
			OrderID:    item.OrderID,
		}
		result = append(result, r)
	}

	return result, nil
}

func (s *service) Withdraw(params WithdrawParams) (WithdrawResult, error) {
	if s.configs.GetEnv() == platform.EnvTest {
		return WithdrawResult{
			ID: params.ToAddress,
		}, nil
	}
	withdrawParams := handler.WithdrawParams{
		Coin:      params.Coin,
		Amount:    params.Amount,
		ToAddress: params.ToAddress,
		Network:   params.Network,
	}

	handlerResult, err := s.exchangeHandler.Withdraw(withdrawParams)

	result := WithdrawResult{
		ID:                 handlerResult.ID,
		ErrorMessage:       handlerResult.ErrorMessage,
		ExternalExchangeID: s.enabledExchange.ID,
	}
	return result, err
}

func NewExternalExchangeService(repo Repository, rc platform.RedisClient, httpClient platform.HTTPClient, pg currency.PriceGenerator, configs platform.Configs, logger platform.Logger) Service {
	s := &service{
		repo:       repo,
		rc:         rc,
		httpClient: httpClient,
		pg:         pg,
		configs:    configs,
		logger:     logger,
	}
	s.setExtraService()
	return s
}

func (s *service) setExtraService() {
	ee := ExternalExchange{}
	_ = s.repo.GetEnabledPrivateExternalExchange(&ee)
	s.enabledExchange = ee
	s.exchangeHandler = s.getExchangeHandler(ee.Name, ee.MetaData)
}

func (s *service) getExchangeHandler(name string, metadata string) handler.ExchangeHandler {
	var h handler.ExchangeHandler
	var err error
	switch name {
	case ExchangeBinance:
		h, err = binance.NewBinanceService(s.httpClient, s.rc, s.configs, s.logger, metadata)
	default:
		h, err = binance.NewBinanceService(s.httpClient, s.rc, s.configs, s.logger, metadata)
	}
	if err != nil {
		s.logger.Error2("getExchangeHandler: failed to create exchange handler", err)
		return nil
	}
	return h
}
