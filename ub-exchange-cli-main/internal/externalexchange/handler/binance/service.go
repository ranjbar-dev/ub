package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"exchange-go/internal/externalexchange/handler"
	"exchange-go/internal/platform"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	APIHost             = "https://api.binance.com"
	NewOrderURI         = "/api/v3/order"
	WithdrawURI         = "/sapi/v1/capital/withdraw/apply"
	NewTestOrderURI     = "/api/v3/order/test"
	FetchKlineURI       = "/api/v3/klines?"
	FetchWithdrawalsURI = "/sapi/v1/capital/withdraw/history?"
	FetchOrdersURI      = "/api/v3/allOrders?"
	FetchTradesURI      = "/api/v3/myTrades?"
	ExchangeInfoURI     = "/api/v3/exchangeInfo"
	ExchangeTypeLimit   = "LIMIT"

	NewOrderReuqestWeight        = 1
	KlineReuqestWeight           = 1
	WithdrawHistoryReuqestWeight = 1
	OrdersReuqestWeight          = 10
	TradesReuqestWeight          = 10
	ExchangeInfoReuqestWeight    = 10
)

var binanceTxStatusToOurStatus = map[int]string{
	0: "CREATED",
	1: "CANCELED",
	2: "CREATED",
	3: "REJECTED",
	4: "IN_PROGRESS",
	5: "FAILED",
	6: "COMPLETED",
}

type metadata struct {
	APIKey string `json:"apiKey"`
	Secret string `json:"secret"`
}

type newOrderRequestBody struct {
	Symbol     string  `url:"symbol"`
	Side       string  `url:"side"`
	Type       string  `url:"type"`
	Price      float64 `url:"price,omitempty"`
	Quantity   float64 `url:"quantity"`
	RecvWindow int64   `url:"recvWindow"`
	Timestamp  int64   `url:"timestamp"`
	Signature  string  `url:"-"`
}

type fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}
type newOrderResponseBody struct {
	Symbol             string `json:"symbol"`
	OrderID            int64  `json:"orderId"`
	OrderListID        int64  `json:"orderListId"`
	ClientOrderID      string `json:"clientOrderId"`
	TransactTime       int64  `json:"transactTime"`
	Price              string `json:"price"`
	OrigQty            string `json:"origQty"`
	ExecutedQty        string `json:"executedQty"`
	CumulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status             string `json:"status"`
	TimeInForce        string `json:"timeInForce"`
	Type               string `json:"type"`
	Side               string `json:"side"`
	Fills              []fill `json:"fills"`
}

type fetchKlinesResponseBody [][]interface{}

type fetchWithdrawalsResponseBody struct {
	Address         string `json:"address"`
	Amount          string `json:"amount"`
	ApplyTime       string `json:"applyTime"`
	Coin            string `json:"coin"`
	ID              string `json:"id"`
	WithdrawOrderID string `json:"withdrawOrderId"`
	Network         string `json:"network"`
	TransferType    int    `json:"transferType"`
	Status          int    `json:"status"`
	TransactionFee  string `json:"transactionFee"`
	TxID            string `json:"txId"`
}

type fetchTradesResponseBody struct {
	Symbol          string `json:"symbol"`
	ID              int64  `json:"id"`
	OrderID         int64  `json:"orderId"`
	OrderListID     int64  `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

type fetchOrdersResponseBody struct {
	Symbol              string `json:"symbol"`
	OrderID             int64  `json:"orderId"`
	OrderListID         int64  `json:"orderListId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

type withdrawRequestBody struct {
	Coin      string `url:"coin"`
	Address   string `url:"address"`
	Amount    string `url:"amount"`
	Network   string `url:"network"`
	Timestamp int64  `url:"timestamp"`
}

type withdrawResponseBody struct {
	ID string `json:"id"`
}

type manipulatedOrderData struct {
	price       string
	amount      string
	isPlaceable bool
}

type service struct {
	httpClient       platform.HTTPClient
	configs          platform.Configs
	logger           platform.Logger
	exchangeInfo     *exchangeInfoService
	rateLimitHandler *rateLimitHandler
	metadata         metadata
}

func (s *service) NewOrder(params handler.NewOrderParams) (res handler.NewOrderResult, err error) {
	if !s.rateLimitHandler.canPlaceOrder() {
		err := handler.CustomError{
			Type:    handler.ErrTypeRateLimit,
			Message: "ratelimit reached",
		}
		res.CustomError = err
		return res, err
	}
	newData, err := s.tryToMakeOrderPlaceable(params)
	if err != nil {
		err := handler.CustomError{
			Type:    handler.ErrTypeInvalidParams,
			Message: err.Error(),
		}
		res.CustomError = err
		return res, err
	}

	if !newData.isPlaceable {
		err := handler.CustomError{
			Type:    handler.ErrTypeInvalidParams,
			Message: "not matching with binance filter",
		}
		res.CustomError = err
		return res, err
	}

	//if we are here the order is placeable
	//currently all of our orders are market
	//we do not care about price
	params.Amount = newData.amount
	ctx := context.Background()
	url := s.getNewOrderURL()
	body, err := s.getNewOrderRequestBody(params)
	if err != nil {
		return res, err
	}

	header := getHeader(s.metadata.APIKey, http.MethodPost)
	resp, respHeader, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)
	if err != nil {
		res.IsOrderPlaced = false
		return res, err
	}

	if statusCode != http.StatusOK {
		res.IsOrderPlaced = false
		s.logger.Warn("error in new order in binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
			zap.String("pair", params.Pair),
		)

		errorType, message := s.getCustomError(resp)
		res.CustomError = handler.CustomError{
			Type:    errorType,
			Message: message,
		}
		return res, fmt.Errorf("status code is not 200 it is %d", statusCode)
	}

	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	resBody := newOrderResponseBody{}
	err = json.Unmarshal(resp, &resBody)

	if err != nil {
		res.IsOrderPlaced = true
		return res, err
	}

	ExternalExchangeOrderIDString := strconv.FormatInt(resBody.OrderID, 10)
	res = handler.NewOrderResult{
		IsOrderPlaced:           true,
		ExternalExchangeOrderID: ExternalExchangeOrderIDString,
		Data:                    string(resp),
	}

	return res, nil
}

//getting new price and quantity so the order be matched with binance filters
//since all of our orders are market, we only check filters related to market orders
func (s *service) tryToMakeOrderPlaceable(params handler.NewOrderParams) (newData manipulatedOrderData, err error) {
	symbol := strings.Replace(params.Pair, "-", "", -1)
	amountDecimal, err := decimal.NewFromString(params.Amount)
	if err != nil {
		return newData, fmt.Errorf("tryToMakeOrderPlaceable: our amount is not valid: %w", err)
	}
	//checking marketLotSize
	marketMin, _, stepSize, err := s.exchangeInfo.getMarketLotSizeFilters(symbol)
	if err != nil {
		return newData, fmt.Errorf("tryToMakeOrderPlaceable: can not get market lot size: %w", err)
	}
	marketMinDecimal, err := decimal.NewFromString(marketMin)
	if err != nil {
		return newData, fmt.Errorf("tryToMakeOrderPlaceable: marketMin is not valid: %w", err)
	}
	stepSizeDecimal, err := decimal.NewFromString(stepSize)
	if err != nil {
		return newData, fmt.Errorf("tryToMakeOrderPlaceable: step size is not valid: %w", err)
	}
	newAmountDecimal := amountDecimal
	if !stepSizeDecimal.IsZero() {
		n := amountDecimal.Sub(marketMinDecimal).Div(stepSizeDecimal).Round(0)
		newAmountDecimal = marketMinDecimal.Add(stepSizeDecimal.Mul(n))
	}
	//checking the lot size here
	minNotional, applyToMarket, err := s.exchangeInfo.getMinNotional(symbol)
	if err != nil {
		return newData, fmt.Errorf("tryToMakeOrderPlaceable: can not get minNotional: %w", err)
	}
	if applyToMarket {
		minNotionalDecimal, err := decimal.NewFromString(minNotional)
		if err != nil {
			return newData, fmt.Errorf("tryToMakeOrderPlaceable: min notional is not valid: %w", err)
		}
		currentMarketPriceDecimal, err := decimal.NewFromString(params.CurrentMarketPrice)
		if err != nil {
			return newData, fmt.Errorf("tryToMakeOrderPlaceable: current market price is not valid: %w", err)
		}
		if currentMarketPriceDecimal.Mul(newAmountDecimal).LessThan(minNotionalDecimal) {
			return manipulatedOrderData{
				price:       params.Price,
				amount:      newAmountDecimal.String(),
				isPlaceable: false,
			}, nil
		}
	}
	//if the order is less than minimum then the order is not placeable
	if newAmountDecimal.LessThan(marketMinDecimal) {
		return manipulatedOrderData{
			price:       params.Price,
			amount:      params.Amount,
			isPlaceable: false,
		}, nil
	}
	return manipulatedOrderData{
		price:       params.Price,
		amount:      newAmountDecimal.String(),
		isPlaceable: true,
	}, nil

}

func (s *service) getCustomError(body []byte) (errorType, message string) {
	type errorBody struct {
		Code    int64  `json:"code"`
		Message string `json:"msg"`
	}
	resBody := errorBody{}

	err := json.Unmarshal(body, &resBody)
	if err != nil {
		return handler.ErrTypeUnknown, resBody.Message
	}

	switch resBody.Code {
	case int64(1013):
		return handler.ErrTypeInvalidParams, resBody.Message
	default:
		return handler.ErrTypeUnknown, resBody.Message
	}

}

func (s *service) getNewOrderURL() string {
	if s.configs.GetEnv() == platform.EnvTest {
		return APIHost + NewTestOrderURI
	}
	return APIHost + NewOrderURI
}

func getHeader(apiKey string, reqMethod string) map[string]string {
	header := make(map[string]string)
	if reqMethod != http.MethodGet {
		header["Content-Type"] = "application/x-www-form-urlencoded"
	}
	header["X-MBX-APIKEY"] = apiKey
	return header
}

func (s *service) getNewOrderRequestBody(params handler.NewOrderParams) (body []byte, err error) {
	symbol := strings.Replace(params.Pair, "-", "", -1)
	price := float64(0)
	amountDecimal, err := decimal.NewFromString(params.Amount)
	if err != nil {
		return body, err
	}

	if params.ExchangeType == ExchangeTypeLimit {
		price, _ = strconv.ParseFloat(params.Price, 64)
	}

	lotSize, err := s.exchangeInfo.getLotSizeForSymbol(symbol)
	if err != nil {
		return body, err
	}

	lotSizeFloat64, err := strconv.ParseFloat(lotSize, 64)
	if err != nil {
		return body, err
	}

	log := -1 * math.Log10(lotSizeFloat64)

	logDecimal := decimal.NewFromFloat(log)
	precision := logDecimal.Round(0).IntPart()

	amount, _ := amountDecimal.Truncate(int32(precision)).Float64()

	timestamp := s.getTimestamp()

	rb := newOrderRequestBody{
		Symbol:     symbol,
		Side:       params.Type,
		Type:       params.ExchangeType,
		Price:      price,
		Quantity:   amount,
		RecvWindow: 60000,
		Timestamp:  timestamp,
	}

	v, err := query.Values(rb)
	if err != nil {
		return body, err
	}

	hash := hmac.New(sha256.New, []byte(s.metadata.Secret))

	hash.Write([]byte(v.Encode()))
	signature := hex.EncodeToString(hash.Sum(nil))

	final := v.Encode() + "&signature=" + signature
	return []byte(final), nil

}

func (s *service) getTimeFrameMap(timeFrame string) string {
	timeFrames := map[string]string{
		"1minute":  "1m",
		"5minutes": "5m",
		"1hour":    "1h",
		"1day":     "1d",
	}

	return timeFrames[timeFrame]
}

func (s *service) FetchKlines(params handler.FetchKlinesParams) ([]handler.FetchKlinesResult, error) {
	result := make([]handler.FetchKlinesResult, 0)
	if !s.rateLimitHandler.canRequest(KlineReuqestWeight) {
		err := handler.CustomError{
			Type:    handler.ErrTypeRateLimit,
			Message: "ratelimit reached",
		}
		return result, err
	}
	ctx := context.Background()
	symbol := strings.Replace(params.Pair, "-", "", -1)
	from := strconv.FormatInt(params.From.Unix()*1000, 10) //multiply to 1000 to have millisecond too
	to := strconv.FormatInt(params.To.Unix()*1000, 10)     //multiply to 1000 to have millisecond too

	interval := s.getTimeFrameMap(params.TimeFrame)
	queryParams := url.Values{}
	queryParams.Set("symbol", symbol)
	queryParams.Set("interval", interval)
	queryParams.Set("startTime", from)
	queryParams.Set("endTime", to)

	queryParams.Set("limit", "1000")
	paramsString := queryParams.Encode()

	url := APIHost + FetchKlineURI + paramsString
	header := make(map[string]string)
	resp, respHeader, statusCode, err := s.getWithRetry(ctx, url, header)
	if err != nil {
		return result, err
	}
	if statusCode != http.StatusOK {
		s.logger.Warn("error in fetch klines from binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
		)
		return result, fmt.Errorf("status code is not 200 it is %d", statusCode)
	}
	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	resBody := fetchKlinesResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return result, err
	}

	for _, item := range resBody {
		//todo it would be better if we cast with two parameter and handle error
		//something like binanceStartTime,err := resBody[0].(int64)
		binanceStartTime := item[0].(float64)
		binanceEndTime := item[6].(float64)

		startTime := time.Unix(int64(binanceStartTime)/1000, 0)
		endTime := time.Unix(int64(binanceEndTime)/1000, 0)

		res := handler.FetchKlinesResult{
			StartTime:           startTime,
			EndTime:             endTime,
			OpenPrice:           item[1].(string),
			HighPrice:           item[2].(string),
			LowPrice:            item[3].(string),
			ClosePrice:          item[4].(string),
			BaseVolume:          item[5].(string),
			QuoteVolume:         item[7].(string),
			TakerBuyBaseVolume:  item[9].(string),
			TakerBuyQuoteVolume: item[10].(string),
		}
		result = append(result, res)
	}

	return result, nil
}
// getWithRetry performs an HTTP GET with up to 3 attempts (exponential backoff: 500ms, 1s).
// Use only for idempotent read requests; never call for writes.
func (s *service) getWithRetry(ctx context.Context, url string, header map[string]string) ([]byte, http.Header, int, error) {
	const maxAttempts = 3
	backoff := 500 * time.Millisecond
	var (
		resp       []byte
		respHeader http.Header
		statusCode int
		err        error
	)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}
		resp, respHeader, statusCode, err = s.httpClient.HTTPGet(ctx, url, header)
		if err == nil && statusCode == http.StatusOK {
			return resp, respHeader, statusCode, nil
		}
	}
	return resp, respHeader, statusCode, err
}

func (s *service) getTimestamp() int64 {
	serverTime := s.exchangeInfo.getServerTime() //in milliseconds
	fetchedAt := s.exchangeInfo.getFetchedAt()   //in milliseconds
	timeDiff := fetchedAt - serverTime
	now := time.Now().Unix() * 1000
	return now + int64(math.Abs(float64(timeDiff)))

}

func (s *service) FetchWithdrawals(params handler.FetchWithdrawalsParams) ([]handler.FetchWithdrawalsResult, error) {
	result := make([]handler.FetchWithdrawalsResult, 0)
	if !s.rateLimitHandler.canRequest(WithdrawHistoryReuqestWeight) {
		err := handler.CustomError{
			Type:    handler.ErrTypeRateLimit,
			Message: "ratelimit reached",
		}
		return result, err
	}
	ctx := context.Background()

	timestamp := s.getTimestamp()
	timestampString := strconv.FormatInt(timestamp, 10)

	startTime := strconv.FormatInt(params.From.Unix()*1000, 10)
	queryParams := url.Values{}
	queryParams.Set("startTime", startTime)
	queryParams.Set("timestamp", timestampString)
	paramsString := queryParams.Encode()

	//this is security constraint binance wants
	hash := hmac.New(sha256.New, []byte(s.metadata.Secret))
	hash.Write([]byte(paramsString))
	signature := hex.EncodeToString(hash.Sum(nil))
	finalParamsString := paramsString + "&signature=" + signature

	url := APIHost + FetchWithdrawalsURI + finalParamsString
	header := getHeader(s.metadata.APIKey, http.MethodGet)
	resp, respHeader, statusCode, err := s.getWithRetry(ctx, url, header)
	if err != nil {
		return result, err
	}

	if statusCode != http.StatusOK {
		s.logger.Warn("error in fetching withdrawals from binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
		)
		return result, fmt.Errorf("status code is not 200 it is %d", statusCode)
	}
	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	var resBody []fetchWithdrawalsResponseBody
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return result, err
	}

	for _, item := range resBody {
		data, err := json.Marshal(item)
		if err != nil {
			return result, err
		}

		res := handler.FetchWithdrawalsResult{
			TxID:               item.TxID,
			Status:             s.mapBinanceStatusToOurStatus(item.Status),
			ExternalWithdrawID: item.ID,
			Data:               string(data),
		}

		result = append(result, res)
	}

	return result, nil
}

func (s *service) FetchOrders(params handler.FetchOrdersParams) ([]handler.FetchOrdersResult, error) {
	result := make([]handler.FetchOrdersResult, 0)
	if !s.rateLimitHandler.canRequest(OrdersReuqestWeight) {
		err := handler.CustomError{
			Type:    handler.ErrTypeRateLimit,
			Message: "ratelimit reached",
		}
		return result, err
	}
	ctx := context.Background()
	timestamp := s.getTimestamp()
	timestampString := strconv.FormatInt(timestamp, 10)

	symbol := strings.Replace(params.Pair, "-", "", -1)

	queryParams := url.Values{}
	queryParams.Set("symbol", symbol)
	queryParams.Set("timestamp", timestampString)

	if params.From != int64(0) {
		startTime := strconv.FormatInt(params.From+1, 10)
		queryParams.Set("startTime", startTime)
	}

	paramsString := queryParams.Encode()
	//this is security constraint binance wants
	hash := hmac.New(sha256.New, []byte(s.metadata.Secret))
	hash.Write([]byte(paramsString))
	signature := hex.EncodeToString(hash.Sum(nil))
	finalParamsString := paramsString + "&signature=" + signature

	url := APIHost + FetchOrdersURI + finalParamsString
	header := getHeader(s.metadata.APIKey, http.MethodGet)
	resp, respHeader, statusCode, err := s.getWithRetry(ctx, url, header)
	if err != nil {
		return result, err
	}

	if statusCode != http.StatusOK {
		s.logger.Warn("error in fetching orders from binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
		)
		return result, fmt.Errorf("status code is not 200 it is %d", statusCode)
	}
	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	var resBody []fetchOrdersResponseBody
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return result, err
	}

	for _, item := range resBody {
		data, err := json.Marshal(item)
		if err != nil {
			return result, err
		}

		dateTime := time.Unix(int64(item.Time)/1000, 0)

		res := handler.FetchOrdersResult{
			OrderID:       item.OrderID,
			ClientOrderID: item.ClientOrderID,
			ExchangeType:  item.Type,
			Type:          item.Side,
			Price:         item.Price,
			Amount:        item.OrigQty,
			Status:        item.Status,
			DateTime:      dateTime,
			Timestamp:     item.Time,
			Data:          string(data),
		}

		result = append(result, res)
	}

	return result, nil
}

func (s *service) FetchTrades(params handler.FetchTradesParams) ([]handler.FetchTradesResult, error) {
	result := make([]handler.FetchTradesResult, 0)
	if !s.rateLimitHandler.canRequest(TradesReuqestWeight) {
		err := handler.CustomError{
			Type:    handler.ErrTypeRateLimit,
			Message: "ratelimit reached",
		}
		return result, err
	}
	ctx := context.Background()
	timestamp := s.getTimestamp()
	timestampString := strconv.FormatInt(timestamp, 10)

	symbol := strings.Replace(params.Pair, "-", "", -1)

	queryParams := url.Values{}
	queryParams.Set("symbol", symbol)
	queryParams.Set("timestamp", timestampString)

	if params.From != int64(0) {
		startTime := strconv.FormatInt(params.From+1, 10)
		queryParams.Set("startTime", startTime)
	}

	paramsString := queryParams.Encode()

	//this is security constraint binance wants
	hash := hmac.New(sha256.New, []byte(s.metadata.Secret))
	hash.Write([]byte(paramsString))
	signature := hex.EncodeToString(hash.Sum(nil))
	finalParamsString := paramsString + "&signature=" + signature

	url := APIHost + FetchTradesURI + finalParamsString
	header := getHeader(s.metadata.APIKey, http.MethodGet)
	resp, respHeader, statusCode, err := s.getWithRetry(ctx, url, header)
	if err != nil {
		return result, err
	}

	if statusCode != http.StatusOK {
		s.logger.Warn("error in fetching trades from binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
		)
		return result, fmt.Errorf("status code is not 200 it is %d", statusCode)
	}
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	var resBody []fetchTradesResponseBody
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return result, err
	}

	for _, item := range resBody {
		data, err := json.Marshal(item)
		if err != nil {
			return result, err
		}

		dateTime := time.Unix(int64(item.Time)/1000, 0)

		res := handler.FetchTradesResult{
			ID:         item.ID,
			Price:      item.Price,
			Amount:     item.Qty,
			Commission: item.Commission,
			Coin:       item.CommissionAsset,
			DateTime:   dateTime,
			Timestamp:  item.Time,
			OrderID:    item.OrderID,
			Data:       string(data),
		}

		result = append(result, res)
	}

	return result, nil
}

func (s *service) mapBinanceStatusToOurStatus(binanceStatus int) string {
	return binanceTxStatusToOurStatus[binanceStatus]
}

func (s *service) Withdraw(params handler.WithdrawParams) (res handler.WithdrawResult, err error) {
	ctx := context.Background()
	url := APIHost + WithdrawURI
	body, err := s.getWithdrawRequestBody(params)
	if err != nil {
		return res, err
	}
	header := getHeader(s.metadata.APIKey, http.MethodPost)
	resp, respHeader, statusCode, err := s.httpClient.HTTPPost(ctx, url, body, header)
	if err != nil {
		return res, err
	}
	if statusCode != http.StatusOK {
		s.logger.Warn("error in withdraw request in binance",
			zap.Int("statusCode", statusCode),
			zap.String("body", string(resp)),
			zap.String("coin", params.Coin),
			zap.String("amount", params.Amount),
			zap.String("toAddress", params.ToAddress),
			zap.String("network", params.Network),
		)

		return res, fmt.Errorf("status code is not 200 it is %d with body %s", statusCode, string(resp))
	}
	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()
	resBody := withdrawResponseBody{}
	err = json.Unmarshal(resp, &resBody)
	if err != nil {
		return res, err
	}
	res = handler.WithdrawResult{
		ID:           resBody.ID,
		ErrorMessage: "",
	}
	return res, nil
}

func (s *service) getWithdrawRequestBody(params handler.WithdrawParams) (body []byte, err error) {
	timestamp := s.getTimestamp()
	rb := withdrawRequestBody{
		Coin:      params.Coin,
		Address:   params.ToAddress,
		Amount:    params.Amount,
		Network:   params.Network,
		Timestamp: timestamp,
	}

	v, err := query.Values(rb)
	if err != nil {
		return body, err
	}

	hash := hmac.New(sha256.New, []byte(s.metadata.Secret))

	hash.Write([]byte(v.Encode()))
	signature := hex.EncodeToString(hash.Sum(nil))

	final := v.Encode() + "&signature=" + signature
	return []byte(final), nil
}

func NewBinanceService(httpClient platform.HTTPClient, rc platform.RedisClient, configs platform.Configs, logger platform.Logger, metadataString string) handler.ExchangeHandler {
	md := metadata{}
	if strings.HasPrefix(metadataString, "\"") {
		metadataString = metadataString[1:]
	}

	if strings.HasSuffix(metadataString, "\"") {
		metadataString = metadataString[0 : len(metadataString)-1]
	}

	metadataString = strings.Replace(metadataString, "\\", "", -1)

	err := json.Unmarshal([]byte(metadataString), &md)

	if err != nil {
		panic("we should never reach here if our metada is set correctly in database")
	}

	rateLimitHandler := newRateLimitHandler(logger)
	exchangeInfo := newExchangeInfoService(rc, httpClient, rateLimitHandler, configs, logger, md.APIKey)
	return &service{
		httpClient:       httpClient,
		configs:          configs,
		logger:           logger,
		exchangeInfo:     exchangeInfo,
		rateLimitHandler: rateLimitHandler,
		metadata:         md,
	}

}
