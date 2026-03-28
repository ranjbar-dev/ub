package order

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/user"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TradeFilters struct {
	UserID    int64
	PairID    int64
	StartDate string
	EndDate   string
}

type Order struct {
	ID                            int64
	UserID                        int       `gorm:"column:creator_user_id"`
	User                          user.User `gorm:"foreignKey:UserID"`
	ParentID                      sql.NullInt64
	Type                          string
	ExchangeType                  string
	Price                         sql.NullString
	Status                        string
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
	DemandedAmount                sql.NullString `gorm:"column:demanded_money_amount"`
	DemandedCoin                  string         `gorm:"column:demanded_money_currency"`
	PayedByAmount                 sql.NullString `gorm:"column:payed_by_money_amount"`
	PayedByCoin                   string         `gorm:"column:payed_by_money_currency"`
	PairID                        int64          `gorm:"column:pair_currency_id"`
	Pair                          currency.Pair  `gorm:"foreignKey:PairID"`
	ExtraInfoID                   sql.NullInt64
	FinalDemanded                 sql.NullString `gorm:"column:final_demanded_money"`
	TradePrice                    sql.NullString
	IsMaker                       sql.NullBool
	FeePercentage                 sql.NullFloat64
	ExternalExchangeFeePercentage sql.NullFloat64
	FinalPayedBy                  sql.NullString `gorm:"column:final_payed_by_money"`
	Level                         sql.NullInt64
	Path                          sql.NullString
	FinalDemandedAmount           sql.NullString
	FinalPayedByAmount            sql.NullString
	StopPointPrice                sql.NullString
	IsSubmitted                   sql.NullBool
	IsTradedWithBot               sql.NullBool
	CurrentMarketPrice            sql.NullString
	IsFastExchange                bool `gorm:"default:false"`
}

func (o Order) IsStopOrder() bool {
	if o.StopPointPrice.Valid {
		return true
	}

	return false
}

func (o Order) getAmount() string {
	if o.Type == TypeBuy {
		return o.DemandedAmount.String
	}
	return o.PayedByAmount.String
}

func (o Order) getStringID() string {
	return strconv.FormatInt(o.ID, 10)
}

func (o Order) isMarket() bool {
	if !o.Price.Valid {
		return true
	}
	return false
}

// Repository provides data access methods for order records.
type Repository interface {
	// GetOrdersDataByIdsWithJoinUsingTx retrieves matching-related order data with joined
	// associations for the given order IDs within the provided transaction.
	GetOrdersDataByIdsWithJoinUsingTx(tx *gorm.DB, orderIds []int64) []MatchingNeededQueryFields
	// GetOrdersByIds retrieves multiple orders by their IDs in a single batch query.
	GetOrdersByIds(orderIds []int64) []Order
	// GetOrderByID loads a single order by its ID into the provided Order pointer.
	GetOrderByID(id int64, o *Order) error
	// GetOrderByIDUsingTx loads a single order by its ID within the provided transaction.
	GetOrderByIDUsingTx(tx *gorm.DB, id int64, o *Order) error
	// GetUserOpenOrders returns all open orders for a given user and currency pair.
	GetUserOpenOrders(userID int, pairID int64) []Order
	// GetOrdersAncestors retrieves the parent orders for partial-fill ancestor IDs.
	GetOrdersAncestors(ancestorsIds []int64) []Order
	// GetLeafOrders returns terminal (leaf) orders matching the given history filters.
	GetLeafOrders(filters HistoryFilters) []HistoryNeededField
	// GetUserTradedOrders returns completed trade orders matching the given trade history filters.
	GetUserTradedOrders(filters TradeHistoryFilters) []Order
	// GetUserOrderDetailsByID returns the detailed order information for a specific order and user.
	GetUserOrderDetailsByID(id int64, userID int) []Order
	// GetOpenOrders returns all open orders created before the specified date.
	GetOpenOrders(date string) []Order
}

type ExtraInfo struct {
	ID                              int64
	UserAgentInfo                   sql.NullString
	ExternalExchangeOtherInfo       sql.NullString
	ExternalExchangeID              sql.NullInt64
	IsMarketOrderInExternalExchange sql.NullBool
	PayedByDiff                     sql.NullString
	ExternalExchangeOrderID         sql.NullString
	AutoExchange                    sql.NullBool
}

func (ExtraInfo) TableName() string {
	return "orders_extra_info"
}

type Trade struct {
	ID           int64
	Price        sql.NullString
	Amount       sql.NullString
	PairID       int64 `gorm:"column:pair_currency_id"`
	BuyOrderID   sql.NullInt64
	SellOrderID  sql.NullInt64
	BotOrderType sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TradeRepository provides data access methods for trade records.
type TradeRepository interface {
	// Create persists a new trade record to the database.
	Create(trade *Trade) error
	// GetTradesOfUserBetweenTimes returns trades for a user within the specified date range.
	GetTradesOfUserBetweenTimes(userID int, startTime, endTime string) []Trade
	// GetBotTradesByIDAndCreatedAtGreaterThan returns bot trades for a pair after the given trade ID and timestamp.
	GetBotTradesByIDAndCreatedAtGreaterThan(pairID int64, tradeID int64, createdAt time.Time) []Trade
}
