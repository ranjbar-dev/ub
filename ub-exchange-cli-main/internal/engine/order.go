package engine

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Order struct {
	Pair                 string `json:"pair"`
	ID                   string `json:"id"`
	Side                 string `json:"side"`
	Quantity             string `json:"quantity"`
	Price                string `json:"price"`
	Timestamp            int64  `json:"timestamp"`
	TradedWithOrderID    string `json:"-"`
	QuantityTraded       string `json:"-"`
	TradePrice           string `json:"-"`
	IsAlreadyInOrderBook bool   `json:"-"`
	MarketPrice          string `json:"marketPrice"`
	MinThresholdPrice    string `json:"minPrice"`
	MaxThresholdPrice    string `json:"maxPrice"`
}

func (o Order) GetPrice() (decimal.Decimal, error) {
	return decimal.NewFromString(o.Price)
}

func (o Order) GetMinPrice() (decimal.Decimal, error) {
	return decimal.NewFromString(o.MinThresholdPrice)
}

func (o Order) GetMaxPrice() (decimal.Decimal, error) {
	return decimal.NewFromString(o.MaxThresholdPrice)
}

func (o Order) GetQuantity() (decimal.Decimal, error) {
	return decimal.NewFromString(o.Quantity)
}

func (o Order) GetMarketPrice() (decimal.Decimal, error) {
	return decimal.NewFromString(o.MarketPrice)
}

func (o Order) IsBidMarket() bool {
	return o.Price == "" && o.Side == SideBid
}

//the only difference is here we do not have max and minPrice
func (o Order) MarshalForOrderbook() ([]byte, error) {
	return json.Marshal(&struct {
		Pair        string `json:"pair"`
		ID          string `json:"id"`
		Side        string `json:"side"`
		Quantity    string `json:"quantity"`
		Price       string `json:"price"`
		MarketPrice string `json:"marketPrice"`
		Timestamp   int64  `json:"timestamp"`
	}{
		Pair:        o.Pair,
		ID:          o.ID,
		Side:        o.Side,
		Quantity:    o.Quantity,
		Price:       o.Price,
		MarketPrice: o.MarketPrice,
		Timestamp:   o.Timestamp,
	})
}

func newOrder(pair string, id string, side string, quantity string, price string, marketPrice string, timestamp int64) Order {
	return Order{
		Pair:        pair,
		ID:          id,
		Side:        side,
		Quantity:    quantity,
		Price:       price,
		MarketPrice: marketPrice,
		Timestamp:   timestamp,
	}
}
