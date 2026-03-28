package candle

import (
	context "context"
	"exchange-go/internal/platform"
	"time"

	grpc "google.golang.org/grpc"
)

type AveragePairPrice struct {
	PairName string
	Price    string
	Day      string
}

type HighAndLowPrice struct {
	High string
	Low  string
}

type CandleTrend struct {
	Pair      string
	Price     string
	StartTime string
	EndTime   string
}
// CandleGRPCClient communicates with the external candle-price gRPC service to
// retrieve historical OHLC candlestick data for trading pairs.
type CandleGRPCClient interface {
	// GetLastPriceForPair returns the closing price of the most recent candle
	// for the given pair at or before currentTime (Unix timestamp).
	GetLastPriceForPair(pairName string, currentTime int64) (string, error)
	// GetAveragePriceOfPairs returns average prices for each pair over the specified
	// time range (Unix timestamps).
	GetAveragePriceOfPairs(pairNames []string, fromTime int64, toTime int64) ([]AveragePairPrice, error)
	// GetHighAndLowPriceForPairFromDate returns the highest and lowest prices for
	// a pair since fromTime (Unix timestamp).
	GetHighAndLowPriceForPairFromDate(pairName string, fromTime int64) (HighAndLowPrice, error)
	// GetCandleTrends retrieves candlestick trend data for multiple pairs within
	// the given time range and time frame.
	GetCandleTrends(pairNames []string, timeFrame string, fromTime int64, toTime int64) ([]CandleTrend, error)
}
type candleGRPCClient struct {
	configs     platform.Configs
	logger      platform.Logger
	grpcAddress string
}

func (c *candleGRPCClient) GetLastPriceForPair(pairName string, currentTime int64) (string, error) {
	conn, err := grpc.Dial(c.grpcAddress, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer conn.Close()
	client := NewCandlePriceClient(conn)
	req := &LastCandleRequest{
		Pair:      pairName,
		Timestamp: currentTime,
	}
	res, err := client.GetLastCandle(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.Close, nil

}

func (c *candleGRPCClient) GetAveragePriceOfPairs(pairNames []string, fromTime int64, toTime int64) ([]AveragePairPrice, error) {
	prices := make([]AveragePairPrice, 0)
	if c.configs.GetEnv() == platform.EnvTest {
		lastDay := time.Now().Add(-1 * 24 * time.Hour)
		prices = []AveragePairPrice{
			{
				PairName: "BTC-USDT",
				Price:    "50000.0",
				Day:      lastDay.Format("2006-01-02"),
			},
			{
				PairName: "ETH-BTC",
				Price:    "0.04",
				Day:      lastDay.Format("2006-01-02"),
			},
		}
		return prices, nil
	}
	conn, err := grpc.Dial(c.grpcAddress, grpc.WithInsecure())
	if err != nil {
		return prices, err
	}
	defer conn.Close()
	client := NewCandlePriceClient(conn)
	req := &AveragePriceOfPairsRequest{
		Pairs:    pairNames,
		FromTime: fromTime,
		ToTime:   toTime,
	}
	res, err := client.GetAveragePriceOfPairs(context.Background(), req)
	if err != nil {
		return prices, err
	}
	for _, p := range res.AveragePrices {
		price := AveragePairPrice{
			PairName: p.Pair,
			Price:    p.Price,
			Day:      p.Day,
		}
		prices = append(prices, price)
	}
	return prices, nil

}

func (c *candleGRPCClient) GetHighAndLowPriceForPairFromDate(pairName string, fromTime int64) (HighAndLowPrice, error) {
	highAndLowPrice := HighAndLowPrice{}
	conn, err := grpc.Dial(c.grpcAddress, grpc.WithInsecure())
	if err != nil {
		return highAndLowPrice, err
	}
	defer conn.Close()
	client := NewCandlePriceClient(conn)
	req := &HighAndLowRequest{
		Pair:      pairName,
		Timestamp: fromTime,
	}
	res, err := client.GetHighAndLowPriceOfPairFromDate(context.Background(), req)
	if err != nil {
		return highAndLowPrice, err
	}
	highAndLowPrice.High = res.High
	highAndLowPrice.Low = res.Low
	return highAndLowPrice, nil

}
func (c *candleGRPCClient) GetCandleTrends(pairNames []string, timeFrame string, fromTime int64, toTime int64) ([]CandleTrend, error) {
	trends := make([]CandleTrend, 0)
	if c.configs.GetEnv() == platform.EnvTest {
		trends := []CandleTrend{
			{
				Pair:      "BTC-USDT",
				Price:     "30000.00000000",
				StartTime: "",
				EndTime:   "",
			},
			{
				Pair:      "ETH-USDT",
				Price:     "2000.00000000",
				StartTime: "",
				EndTime:   "",
			},
		}
		return trends, nil
	}
	conn, err := grpc.Dial(c.grpcAddress, grpc.WithInsecure())
	if err != nil {
		return trends, err
	}
	defer conn.Close()
	client := NewCandlePriceClient(conn)
	req := &GetCandlesRequest{
		Pairs:     pairNames,
		TimeFrame: timeFrame,
		FromTime:  fromTime,
		ToTime:    toTime,
	}
	res, err := client.GetCandles(context.Background(), req)
	if err != nil {
		return trends, err
	}
	for _, c := range res.Candles {
		trend := CandleTrend{
			Pair:      c.Pair,
			Price:     c.Price,
			StartTime: c.StartTime,
			EndTime:   c.EndTime,
		}
		trends = append(trends, trend)
	}
	return trends, nil
}

func NewCandleGRPCClient(configs platform.Configs, logger platform.Logger) CandleGRPCClient {
	address := configs.GetString("candle.grpcaddr")
	return &candleGRPCClient{
		configs:     configs,
		logger:      logger,
		grpcAddress: address,
	}
}
