package order

import (
	"exchange-go/internal/engine"
	"fmt"

	"github.com/shopspring/decimal"
)

// EngineCommunicator is the bridge between the order service and the matching engine,
// translating domain orders into engine-specific representations.
type EngineCommunicator interface {
	// SubmitOrder converts the order to an engine order and pushes it into the matching engine queue.
	SubmitOrder(o Order) error
	// RemoveOrder cancels an order by removing it from the matching engine's order book.
	RemoveOrder(o Order) error
	// RetrieveOrder fetches the current state of an order from the matching engine.
	RetrieveOrder(o Order) error
}

type engineCommunicator struct {
	engine      engine.Engine
	forceTrader ForceTrader
}

func (ec *engineCommunicator) SubmitOrder(o Order) error {
	engineOrder, err := ec.getEngineOrder(o)
	if err != nil {
		return err
	}
	return ec.engine.SubmitOrder(engineOrder)
}

func (ec *engineCommunicator) RemoveOrder(o Order) error {
	engineOrder, err := ec.getEngineOrder(o)
	if err != nil {
		return err
	}
	err = ec.engine.RemoveOrder(engineOrder)
	return err

}

func (ec *engineCommunicator) RetrieveOrder(o Order) error {
	engineOrder, err := ec.getEngineOrder(o)
	if err != nil {
		return err
	}
	return ec.engine.RetrieveOrder(engineOrder)
}

func (ec *engineCommunicator) getEngineOrder(o Order) (engine.Order, error) {
	engineOrder := engine.Order{}
	id := o.getStringID()
	//we are doing this so all ids have same len, because in engine we need lexicographical order
	id = fmt.Sprintf("%011s", id)
	side := mapTypeToSide(o.Type)
	quantityDecimal, err := decimal.NewFromString(o.getAmount())
	if err != nil {
		return engineOrder, err
	}
	quantity := quantityDecimal.StringFixed(8)
	price := o.Price.String
	minPrice, maxPrice, err := ec.forceTrader.GetMinAndMaxPrice(o.Pair.Name, o.Type, o.CurrentMarketPrice.String)
	if err != nil {
		return engineOrder, err
	}
	marketPriceDecimal, err := decimal.NewFromString(o.CurrentMarketPrice.String)
	if err != nil {
		return engineOrder, err
	}
	marketPrice := marketPriceDecimal.StringFixed(8)
	engineOrder = engine.Order{
		Pair:              o.Pair.Name,
		ID:                id,
		Side:              side,
		Quantity:          quantity,
		Price:             price,
		Timestamp:         o.CreatedAt.Unix(),
		MarketPrice:       marketPrice,
		MinThresholdPrice: minPrice,
		MaxThresholdPrice: maxPrice,
		UserID:            fmt.Sprint(o.UserID),
	}

	return engineOrder, nil
}

func NewEngineCommunicator(ft ForceTrader, e engine.Engine) EngineCommunicator {
	return &engineCommunicator{
		engine:      e,
		forceTrader: ft,
	}

}
