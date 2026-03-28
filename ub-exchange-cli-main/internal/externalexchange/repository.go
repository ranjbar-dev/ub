package externalexchange

import (
	"database/sql"
	"time"
)

type ExternalExchange struct {
	ID       int64
	Name     string
	MetaData string
	Status   string
	Type     string
}

// Repository provides data access for external exchange configuration records.
type Repository interface {
	// GetEnabledPrivateExternalExchange retrieves the currently enabled private external exchange.
	GetEnabledPrivateExternalExchange(ee *ExternalExchange) error
}

type Order struct {
	ID                 int64
	PairID             sql.NullInt64 `gorm:"column:pair_currency_id"`
	ExchangeID         sql.NullInt64 `gorm:"column:external_exchange_id"`
	Type               sql.NullString
	ExchangeType       sql.NullString
	Price              sql.NullString
	Amount             sql.NullString
	OtherInfo          sql.NullString `gorm:"column:external_exchange_other_info"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Status             sql.NullString
	FinalGetAmount     sql.NullString
	FinalPayAmount     sql.NullString
	FinalTradePrice    sql.NullString
	FinalFeePercentage sql.NullFloat64
	LastTradeID        sql.NullInt64
	FailReason         sql.NullString
	BuyAmount          sql.NullString
	BuyPrice           sql.NullString
	SellAmount         sql.NullString
	SellPrice          sql.NullString
	OrderIds           sql.NullString
	Source             sql.NullString
	ExceptionMessage   sql.NullString
	MetaID             sql.NullString
	UserOrderID        sql.NullInt64
}

func (Order) TableName() string {
	return "external_exchange_orders"
}

// OrderRepository manages external exchange order records in the local database.
type OrderRepository interface {
	// Create inserts a new external exchange order record.
	Create(o *Order) error
	// Update persists changes to an existing external exchange order record.
	Update(o *Order) error
	// GetExternalExchangeOrdersLastTradeIds returns the last trade ID per pair
	// across all external exchange orders.
	GetExternalExchangeOrdersLastTradeIds() []LastTradeIDAndPair
}

type OrderFromExternal struct {
	ID              int64
	PairID          sql.NullInt64 `gorm:"column:pair_currency_id"`
	ExternalOrderID int64
	ClientOrderID   string
	Type            string `gorm:"column:side"`
	ExchangeType    string `gorm:"column:type"`
	Price           sql.NullString
	Amount          sql.NullString
	Status          sql.NullString
	MetaData        sql.NullString
	Time            sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Timestamp       sql.NullInt64
}

func (OrderFromExternal) TableName() string {
	return "order_from_external"
}

// OrderFromExternalRepository manages order records that were fetched from the
// external exchange (e.g. Binance).
type OrderFromExternalRepository interface {
	// GetLastOrderFromExternalByPairID retrieves the most recent external order for a pair.
	GetLastOrderFromExternalByPairID(pairID int64, orderFromExternal *OrderFromExternal) error
	// Create inserts a new order record fetched from the external exchange.
	Create(orderFromExternal *OrderFromExternal) error
	// GetOrderByExternalOrderID looks up an order by its external exchange order ID.
	GetOrderByExternalOrderID(externalOrderID int64, orderFromExternal *OrderFromExternal) error
}

type TradeFromExternal struct {
	ID              int64
	OrderID         sql.NullInt64
	ExternalTradeID int64
	Price           sql.NullString
	Amount          sql.NullString
	Commission      sql.NullString
	Coin            sql.NullString `gorm:"column:commission_coin"`
	MetaData        sql.NullString
	time            sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Timestamp       sql.NullInt64
}

func (TradeFromExternal) TableName() string {
	return "trade_from_external"
}

// TradeFromExternalRepository manages trade records that were fetched from the
// external exchange (e.g. Binance).
type TradeFromExternalRepository interface {
	// GetLastTradeFromExternalByPairID retrieves the most recent external trade for a pair.
	GetLastTradeFromExternalByPairID(pairID int64, tradeFromExternal *TradeFromExternal) error
	// Create inserts a new trade record fetched from the external exchange.
	Create(tradeFromExternal *TradeFromExternal) error
}
