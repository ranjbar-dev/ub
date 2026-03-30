package order

import (
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"sync"

	"gorm.io/gorm"
)

type MatchingResult struct {
	Err                   error
	RemainingPartialOrder *CallBackOrderData
	RemovingDoneOrderIds  []int64
}

type CallBackOrderData struct {
	ID                   int64
	PairName             string
	OrderType            string
	Quantity             string
	Price                string
	Timestamp            int64
	TradedWithOrderID    int64
	QuantityTraded       string
	TradePrice           string
	MarketPrice          string
	IsAlreadyInOrderBook bool
	MinThresholdPrice    string
	MaxThresholdPrice    string
}

// PostOrderMatchingService performs post-trade settlement after the matching engine
// matches orders, including balance updates, trade record creation, and Centrifugo event publishing.
type PostOrderMatchingService interface {
	// HandlePostOrderMatching processes matched orders and an optional partial fill by updating
	// balances, creating trade records, and publishing events. When isFromAdmin is true the
	// matching was triggered by an admin fulfillment rather than the engine.
	HandlePostOrderMatching(ordersData []CallBackOrderData, partial *CallBackOrderData, isFromAdmin bool) MatchingResult
	// HandleExternalTradedOrder processes an order that was fulfilled on an external exchange,
	// updating the local order state and recording the external trade details.
	HandleExternalTradedOrder(data ExternalTradedOrderData) error
}

type postOrderMatchingService struct {
	db                 *gorm.DB
	orderRepository    Repository
	userBalanceService userbalance.Service
	forceTrader        ForceTrader
	priceGenerator     currency.PriceGenerator
	tradeEventsHandler TradeEventsHandler
	mqttManager        communication.CentrifugoManager
	rc                 platform.RedisClient
	currencyService    currency.Service
	userService        user.Service
	userLevelService   user.LevelService
	configs            platform.Configs
	logger             platform.Logger
	mu                 sync.Mutex
	currentMarketPrice string
	tempTrades         []tempTrade
	tradesData         []TradeData
	pushData           []orderPushPayload
}

type MatchingNeededQueryFields struct {
	OrderID            int64
	Price              string
	DemandedAmount     string
	OrderType          string
	OrderExchangeType  string
	PayedByAmount      string
	Path               string
	Status             string
	MarketPrice        string
	CreatedAt          string
	UserAgentInfo      string
	UserID             int
	UserEmail          string
	UserPrivateChannel string
	UserLevelID        int64
	MakerFeePercentage float64
	TakerFeePercentage float64
}

type tempTrade struct {
	price       string
	amount      string
	buyOrderID  int64
	sellOrderID int64
	pair        currency.Pair
}

type tempOrder struct {
	tradePrice        string
	tradeAmount       string
	orderType         string
	marketPrice       string
	isMaker           bool
	isPartial         bool
	TradedWithOrderID int64
}

type orderGroup struct {
	tempOrders   []tempOrder
	orderItem    MatchingNeededQueryFields
	userBalances [2]*userbalance.UserBalance
}

type partialOrderHandlingResult struct {
	isTraded bool
	err      error
	order    Order
}



type ExternalTradedOrderData struct {
	OrderID                 int64
	ExtraInfoID             int64
	Data                    string
	ExternalExchangeID      int64
	ExternalExchangeOrderID string
	Pair                    currency.Pair
}

func NewPostOrderMatchingService(db *gorm.DB, or Repository, ubs userbalance.Service, ft ForceTrader, pg currency.PriceGenerator, teh TradeEventsHandler,
	mqttManager communication.CentrifugoManager, rc platform.RedisClient, currencyService currency.Service, userService user.Service,
	userLevelService user.LevelService, configs platform.Configs, logger platform.Logger) PostOrderMatchingService {
	return &postOrderMatchingService{
		db:                 db,
		orderRepository:    or,
		userBalanceService: ubs,
		forceTrader:        ft,
		priceGenerator:     pg,
		tradeEventsHandler: teh,
		mqttManager:        mqttManager,
		rc:                 rc,
		currencyService:    currencyService,
		userService:        userService,
		userLevelService:   userLevelService,
		configs:            configs,
		logger:             logger,
	}
}
