package binance

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	FilterTypeLotSize       = "LOT_SIZE"
	FilterTypeMarketLotSize = "MARKET_LOT_SIZE"
	FilterTypeMinNotional   = "MIN_NOTIONAL"
)

type rateLimitData struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int64  `json:"limit"`
}

//type filter map[string]string

type symbol struct {
	Symbol                     string                   `json:"symbol"`
	Status                     string                   `json:"status"`
	BaseAsset                  string                   `json:"baseAsset"`
	BaseAssetPrecision         int                      `json:"baseAssetPrecision"`
	QuoteAsset                 string                   `json:"quoteAsset"`
	QuotePrecision             int                      `json:"quotePrecision"`
	QuoteAssetPrecision        int                      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int                      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int                      `json:"quoteCommissionPrecision"`
	OrderTypes                 []string                 `json:"orderTypes"`
	IcebergAllowed             bool                     `json:"icebergAllowed"`
	OcoAllowed                 bool                     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool                     `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool                     `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool                     `json:"isMarginTradingAllowed"`
	Filters                    []map[string]interface{} `json:"filters"`
	Permissions                []string                 `json:"permissions"`
}

type exchangeInfo struct {
	Timezone   string          `json:"timezone"`
	ServerTime int64           `json:"serverTime"`
	RateLimits []rateLimitData `json:"rateLimits"`
	Symbols    []symbol        `json:"symbols"`
	fetchedAt  int64
	//exchangeFilters []string
}

type exchangeInfoService struct {
	rc               platform.RedisClient
	httpClient       platform.HTTPClient
	rateLimitHandler *rateLimitHandler
	logger           platform.Logger
	info             exchangeInfo
	apiKey           string
}

func (s *exchangeInfoService) getServerTime() int64 {
	return s.info.ServerTime
}

func (s *exchangeInfoService) getFetchedAt() int64 {
	return s.info.fetchedAt
}

func (s *exchangeInfoService) getLotSizeForSymbol(symbolName string) (string, error) {
	symbols := s.info.Symbols
	for _, symbol := range symbols {
		if symbol.Symbol == symbolName {
			for _, filter := range symbol.Filters {
				if filter["filterType"] == FilterTypeLotSize {
					stepSize, _ := filter["stepSize"].(string)
					return stepSize, nil
				}
			}

		}
	}

	return "", fmt.Errorf("symbol not found")
}

func (s *exchangeInfoService) getLotSizeMinAndMaxQty(symbolName string) (string, string, error) {
	symbols := s.info.Symbols
	for _, symbol := range symbols {
		if symbol.Symbol == symbolName {
			for _, filter := range symbol.Filters {
				if filter["filterType"] == FilterTypeLotSize {
					minQty, _ := filter["minQty"].(string)
					maxQty, _ := filter["maxQty"].(string)
					return minQty, maxQty, nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("symbol not found")
}

func (s *exchangeInfoService) getMarketLotSizeFilters(symbolName string) (string, string, string, error) {
	symbols := s.info.Symbols
	for _, symbol := range symbols {
		if symbol.Symbol == symbolName {
			for _, filter := range symbol.Filters {
				if filter["filterType"] == FilterTypeMarketLotSize {
					minQty, _ := filter["minQty"].(string)
					maxQty, _ := filter["maxQty"].(string)
					stepSize, _ := filter["stepSize"].(string)
					return minQty, maxQty, stepSize, nil
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("symbol not found")
}

func (s *exchangeInfoService) getMinNotional(symbolName string) (string, bool, error) {
	symbols := s.info.Symbols
	for _, symbol := range symbols {
		if symbol.Symbol == symbolName {
			for _, filter := range symbol.Filters {
				if filter["filterType"] == FilterTypeMinNotional {
					min, _ := filter["minNotional"].(string)
					applyToMarket, _ := filter["applyToMarket"].(bool)
					return min, applyToMarket, nil
				}
			}
		}
	}

	return "", false, fmt.Errorf("symbol not found")
}

func (s *exchangeInfoService) getFromAPI() {
	if !s.rateLimitHandler.canRequest(ExchangeInfoReuqestWeight) {
		return
	}
	ctx := context.Background()
	url := APIHost + ExchangeInfoURI
	headers := getHeader(s.apiKey, http.MethodGet)
	body, respHeader, statusCode, err := s.httpClient.HTTPGet(ctx, url, headers)
	if err != nil {
		s.logger.Error2("can not get exchange info from binance", err,
			zap.String("service", "exchangeInfoService"),
			zap.String("method", "getFromApi"),
			zap.Int("statusCode", statusCode),
		)
		return
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("status code is not 200 it is %d with body %s", statusCode, string(body))
		s.logger.Error2("unsuccessful response from binance", err,
			zap.String("service", "exchangeInfoService"),
			zap.String("method", "getFromApi"),
			zap.Int("statusCode", statusCode),
		)
		return
	}
	//recalculating ratelimit
	defer func() {
		go s.rateLimitHandler.updateUsingHeader(respHeader)
	}()

	bei := exchangeInfo{}

	err = json.Unmarshal([]byte(body), &bei)
	if err != nil {
		s.logger.Error2("can not unmarshal response from binance", err,
			zap.String("service", "exchangeInfoService"),
			zap.String("method", "getFromApi"),
			zap.String("body", string(body)),
		)
	}

	bei.fetchedAt = time.Now().Unix() * 1000 //in milliseconds

	s.info = bei
}

func (s *exchangeInfoService) loadFromRedis() {
}

func (s *exchangeInfoService) getFilters() {
	return
}

func (s *exchangeInfoService) updateExchangeInfo() {
	s.getFromAPI()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				s.getFromAPI()
			}
		}
	}()

}

func newExchangeInfoService(rc platform.RedisClient, httpClient platform.HTTPClient, rateLimitHandler *rateLimitHandler, configs platform.Configs, logger platform.Logger, apiKey string) *exchangeInfoService {
	s := &exchangeInfoService{
		rc:               rc,
		httpClient:       httpClient,
		rateLimitHandler: rateLimitHandler,
		logger:           logger,
		apiKey:           apiKey,
	}
	if configs.GetEnv() == platform.EnvProd {
		s.updateExchangeInfo()
	}
	return s
}
