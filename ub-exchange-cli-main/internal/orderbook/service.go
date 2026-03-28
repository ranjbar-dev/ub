package orderbook

import (
	"context"
	"encoding/json"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/livedata"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"math"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	Binance = "binance"
	BidType = "bid"
	AskType = "ask"
)

type GetOrderBookParams struct {
	Pair string `form:"pair" binding:"required"`
}

type GetTradeBookParams struct {
	Pair string `form:"pair" binding:"required"`
}

type GetTradeBookResponse struct {
}

// Service provides order book and trade book operations, including syncing
// order book data from external exchanges and serving API responses.
type Service interface {
	// UpdateExternalOrderBook processes raw order book data from an external exchange,
	// applies precision rounding, and returns the formatted order book.
	UpdateExternalOrderBook(ctx context.Context, externalExchangeName string, pairName string, externalExchangePairName string, precision int, data []byte) (OrderBook, error)
	// GetOrderBook returns the current order book for a trading pair as an API response.
	GetOrderBook(params GetOrderBookParams) (apiResponse response.APIResponse, statusCode int)
	// GetTradeBook returns the recent trade history for a trading pair as an API response.
	GetTradeBook(params GetTradeBookParams) (apiResponse response.APIResponse, statusCode int)
}

type RawOrderBook struct {
	Bids [][3]string
	Asks [][3]string
}
type OrderBook struct {
	Bids []BookItem `json:"bids"`
	Asks []BookItem `json:"asks"`
}

type BookItem struct {
	Price      string `json:"price"`
	Amount     string `json:"amount"`
	Value      string `json:"value"`
	Percentage string `json:"percentage"`
	Sum        string `json:"sum"`
	Type       string `json:"type"`

	price      float64 //these are for easy data manipulation and are not exported
	amount     float64
	value      float64
	percentage float64
	sum        float64
}

type externalExchangeOrderBook interface {
	updateOrderBook(ctx context.Context, pairName string, externalExchangePairName string, depth binanceDepth) (RawOrderBook, error)
}

type service struct {
	httpClient                platform.HTTPClient
	liveDataService           livedata.Service
	currencyService           currency.Service
	logger                    platform.Logger
	externalExchangeOrderBook externalExchangeOrderBook
}

func (s *service) UpdateExternalOrderBook(ctx context.Context, externalExchangeName string, pairName string, externalExchangePairName string, precision int, data []byte) (OrderBook, error) {
	orderBook := OrderBook{}
	switch externalExchangeName {
	case Binance:
		depth := binanceDepth{}

		err := json.Unmarshal(data, &depth)
		if err != nil {
			return orderBook, err
		}

		if s.externalExchangeOrderBook == nil {
			s.externalExchangeOrderBook = getBinanceOrderBook(s.httpClient, s.liveDataService)
		}

		rawOrderBook, err := s.externalExchangeOrderBook.updateOrderBook(ctx, pairName, externalExchangePairName, depth)
		if err != nil {
			return orderBook, err
		}
		return EnhanceOrderBookData(rawOrderBook, precision)

	default:
		break
	}

	return orderBook, nil
}

func (s *service) GetOrderBook(params GetOrderBookParams) (apiResponse response.APIResponse, statusCode int) {
	orderbook := OrderBook{
		Bids: make([]BookItem, 0),
		Asks: make([]BookItem, 0),
	}
	pair, err := s.currencyService.GetPairByName(params.Pair)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting pair", err,
			zap.String("service", "orderbook"),
			zap.String("method", "GetOrderBook"),
			zap.String("pairName", params.Pair),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Error("pair not found", http.StatusUnprocessableEntity, nil)
	}

	redisDepthSnapshot, err := s.liveDataService.GetDepthSnapshot(context.Background(), pair.Name)
	if err != nil && err != redis.Nil {
		s.logger.Error2("error getting redisDepthSnapshot", err,
			zap.String("service", "orderbook"),
			zap.String("method", "GetOrderBook"),
			zap.String("pair", params.Pair),
		)
		return response.Success(orderbook, "")
	}
	rawOrderBook := RawOrderBook{
		Bids: redisDepthSnapshot.Bids,
		Asks: redisDepthSnapshot.Asks,
	}
	orderbookResponse, err := EnhanceOrderBookData(rawOrderBook, 8) //second parameter is not important
	if err != nil {
		s.logger.Error2("error enhancing orderbook", err,
			zap.String("service", "orderbook"),
			zap.String("method", "GetOrderBook"),
			zap.String("pair", params.Pair),
		)
		return response.Success(orderbook, "")
	}
	if len(orderbookResponse.Asks) > 0 {
		orderbook.Asks = orderbookResponse.Asks
	}

	if len(orderbookResponse.Bids) > 0 {
		orderbook.Bids = orderbookResponse.Bids
	}
	return response.Success(orderbook, "")
}

func (s *service) GetTradeBook(params GetTradeBookParams) (apiResponse response.APIResponse, statusCode int) {
	tradebook := make([]livedata.RedisTrade, 0)
	pair, err := s.currencyService.GetPairByName(params.Pair)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("error getting pair", err,
			zap.String("service", "orderbook"),
			zap.String("method", "GetOrderBook"),
			zap.String("pairName", params.Pair),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Error("pair not found", http.StatusUnprocessableEntity, nil)
	}

	tradebookResponse, err := s.liveDataService.GetTradeBook(context.Background(), pair.Name)
	if err != nil && err != redis.Nil {
		s.logger.Error2("error getting redisDepthSnapshot", err,
			zap.String("service", "orderbook"),
			zap.String("method", "GetOrderBook"),
			zap.String("pair", params.Pair),
		)
	}
	if len(tradebookResponse) > 0 {
		//just reversing the slice
		for i, j := 0, len(tradebookResponse)-1; i < j; i, j = i+1, j-1 {
			tradebookResponse[i], tradebookResponse[j] = tradebookResponse[j], tradebookResponse[i]
		}
		tradebook = tradebookResponse
	}
	return response.Success(tradebook, "")
}

func NewOrderBookService(httpClient platform.HTTPClient, liveDataService livedata.Service, currencyService currency.Service, logger platform.Logger) Service {
	return &service{
		httpClient:                httpClient,
		liveDataService:           liveDataService,
		currencyService:           currencyService,
		logger:                    logger,
		externalExchangeOrderBook: nil,
	}
}

func EnhanceOrderBookData(rawOrderbook RawOrderBook, precision int) (OrderBook, error) {
	ob := OrderBook{}
	rawBids := rawOrderbook.Bids
	rawAsks := rawOrderbook.Asks

	bids, err := getBids(rawBids, precision)
	if err != nil {
		return ob, err
	}

	asks, err := getAsks(rawAsks, precision)
	if err != nil {
		return ob, err
	}

	if len(bids) > 100 {
		ob.Bids = bids[len(bids)-100:]
	} else {
		ob.Bids = bids
	}

	if len(asks) > 100 {
		ob.Asks = asks[:100]
	} else {
		ob.Asks = asks
	}

	return ob, nil

}

func isInFloat64Slice(element float64, slice []float64) bool {
	for _, uniquePrice := range slice {
		if uniquePrice == element {
			return true
		}
	}
	return false
}

func getFloat64(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return roundToPrecision(f, 8)
}

func getBids(bids [][3]string, precision int) ([]BookItem, error) {
	var bidsData []BookItem
	sumOfBids := float64(0)
	var uniqueBidPrices []float64

	for _, val := range bids {
		p := getFloat64(val[0])
		a := getFloat64(val[1])
		v := a * p
		tempSumOfBids := sumOfBids + a
		sumOfBids = roundToPrecision(tempSumOfBids, 8)
		if isInFloat64Slice(p, uniqueBidPrices) {
			for _, bookItem := range bidsData {
				if bookItem.price == p {
					bookItem.amount += a
					bookItem.value += v
					break
				}
			}

		} else {
			uniqueBidPrices = append(uniqueBidPrices, p)
			bookItem := BookItem{
				price:  p,
				amount: a,
				value:  v,
			}
			bidsData = append(bidsData, bookItem)
		}
	}

	if sumOfBids == 0 {
		return bidsData, nil
	}

	percentage := float64(0)
	remainingSumOfBids := float64(0)

	for i, itemBook := range bidsData {
		if i == 0 {
			bidsData[i].percentage = percentage
			percentage += itemBook.amount / sumOfBids
			bidsData[i].sum = itemBook.amount
			tempRemainingSumOfBids := remainingSumOfBids + itemBook.amount
			remainingSumOfBids = roundToPrecision(tempRemainingSumOfBids, 8)

		} else {
			percentage += itemBook.amount / sumOfBids
			bidsData[i].percentage = percentage
			tempRemainingSumOfBids := remainingSumOfBids + itemBook.amount
			remainingSumOfBids = roundToPrecision(tempRemainingSumOfBids, 8)
			bidsData[i].sum = remainingSumOfBids

		}
		bidsData[i].Price = strconv.FormatFloat(bidsData[i].price, 'f', precision, 64)
		bidsData[i].Amount = strconv.FormatFloat(bidsData[i].amount, 'f', precision, 64)
		bidsData[i].Value = strconv.FormatFloat(bidsData[i].value, 'f', precision, 64)
		bidsData[i].Percentage = strconv.FormatFloat(bidsData[i].percentage, 'f', 2, 64)
		bidsData[i].Sum = strconv.FormatFloat(bidsData[i].sum, 'f', precision, 64)
		bidsData[i].Type = BidType
	}
	sort.Slice(bidsData, func(i, j int) bool {
		return bidsData[i].price < bidsData[j].price
	})
	return bidsData, nil

}

func getAsks(asks [][3]string, precision int) ([]BookItem, error) {
	var asksData []BookItem
	sumOfAsks := float64(0)
	var uniqueAskPrices []float64

	//sort.Slice(asks, func(i, j int) bool {
	//return asks[i][0] > asks[j][0]
	//})

	for _, val := range asks {
		p := getFloat64(val[0])
		a := getFloat64(val[1])
		//if p < 0.98 * bestBidPrice {
		//	continue
		//}
		v := a * p

		sumOfAsks += a
		if isInFloat64Slice(p, uniqueAskPrices) {
			for _, bookItem := range asksData {
				if bookItem.price == p {
					bookItem.amount += a
					bookItem.value += v
					break
				}
			}

		} else {
			uniqueAskPrices = append(uniqueAskPrices, p)
			bookItem := BookItem{
				price:  p,
				amount: a,
				value:  v,
			}
			asksData = append(asksData, bookItem)
		}
	}

	if sumOfAsks == 0 {
		return asksData, nil
	}

	percentage := float64(0)
	remainingSumOfAsks := float64(0)
	for i, itemBook := range asksData {
		if i == 0 {
			asksData[i].percentage = percentage
			percentage += itemBook.amount / sumOfAsks
			asksData[i].sum = itemBook.amount
			tempRemainingSumOfAsks := remainingSumOfAsks + itemBook.amount
			remainingSumOfAsks = roundToPrecision(tempRemainingSumOfAsks, 8)

		} else {
			percentage += itemBook.amount / sumOfAsks
			asksData[i].percentage = percentage
			tempRemainingSumOfAsks := remainingSumOfAsks + itemBook.amount
			remainingSumOfAsks = roundToPrecision(tempRemainingSumOfAsks, 8)
			asksData[i].sum = remainingSumOfAsks

		}
		asksData[i].Price = strconv.FormatFloat(asksData[i].price, 'f', precision, 64)
		asksData[i].Amount = strconv.FormatFloat(asksData[i].amount, 'f', precision, 64)
		asksData[i].Value = strconv.FormatFloat(asksData[i].value, 'f', precision, 64)
		asksData[i].Percentage = strconv.FormatFloat(asksData[i].percentage, 'f', 2, 64)
		asksData[i].Sum = strconv.FormatFloat(asksData[i].sum, 'f', precision, 64)
		asksData[i].Type = AskType
	}

	//sort.Slice(asksData, func(i, j int) bool {
	//	return asksData[i].price < asksData[j].price
	//})

	return asksData, nil

}

func roundToPrecision(number float64, precision int) float64 {
	return math.Round(number*math.Pow10(precision)) / math.Pow10(precision)
}
