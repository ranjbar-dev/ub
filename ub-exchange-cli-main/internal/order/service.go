package order

import (
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	SideAsk            = "ask"
	SideBid            = "bid"
	TypeSell           = "SELL"
	TypeBuy            = "BUY"
	ExchangeTypeLimit  = "LIMIT"
	ExchangeTypeMarket = "MARKET"
	ExchangeTypeFast   = "FAST_EXCHANGE" //not saved in database only showing to user

	MainTypeStopOrder = "stopOrder"
	MainTypeOrder     = "order"

	StatusOpen       = "OPEN"
	StatusFilled     = "FILLED"
	StatusCanceled   = "CANCELED"
	StatusExpired    = "EXPIRED"
	StatusProcessing = "PROCESSING"

	PlaceExternalExchange = "external_exchange"
	PlaceOurExchange      = "our_exchange"

	HistoryPeriod1Day   = "1DAY"
	HistoryPeriod1Week  = "1WEEK"
	HistoryPeriod1Month = "1MONTH"
	HistoryPeriod3Month = "3MONTH"
)

//this is done just for unit test (TestService_CancelOrder_ActionNotAllowed) is able to be use mock
var IsActionAllowed = currency.IsActionAllowed

type UserAgentInfo struct {
	Device  string `json:"device"`
	IP      string `json:"ip"`
	Browser string `json:"browser"`
}

type CreateOrderParams struct {
	Type           string `json:"type" binding:"required,oneof='buy' 'sell'"`
	ExchangeType   string `json:"exchange_type" binding:"required,oneof='market' 'limit'"`
	Amount         string `json:"amount" binding:"required"`
	PairID         int64  `json:"pair_currency_id" binding:"required,gt=0"`
	Price          string `json:"price"`
	IsInstant      bool
	IsFastExchange bool   `json:"is_fast_exchange"`
	StopPointPrice string `json:"stop_point_price"`
	UserAgentInfo  UserAgentInfo
}

type CancelOrderParams struct {
	ID int64 `json:"order_id" binding:"required,gt=0"`
}

type GetOpenOrdersParams struct {
	PairID int64 `json:"pair_currency_id"`
}

type GetOrdersHistoryParams struct {
	PairID         int64  `form:"pair_currency_id"`
	PairName       string `form:"pair_currency_name"`
	LastID         int64  `form:"last_id"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	Type           string `form:"type"`
	Hide           bool   `form:"hide"`
	Period         string `form:"period"`
	IsFastExchange *bool  `form:"is_fast_exchange"`
}

type GetOrderDetailParams struct {
	ID int64 `form:"order_id"`
}
type GetTradesHistoryParams struct {
	PairID    int64  `form:"pair_currency_id"`
	PairName  string `form:"pair_currency_name"`
	LastID    int64  `form:"last_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Type      string `form:"type"`
	Period    string `form:"period"`
}

type CreateOrderResponse struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"createdAt"`
	Price     string `json:"price"`
}

type OpenOrdersResponse struct {
	MainType         string          `json:"mainType"`
	OrderType        string          `json:"type"`
	Pair             string          `json:"pair"`
	ID               int64           `json:"id"`
	Side             string          `json:"side"`
	Price            string          `json:"price"`
	SubUnit          int             `json:"subUnit"`
	Amount           string          `json:"amount"`
	Total            string          `json:"total"`
	total            decimal.Decimal //should not be exported
	Executed         string          `json:"executed"`
	CreatedAt        string          `json:"createdAt"`
	createdAt        time.Time       //should not be exported
	TriggerCondition string          `json:"triggerCondition"`
}

type HistoryNeededField struct {
	OrderID   int64
	Path      string
	CreatedAt string
}

type HistoryFilters struct {
	UserID          int
	PairID          int64
	PairName        string
	DependentCoinID int64
	BasisCoinID     int64
	LastID          int64
	StartDate       string
	EndDate         string
	Type            string
	Hide            bool
	IsFullHistory   bool
	IsFastExchange  *bool
	PageSize        int
}

type TradeHistoryFilters struct {
	UserID          int
	PairID          int64
	PairName        string
	DependentCoinID int64
	BasisCoinID     int64
	LastID          int64
	StartDate       string
	EndDate         string
	Type            string
	IsFullHistory   bool
	PageSize        int
}

type HistoryResponse struct {
	MainType         string `json:"mainType"`
	OrderType        string `json:"type"`
	Pair             string `json:"pair"`
	ID               int64  `json:"id"`
	Side             string `json:"side"`
	Price            string `json:"price"`
	AveragePrice     string `json:"averagePrice"`
	SubUnit          int    `json:"subUnit"`
	Amount           string `json:"amount"`
	Total            string `json:"total"`
	Executed         string `json:"executed"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	Status           string `json:"status"`
	TriggerCondition string `json:"triggerCondition"`
}

type TradeHistoryResponse struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"createdAt"`
	Pair      string `json:"pair"`
	OrderType string `json:"type"`
	Price     string `json:"price"`
	SubUnit   int    `json:"subUnit"`
	Executed  string `json:"executed"`
	Fee       string `json:"fee"`
	Amount    string `json:"amount"`
	Total     string `json:"total"`
}

type FulfillOrderParams struct {
	ID int64 `json:"id" binding:"required"`
}

type DetailResponse struct {
	CreatedAt string `json:"createdAt"`
	Pair      string `json:"pair"`
	Type      string `json:"type"`
	SubUnit   int    `json:"subUnit"`
	Price     string `json:"price"`
	Executed  string `json:"executed"`
	Fee       string `json:"fee"`
	Amount    string `json:"amount"`
}

// Service is the main order service used by API handlers to create, cancel,
// and query orders and trades.
type Service interface {
	// CreateOrder validates input, creates a new order, and triggers post-creation events.
	CreateOrder(u *user.User, params CreateOrderParams) (apiResponse response.APIResponse, statusCode int)
	// CancelOrder cancels an existing open order and releases the frozen balance.
	CancelOrder(u *user.User, params CancelOrderParams) (apiResponse response.APIResponse, statusCode int)
	// GetOpenOrders returns the authenticated user's currently open orders, optionally filtered by pair.
	GetOpenOrders(u *user.User, params GetOpenOrdersParams) (apiResponse response.APIResponse, statusCode int)
	// GetOrdersHistory returns the user's historical orders. When isFullHistory is true, all
	// orders are returned; otherwise results are filtered by the provided parameters.
	GetOrdersHistory(u *user.User, params GetOrdersHistoryParams, isFullHistory bool) (apiResponse response.APIResponse, statusCode int)
	// GetTradesHistory returns the user's trade history. When isFullHistory is true, all
	// trades are returned; otherwise results are filtered by the provided parameters.
	GetTradesHistory(u *user.User, params GetTradesHistoryParams, isFullHistory bool) (apiResponse response.APIResponse, statusCode int)
	// GetOrderDetail returns detailed information for a single order belonging to the user.
	GetOrderDetail(u *user.User, params GetOrderDetailParams) (apiResponse response.APIResponse, statusCode int)

	// FulfillOrder is an admin-only endpoint that manually fulfills an order.
	FulfillOrder(adminUser *user.User, params FulfillOrderParams) (apiResponse response.APIResponse, statusCode int)
}

type service struct {
	db                    *gorm.DB
	orderRepository       Repository
	orderCreateManager    CreateManager
	eventsHandler         EventsHandler
	currencyService       currency.Service
	priceGenerator        currency.PriceGenerator
	userBalanceService    userbalance.Service
	orderRedisManager     RedisManager
	userConfigService     user.ConfigService
	userPermissionManager user.PermissionManager
	adminOrderManager     AdminOrderManager
	engineCommunicator    EngineCommunicator
	configs               platform.Configs
	logger                platform.Logger
}

func NewOrderService(db *gorm.DB, repo Repository, orderCreateManager CreateManager, eh EventsHandler, cs currency.Service,
	pg currency.PriceGenerator, ubs userbalance.Service, orm RedisManager, ucs user.ConfigService, upm user.PermissionManager,
	adminOrderManager AdminOrderManager, engineCommunicator EngineCommunicator, configs platform.Configs,
	logger platform.Logger) Service {
	return &service{
		db:                    db,
		orderRepository:       repo,
		orderCreateManager:    orderCreateManager,
		eventsHandler:         eh,
		currencyService:       cs,
		priceGenerator:        pg,
		userBalanceService:    ubs,
		orderRedisManager:     orm,
		userConfigService:     ucs,
		userPermissionManager: upm,
		adminOrderManager:     adminOrderManager,
		engineCommunicator:    engineCommunicator,
		configs:               configs,
		logger:                logger,
	}
}
