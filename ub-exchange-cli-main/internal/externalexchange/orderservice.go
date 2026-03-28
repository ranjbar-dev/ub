package externalexchange

import (
	"database/sql"
	"exchange-go/internal/platform"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const (
	StatusOpen                         = "OPEN"
	StatusCompleted                    = "COMPLETED"
	StatusFailed                       = "FAILED"
	TypeSell                           = "SELL"
	TypeBuy                            = "BUY"
	OrderFailReasonRateLimit           = "RATE_LIMIT"
	OrderFailReasonInsufficientBalance = "INSUFFICIENT_BALANCE"
	OrderFailReasonExchangeError       = "EXCHANGE_ERROR"
	OrderSourceUser                    = "USER"
	OrderSourceBot                     = "BOT"
	OrderSourceAdmin                   = "ADMIN"
)

type BotOrderParams struct {
	PairID       int64
	PairName     string
	Type         string
	ExchangeType string
	Amount       string
	Price        string
	BuyAmount    string
	BuyPrice     string
	SellAmount   string
	SellPrice    string
	LastTradeID  int64
	OrderIds     []string
}

type UserOrderParams struct {
	PairID       int64
	PairName     string
	Type         string
	ExchangeType string
	Amount       string
	Price        string
	OrderID      int64
}

type UserOrderResult struct {
	IsOrderPlaced           bool
	ExternalExchangeOrderID string
	ExternalExchangeID      int64
	Data                    string
}

type LastTradeIDAndPair struct {
	PairID  int64
	TradeID int64
}

// OrderService manages order placement on the external exchange for both user
// orders and aggregated bot orders.
type OrderService interface {
	// CreateExternalExchangeOrderForBot places an aggregated bot order on the external
	// exchange and records it locally.
	CreateExternalExchangeOrderForBot(o BotOrderParams) (isOrderPlaced bool, err error)
	// CreateExternalExchangeOrderForUser places a user-initiated order on the external
	// exchange and records it locally.
	CreateExternalExchangeOrderForUser(o UserOrderParams) (result UserOrderResult, err error)
	// GetExternalExchangeOrdersLastTradeIds returns the last trade ID per pair for
	// reconciliation with the external exchange.
	GetExternalExchangeOrdersLastTradeIds() []LastTradeIDAndPair
}

type orderService struct {
	externalExchangeOrderRepository OrderRepository
	externalExchangeService         Service
	logger                          platform.Logger
}

func (s *orderService) CreateExternalExchangeOrderForBot(params BotOrderParams) (isOrderPlaced bool, err error) {
	orderIdsString := strings.Join(params.OrderIds, ",")
	order := &Order{
		PairID:       sql.NullInt64{Int64: params.PairID, Valid: true},
		Type:         sql.NullString{String: params.Type, Valid: true},
		ExchangeType: sql.NullString{String: params.ExchangeType, Valid: true},
		Price:        sql.NullString{String: params.Price, Valid: true},
		Amount:       sql.NullString{String: params.Amount, Valid: true},
		Status:       sql.NullString{String: StatusOpen, Valid: true},
		LastTradeID:  sql.NullInt64{Int64: params.LastTradeID, Valid: true},
		BuyAmount:    sql.NullString{String: params.BuyAmount, Valid: true},
		BuyPrice:     sql.NullString{String: params.BuyPrice, Valid: true},
		SellAmount:   sql.NullString{String: params.SellAmount, Valid: true},
		SellPrice:    sql.NullString{String: params.SellPrice, Valid: true},
		OrderIds:     sql.NullString{String: orderIdsString, Valid: true},
		Source:       sql.NullString{String: OrderSourceBot, Valid: true},
	}
	err = s.externalExchangeOrderRepository.Create(order)
	if err != nil {
		return false, err
	}

	od := ExternalOrderData{
		Pair:         params.PairName,
		Type:         params.Type,
		Amount:       params.Amount,
		Price:        params.Price,
		ExchangeType: params.ExchangeType,
	}

	placementResult, err := s.externalExchangeService.OrderPlacement(od)
	if err != nil {
		s.logger.Error2("error in placing order", err,
			zap.String("service", "externalExchangeOrderService"),
			zap.String("method", "CreateExternalExchangeOrderForBot"),
			zap.String("pairName", params.PairName),
			zap.String("amount", params.Amount),
			zap.String("type", params.Type),
		)
	}

	if err == nil && placementResult.IsOrderPlaced {
		isOrderPlaced = true
		order.ExchangeID = sql.NullInt64{Int64: placementResult.ExternalExchangeID, Valid: true}
		order.MetaID = sql.NullString{String: placementResult.ExternalExchangeOrderID, Valid: true}
		order.Status = sql.NullString{String: StatusCompleted, Valid: true}
	} else {
		isOrderPlaced = false
		if err == nil && !placementResult.IsOrderPlaced {
			err = fmt.Errorf("order is not placed in external exchange for pair %s ", params.PairName)
		}

		order.FailReason = sql.NullString{String: placementResult.ErrorType, Valid: true}
		order.ExceptionMessage = sql.NullString{String: placementResult.ErrorMessage, Valid: true}
		order.Status = sql.NullString{String: StatusFailed, Valid: true}
	}
	err = s.externalExchangeOrderRepository.Update(order)
	return isOrderPlaced, err
}

func (s *orderService) CreateExternalExchangeOrderForUser(params UserOrderParams) (result UserOrderResult, err error) {
	buyAmount := sql.NullString{String: "", Valid: false}
	buyPrice := sql.NullString{String: "", Valid: false}
	sellAmount := sql.NullString{String: "", Valid: false}
	sellPrice := sql.NullString{String: "", Valid: false}

	if params.Type == TypeBuy {
		buyAmount = sql.NullString{String: params.Amount, Valid: true}
		buyPrice = sql.NullString{String: params.Price, Valid: true}
	} else {
		sellAmount = sql.NullString{String: params.Amount, Valid: false}
		sellPrice = sql.NullString{String: params.Price, Valid: false}
	}

	order := &Order{
		PairID:       sql.NullInt64{Int64: params.PairID, Valid: true},
		Type:         sql.NullString{String: params.Type, Valid: true},
		ExchangeType: sql.NullString{String: params.ExchangeType, Valid: true},
		Price:        sql.NullString{String: params.Price, Valid: true},
		Amount:       sql.NullString{String: params.Amount, Valid: true},
		Status:       sql.NullString{String: StatusOpen, Valid: true},
		BuyAmount:    buyAmount,
		BuyPrice:     buyPrice,
		SellAmount:   sellAmount,
		SellPrice:    sellPrice,
		Source:       sql.NullString{String: OrderSourceUser, Valid: true},
		UserOrderID:  sql.NullInt64{Int64: params.OrderID, Valid: true},
	}
	err = s.externalExchangeOrderRepository.Create(order)
	if err != nil {
		return result, err
	}

	od := ExternalOrderData{
		Pair:         params.PairName,
		Type:         params.Type,
		Amount:       params.Amount,
		Price:        params.Price,
		ExchangeType: params.ExchangeType,
	}

	placementResult, err := s.externalExchangeService.OrderPlacement(od)
	if err != nil {
		s.logger.Error2("error in placing order", err,
			zap.String("service", "externalExchangeOrderService"),
			zap.String("method", "CreateExternalExchangeOrderForUser"),
			zap.String("pairName", params.PairName),
			zap.String("amount", params.Amount),
			zap.String("type", params.Type),
		)
	}

	if err == nil && placementResult.IsOrderPlaced {
		order.ExchangeID = sql.NullInt64{Int64: placementResult.ExternalExchangeID, Valid: true}
		order.MetaID = sql.NullString{String: placementResult.ExternalExchangeOrderID, Valid: true}
		order.Status = sql.NullString{String: StatusCompleted, Valid: true}
	} else {
		if err == nil && !placementResult.IsOrderPlaced {
			err = fmt.Errorf("order is not placed in external exchange for pair %s ", params.PairName)
		}
		order.FailReason = sql.NullString{String: placementResult.ErrorType, Valid: true}
		order.ExceptionMessage = sql.NullString{String: placementResult.ErrorMessage, Valid: true}
		order.Status = sql.NullString{String: StatusFailed, Valid: true}
	}
	err = s.externalExchangeOrderRepository.Update(order)

	result = UserOrderResult{
		IsOrderPlaced:           placementResult.IsOrderPlaced,
		ExternalExchangeOrderID: placementResult.ExternalExchangeOrderID,
		ExternalExchangeID:      placementResult.ExternalExchangeID,
		Data:                    placementResult.Data,
	}

	return result, err
}

func (s *orderService) GetExternalExchangeOrdersLastTradeIds() []LastTradeIDAndPair {
	return s.externalExchangeOrderRepository.GetExternalExchangeOrdersLastTradeIds()
}

func NewOrderService(externalExchangeOrderRepository OrderRepository, externalExchangeService Service,
	logger platform.Logger) OrderService {
	return &orderService{
		externalExchangeOrderRepository: externalExchangeOrderRepository,
		externalExchangeService:         externalExchangeService,
		logger:                          logger,
	}
}
