package processor

import (
	"context"
	"encoding/json"
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/livedata"
	"exchange-go/internal/order"
	"exchange-go/internal/orderbook"
	"exchange-go/internal/platform"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const RedisChannel = "channel:ticker"

var mutex = &sync.Mutex{}

// Processor handles incoming market data from external exchange WebSocket streams,
// transforming and storing trade, depth, kline, and ticker updates.
type Processor interface {
	// ProcessTrade processes a single trade event and updates the live trade book.
	ProcessTrade(ctx context.Context, trade Trade)
	// ProcessDepth processes a depth (order book) update from the external exchange.
	ProcessDepth(ctx context.Context, externalExchangeName string, pairName string, externalExchangePairName string, data []byte)
	// ProcessKline processes a kline (candlestick) update and stores it in Redis.
	ProcessKline(ctx context.Context, kline Kline)
	// ProcessTicker processes a ticker price update and publishes it via Redis pub/sub.
	ProcessTicker(ctx context.Context, ticker Ticker)
}

type processor struct {
	redisClient                platform.RedisClient
	liveDataService            livedata.Service
	priceGenerator             currency.PriceGenerator
	klineService               currency.KlineService
	orderBookService           orderbook.Service
	mqttManager                communication.MqttManager
	stopOrderSubmissionManager order.StopOrderSubmissionManager
	inQueueOrderManager        order.InQueueOrderManager
	queueMnager                communication.QueueManager
	logger                     platform.Logger
	currencyService            currency.Service
	pairCurrentTradeCounts     map[string]uint32
}

type Trade struct {
	Pair      string `json:"pair"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"createdAt"`
	IsMaker   bool   `json:"isMaker"`
	Ignore    bool   `json:"ignore"`
}

type Depth struct {
	Pair           string
	FirstUpdatedID int64
	FinalUpdatedID int64
	Bids           [][]string
	Asks           [][]string
}

type Ticker struct {
	Pair            string `json:"name"`
	Price           string `json:"price"`
	Percentage      string `json:"percentage"`
	ID              int64  `json:"id"`
	EquivalentPrice string `json:"equivalentPrice"`
	Volume          string `json:"volume"`
	High            string `json:"high"`
	Low             string `json:"low"`
}

type Kline struct {
	Pair                string `json:"pair"`
	TimeFrame           string `json:"timeFrame"`
	KlineStartTime      string `json:"startTime"`
	KlineCloseTime      string `json:"closeTime"`
	OpenPrice           string `json:"openPrice"`
	ClosePrice          string `json:"closePrice"`
	HighPrice           string `json:"highPrice"`
	LowPrice            string `json:"lowPrice"`
	BaseVolume          string `json:"baseVolume"`
	QuoteVolume         string `json:"quoteVolume"`
	TakerBuyBaseVolume  string `json:"takerBuyBaseVolume"`
	TakerBuyQuoteVolume string `json:"takerBuyQuoteVolume"`
	Spread              string `json:"spread"`
	IsOld               bool   `json:"isOld"`
}

func (p *processor) ProcessTrade(ctx context.Context, trade Trade) {
	payload, err := json.Marshal(trade)
	if err != nil {
		p.logger.Warn("can not marshal trade",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessTrade"),
		)
		return
	}

	redisTrade := livedata.RedisTrade{
		Price:     trade.Price,
		Amount:    trade.Amount,
		CreatedAt: trade.CreatedAt,
		IsMaker:   trade.IsMaker,
		Ignore:    trade.Ignore,
	}
	err = p.liveDataService.UpdateTradeBook(ctx, trade.Pair, redisTrade)
	if err != nil {
		p.logger.Warn("can not update tradebook",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessTrade"),
			zap.String("pairName", trade.Pair),
		)
		return
	}

	//since the rate of trades published by binance is too much we only publish 1 out of 10 trade
	mutex.Lock()
	count, exists := p.pairCurrentTradeCounts[trade.Pair]
	if exists {
		if count%10 == 0 {
			go p.mqttManager.PublishTrades(ctx, trade.Pair, payload)
			p.pairCurrentTradeCounts[trade.Pair] = 0
		}
		p.pairCurrentTradeCounts[trade.Pair]++
	} else {
		go p.mqttManager.PublishTrades(ctx, trade.Pair, payload)
		p.pairCurrentTradeCounts[trade.Pair] = 1
	}
	mutex.Unlock()

}

func (p *processor) ProcessDepth(ctx context.Context, externalExchangeName string, pairName string, externalExchangePairName string, data []byte) {
	precision := p.getPairPrecision(pairName)
	//we should create order book and then publish to mqtt
	orderBook, err := p.orderBookService.UpdateExternalOrderBook(ctx, externalExchangeName, pairName, externalExchangePairName, precision, data)
	if err != nil {
		p.logger.Warn("can not update external orderbook",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessDepth"),
			zap.String("pairName", pairName),
			zap.String("data", string(data)),
		)
		return
	}

	if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
		return
	}
	payload, err := json.Marshal(orderBook)
	if err != nil {
		p.logger.Warn("can not marshal orderbook",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessDepth"),
			zap.String("pairName", pairName),
		)
		return
	}

	go p.mqttManager.PublishOrderBook(ctx, pairName, payload)

}

func (p *processor) ProcessTicker(ctx context.Context, ticker Ticker) {
	formerPrice := p.getFormerPriceOfPair(ctx, ticker.Pair)
	p.setTickerDataInRedis(ctx, ticker)

	ticker.ID = p.getPairIDFromPairName(ticker.Pair)
	ticker.EquivalentPrice = p.getEquivalentPriceForPair(ctx, ticker.Pair, ticker.Price)

	payload, err := json.Marshal(ticker)
	if err != nil {
		p.logger.Warn("can not marshal ticker",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessTicker"),
			zap.String("pairName", ticker.Pair),
		)
		return
	}

	//p.publishTickerToRedisChannel(ctx, string(payload))
	go p.mqttManager.PublishTicker(ctx, ticker.Pair, payload)
	if ticker.Price != formerPrice {
		go p.stopOrderSubmissionManager.Submit(ctx, ticker.Pair, ticker.Price, formerPrice)
		go p.inQueueOrderManager.HandleInQueueOrders(ctx, ticker.Pair, ticker.Price)
	}
}

func (p *processor) ProcessKline(ctx context.Context, kline Kline) {
	payload, err := json.Marshal(kline)
	if err != nil {
		p.logger.Warn("can not marshal kline",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "ProcessKline"),
			zap.String("pairName", kline.Pair),
		)
		return
	}

	p.queueMnager.PublishKline(payload)
	go p.mqttManager.PublishKline(ctx, kline.Pair, kline.TimeFrame, payload)

}

func (p *processor) getPairIDFromPairName(pairName string) int64 {
	pair, err := p.currencyService.GetPairByName(pairName)
	if err != nil {
		p.logger.Warn("can not get pair by name",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "getPairPrecision"),
			zap.String("pairName", pairName),
		)
	}
	return pair.ID
}

func (p *processor) publishTickerToRedisChannel(ctx context.Context, payload string) {
	err := p.redisClient.Publish(ctx, RedisChannel, string(payload))
	if err != nil {
		p.logger.Warn("can not publish ticker to redis",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "publishTickerToRedisChannel"),
			zap.String("payload", string(payload)),
		)
	}

}

func (p *processor) setTickerDataInRedis(ctx context.Context, ticker Ticker) {
	err := p.liveDataService.SetPriceData(ctx, ticker.Pair, ticker.Price, ticker.Percentage, ticker.Volume)
	if err != nil {
		p.logger.Warn("can not set price data in redis",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "setTickerDataInRedis"),
			zap.String("pairName", ticker.Pair),
		)
	}
}

func (p *processor) getEquivalentPriceForPair(ctx context.Context, pairName string, price string) string {
	equivalentPrice, err := p.priceGenerator.GetPairPriceBasedOnUSDT(ctx, pairName)
	if err != nil {
		p.logger.Warn("can not get pair price based on usdt",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "getEquivalentPriceForPair"),
			zap.String("pairName", pairName),
		)
	}
	return equivalentPrice
}

func (p *processor) getPairPrecision(pairName string) int {
	pair, err := p.currencyService.GetPairByName(pairName)
	if err != nil {
		p.logger.Warn("can not get pair by name",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "getPairPrecision"),
			zap.String("pairName", pairName),
		)
		return 2
	}
	return pair.ShowDigits
}

func (p *processor) getFormerPriceOfPair(ctx context.Context, pairName string) string {
	price, err := p.liveDataService.GetPrice(ctx, pairName)
	if err != nil && err != redis.Nil {
		p.logger.Warn("can not get pair price",
			zap.Error(err),
			zap.String("service", "processor"),
			zap.String("method", "getFormerPriceOfPair"),
			zap.String("pairName", pairName),
		)
	}
	return price
}

func NewProcessor(redisClient platform.RedisClient, liveDataService livedata.Service, priceGenerator currency.PriceGenerator, klineService currency.KlineService, orderBookService orderbook.Service, mqttManager communication.MqttManager, stopOrderSubmissionManager order.StopOrderSubmissionManager, inQueueOrderManager order.InQueueOrderManager, queueMnager communication.QueueManager, logger platform.Logger, currencyService currency.Service) Processor {
	pairCurrentTradeCounts := make(map[string]uint32)
	return &processor{
		redisClient:                redisClient,
		liveDataService:            liveDataService,
		priceGenerator:             priceGenerator,
		klineService:               klineService,
		orderBookService:           orderBookService,
		mqttManager:                mqttManager,
		stopOrderSubmissionManager: stopOrderSubmissionManager,
		inQueueOrderManager:        inQueueOrderManager,
		queueMnager:                queueMnager,
		logger:                     logger,
		currencyService:            currencyService,
		pairCurrentTradeCounts:     pairCurrentTradeCounts,
	}
}
