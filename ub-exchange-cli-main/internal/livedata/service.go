package livedata

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	Prefix                      = "live_data"
	Price                       = "price"
	Percentage                  = "change_price_percentage"
	Volume                      = "volume"
	Pair                        = "pair_currency"
	TradeBook                   = "trade_book"
	PreKline                    = "pre_kline"
	Kline                       = "kline"
	DepthSnapshot               = "depth_snapshot"
	OrderBook                   = "order_book"
	PriceLastInsertTime         = "price_last_insert_time"
	TradeBookLastInsertTime     = "trade_book_last_insert_time"
	DepthSnapshotLastInsertTime = "depth_snapshot_last_insert_time"
	KlineLastInsertTime         = "kline_last_insert_time"
	LastAggregationTime         = "last_aggregation_time"
)

type RedisTrade struct {
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"createdAt"`
	IsMaker   bool   `json:"isMaker"`
	Ignore    bool   `json:"ignore"`
}

type RedisKline struct {
	TimeFrame           string `json:"timeFrame"`
	KlineStartTime      string `json:"ohlcStartTime"`
	KlineCloseTime      string `json:"ohlcCloseTime"`
	OpenPrice           string `json:"openPrice"`
	ClosePrice          string `json:"closePrice"`
	HighPrice           string `json:"highPrice"`
	LowPrice            string `json:"lowPrice"`
	BaseVolume          string `json:"baseVolume"`
	QuoteVolume         string `json:"quoteVolume"`
	TakerBuyBaseVolume  string `json:"takerBuyBaseVolume"`
	TakerBuyQuoteVolume string `json:"takerBuyQuoteVolume"`
}

type RedisDepthSnapshot struct {
	LastUpdatedID int64       `json:"lastUpdatedId"`
	UpdatedAt     int64       `json:"updatedAt"`
	Bids          [][3]string `json:"bids"`
	Asks          [][3]string `json:"asks"`
}

type RedisOrderBook struct {
	Bids map[string]string `json:"bids"`
	Asks map[string]string `json:"asks"`
}

type RedisPairPriceData struct {
	PairName   string
	Price      string
	Percentage string
	Volume     string
}

// Service manages real-time market data in Redis, including prices, trade books,
// kline candles, depth snapshots, and order books for all trading pairs.
type Service interface {
	// GetPrice returns the current price for a trading pair from Redis.
	GetPrice(ctx context.Context, pairName string) (string, error)
	// SetPriceData stores price, percentage change, and volume for a trading pair in Redis.
	SetPriceData(ctx context.Context, pairName string, price string, percentage string, volume string) error
	// UpdateTradeBook appends a new trade to the pair's trade book in Redis, keeping the last 20 trades.
	UpdateTradeBook(ctx context.Context, pairName string, trade RedisTrade) error
	// GetTradeBook returns the recent trade history for a trading pair from Redis.
	GetTradeBook(ctx context.Context, pairName string) ([]RedisTrade, error)
	// SetKline stores a new kline candle in Redis, rotating the current kline to pre-kline.
	SetKline(ctx context.Context, pairName string, kline RedisKline) error
	// GetKline returns the current kline candle for a pair and time frame from Redis.
	GetKline(ctx context.Context, pairName string, timeFrame string) (RedisKline, error)
	// GetPreKline returns the previous kline candle for a pair and time frame from Redis.
	GetPreKline(ctx context.Context, pairName string, timeFrame string) (RedisKline, error)
	// SetDepthSnapshot stores the full depth (order book) snapshot for a pair in Redis.
	SetDepthSnapshot(ctx context.Context, pairName string, snapshot RedisDepthSnapshot) error
	// GetDepthSnapshot returns the depth (order book) snapshot for a pair from Redis.
	GetDepthSnapshot(ctx context.Context, pairName string) (RedisDepthSnapshot, error)
	// SetOrderBook stores the aggregated order book for a pair in Redis.
	SetOrderBook(ctx context.Context, pairName string, orderBook RedisOrderBook) error
	// GetOrderBook returns the aggregated order book for a pair from Redis.
	GetOrderBook(ctx context.Context, pairName string) (RedisOrderBook, error)
	// GetPairsPriceData returns price, percentage change, and volume for multiple pairs in a single pipeline.
	GetPairsPriceData(ctx context.Context, pairNames []string) ([]RedisPairPriceData, error)
	// GetLastInsertTime returns the Unix timestamp of the last data insert for the given field and pair.
	GetLastInsertTime(ctx context.Context, pairName string, fieldName string) (string, error)
	// GetLastAggregationTime returns the Unix timestamp of the last order aggregation for a pair.
	GetLastAggregationTime(ctx context.Context, pairName string) (string, error)
	// SetLastAggregationTime records the Unix timestamp of the latest order aggregation for a pair.
	SetLastAggregationTime(ctx context.Context, pairName string, timestamp int64) error
}

type service struct {
	redisClient platform.RedisClient
}

func (s *service) setLiveData(ctx context.Context, key string, values ...interface{}) error {
	return s.redisClient.HSet(ctx, key, values...)

}

func (s *service) getLiveData(ctx context.Context, key string, field string) (string, error) {
	return s.redisClient.HGet(ctx, key, field)
}

func (s *service) GetPrice(ctx context.Context, pairName string) (string, error) {
	key := getKeyForPair(pairName)
	return s.getLiveData(ctx, key, Price)
}

func (s *service) SetPriceData(ctx context.Context, pairName string, price string, percentage string, volume string) error {
	key := getKeyForPair(pairName)
	return s.setLiveData(ctx, key, Price, price, Percentage, percentage, Volume, volume, PriceLastInsertTime, time.Now().Unix())
}

func (s *service) UpdateTradeBook(ctx context.Context, pairName string, trade RedisTrade) error {
	key := getKeyForPair(pairName)
	formerTrades, err := s.GetTradeBook(ctx, pairName)
	if err != nil && err != redis.Nil {
		return err
	}

	allTrades := append(formerTrades, trade)
	if length := len(allTrades); length > 20 {
		allTrades = allTrades[length-20 : length-1]
	}

	finalData, err := json.Marshal(allTrades)
	if err != nil {
		return err
	}

	return s.setLiveData(ctx, key, TradeBook, finalData, TradeBookLastInsertTime, time.Now().Unix())
}

func (s *service) GetTradeBook(ctx context.Context, pairName string) ([]RedisTrade, error) {
	var trades []RedisTrade
	key := getKeyForPair(pairName)

	if s.redisClient.Exists(ctx, key) {
		tradesString, err := s.getLiveData(ctx, key, TradeBook)
		if err != nil {
			return trades, err
		}
		err = json.Unmarshal([]byte(tradesString), &trades)

		return trades, err
	}
	return trades, nil

}

func (s *service) SetKline(ctx context.Context, pairName string, kline RedisKline) error {
	//in this function we get current kline and set it as pre kline then set the new receiving kline as kline
	key := getKeyForPair(pairName)
	preKlineField := PreKline + "_" + kline.TimeFrame
	klineField := Kline + "_" + kline.TimeFrame

	data, err := json.Marshal(kline)
	if err != nil {
		return err
	}

	currentKline, err := s.GetKline(ctx, pairName, kline.TimeFrame)

	//this means key exists
	if err == nil {
		currentKlineData, err := json.Marshal(currentKline)
		if err != nil {
			return err
		}
		err = s.setLiveData(ctx, key, preKlineField, currentKlineData)
		if err != nil {
			return err
		}
		err = s.setLiveData(ctx, key, klineField, data, KlineLastInsertTime, time.Now().Unix())
		if err != nil {
			return err
		}
		return nil

	} else {
		if err == redis.Nil {
			//reaching here it means we have no current kline in redis so we set it
			err = s.setLiveData(ctx, key, klineField, data, KlineLastInsertTime, time.Now().Unix())
			if err != nil {
				return err
			}

			err = s.setLiveData(ctx, key, preKlineField, data)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return err

}

func (s *service) GetPreKline(ctx context.Context, pairName string, timeFrame string) (RedisKline, error) {
	redisKline := RedisKline{}
	key := getKeyForPair(pairName)
	field := PreKline + "_" + timeFrame
	res, err := s.getLiveData(ctx, key, field)
	if err != nil {
		return redisKline, err
	}
	if res == "" {
		return redisKline, nil
	}
	err = json.Unmarshal([]byte(res), &redisKline)
	if err != nil {
		return redisKline, err
	}
	return redisKline, nil

}

func (s *service) GetKline(ctx context.Context, pairName string, timeFrame string) (RedisKline, error) {
	redisKline := RedisKline{}
	key := getKeyForPair(pairName)
	field := Kline + "_" + timeFrame
	res, err := s.getLiveData(ctx, key, field)
	if err != nil {
		return redisKline, err
	}
	err = json.Unmarshal([]byte(res), &redisKline)
	if err != nil {
		return redisKline, err
	}
	return redisKline, nil

}

func (s *service) SetDepthSnapshot(ctx context.Context, pairName string, snapshot RedisDepthSnapshot) error {
	key := getKeyForPair(pairName)
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	err = s.setLiveData(ctx, key, DepthSnapshot, data, DepthSnapshotLastInsertTime, time.Now().Unix())
	return err

}

func (s *service) GetDepthSnapshot(ctx context.Context, pairName string) (RedisDepthSnapshot, error) {
	redisDepthSnapshot := RedisDepthSnapshot{}
	key := getKeyForPair(pairName)
	res, err := s.getLiveData(ctx, key, DepthSnapshot)
	if err != nil {
		return redisDepthSnapshot, err
	}

	err = json.Unmarshal([]byte(res), &redisDepthSnapshot)
	if err != nil {
		return redisDepthSnapshot, err
	}

	return redisDepthSnapshot, nil

}

func (s *service) SetOrderBook(ctx context.Context, pairName string, orderBook RedisOrderBook) error {
	key := getKeyForPair(pairName)
	data, err := json.Marshal(orderBook)
	if err != nil {
		return err
	}

	err = s.setLiveData(ctx, key, OrderBook, data)
	return err
}

func (s *service) GetOrderBook(ctx context.Context, pairName string) (RedisOrderBook, error) {
	orderBook := RedisOrderBook{}
	key := getKeyForPair(pairName)
	res, err := s.getLiveData(ctx, key, OrderBook)
	if err != nil {
		return orderBook, err
	}
	err = json.Unmarshal([]byte(res), &orderBook)
	return orderBook, err
}

func (s *service) GetPairsPriceData(ctx context.Context, pairNames []string) ([]RedisPairPriceData, error) {
	var result []RedisPairPriceData
	pipe := s.redisClient.TxPipeline()
	m := map[string]*redis.SliceCmd{}
	for _, pairName := range pairNames {
		key := getKeyForPair(pairName)
		m[pairName] = pipe.HMGet(ctx, key, Price, Percentage, Volume)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return result, err
	}
	for pairName, v := range m {
		res, _ := v.Result()

		price, _ := res[0].(string)
		percentage, _ := res[1].(string)
		volume, _ := res[2].(string)
		pairPriceData := RedisPairPriceData{
			PairName:   pairName,
			Price:      price,
			Percentage: percentage,
			Volume:     volume,
		}
		result = append(result, pairPriceData)
	}
	return result, nil

}

func (s *service) GetLastInsertTime(ctx context.Context, pairName string, fieldName string) (string, error) {
	key := getKeyForPair(pairName)
	return s.getLiveData(ctx, key, fieldName)
}

func (s *service) GetLastAggregationTime(ctx context.Context, pairName string) (string, error) {
	key := getKeyForPair(pairName)
	return s.getLiveData(ctx, key, LastAggregationTime)
}

func (s *service) SetLastAggregationTime(ctx context.Context, pairName string, timestamp int64) error {
	key := getKeyForPair(pairName)
	return s.setLiveData(ctx, key, LastAggregationTime, timestamp)
}

func getKeyForPair(pairName string) string {
	return Prefix + ":" + Pair + ":" + pairName
}

func NewLiveDataService(redisClient platform.RedisClient) Service {
	return &service{redisClient}
}
