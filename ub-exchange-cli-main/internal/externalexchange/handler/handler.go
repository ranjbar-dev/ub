package handler

import (
	"time"
)

var (
	ErrTypeInvalidParams       = "INVALID_PARAMS"
	ErrTypeInsufficientBalance = "INSUFFICIENT_BALANCE"
	ErrTypeRateLimit           = "RATE_LIMIT"
	ErrTypeUnknown             = "UNKNOWN"
)

// ExchangeHandler defines a pluggable interface for interacting with a specific
// external exchange implementation (e.g. Binance). It abstracts order placement,
// data fetching, and withdrawal operations.
type ExchangeHandler interface {
	// NewOrder places a new order on the external exchange.
	NewOrder(params NewOrderParams) (NewOrderResult, error)
	// FetchKlines retrieves historical OHLC candlestick data from the exchange.
	FetchKlines(params FetchKlinesParams) ([]FetchKlinesResult, error)
	// FetchWithdrawals retrieves withdrawal records from the exchange.
	FetchWithdrawals(params FetchWithdrawalsParams) ([]FetchWithdrawalsResult, error)
	// FetchOrders retrieves order records from the exchange for a given pair.
	FetchOrders(params FetchOrdersParams) ([]FetchOrdersResult, error)
	// FetchTrades retrieves trade execution records from the exchange for a given pair.
	FetchTrades(params FetchTradesParams) ([]FetchTradesResult, error)
	// Withdraw initiates a cryptocurrency withdrawal on the exchange.
	Withdraw(params WithdrawParams) (WithdrawResult, error)
}

type NewOrderParams struct {
	Pair               string
	Type               string
	Amount             string
	Price              string
	CurrentMarketPrice string
	ExchangeType       string
}

type CustomError struct {
	Type    string
	Message string
}

func (ce CustomError) Error() string {
	return ce.Message
}

type NewOrderResult struct {
	IsOrderPlaced           bool
	ExternalExchangeOrderID string
	Data                    string
	CustomError             CustomError
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
	From time.Time
	Coin string
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
	ID           string
	ErrorMessage string
}
