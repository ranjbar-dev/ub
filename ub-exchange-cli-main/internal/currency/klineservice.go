package currency

import (
	"context"
	"database/sql"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/livedata"
	"time"
)

//Time frames
const (
	Timeframe1minute     = "1minute"
	Timeframe3minutes    = "3minutes"
	Timeframe5minutes    = "5minutes"
	Timeframe15minutes   = "15minutes"
	Timeframe30minutes   = "30minutes"
	Timeframe45minutes   = "45minutes"
	Timeframe1hour       = "1hour"
	Timeframe2hours      = "2hours"
	Timeframe3hours      = "3hours"
	Timeframe4hours      = "4hours"
	Timeframe1day        = "1day"
	Timeframe1week       = "1week"
	Timeframe1month      = "1month"
	SyncStatusCreated    = "CREATED"
	SyncStatusProcessing = "PROCESSING"
	SyncStatusDone       = "DONE"
	SyncStatusError      = "ERROR"
	SyncTypeManual       = "MANUAL"
	SyncTypeAuto         = "AUTO"

	//MAX period allowed for each time-frame
	MaxPeriod1minute  = 604800    //7 days
	MaxPeriod5minutes = 4320000   //50 days
	MaxPeriod1hour    = 31104000  //360 days
	MaxPeriod1day     = 315360000 //3650 days
)

type PartialKline struct {
	TimeFrame           string
	StartTime           string
	EndTime             string
	OpenPrice           string
	ClosePrice          string
	LowPrice            string
	HighPrice           string
	BaseVolume          string
	QuoteVolume         string
	TakerBuyBaseVolume  string
	TakerBuyQuoteVolume string
}

type CreateKlineSyncParams struct {
	StartTime  string
	EndTime    string
	PairID     int64
	TimeFrame  string
	WithUpdate bool
	Type       string
}

type KlineTrendData struct {
	PairID int64  `json:"-"`
	Price  string `json:"price"`
	Time   string `json:"time"`
}

// KlineService provides OHLC (candlestick) data operations via gRPC, including
// price lookups, trend retrieval, and kline sync task management.
type KlineService interface {
	// GetLastPriceForPair returns the most recent closing price for a pair at the given date.
	GetLastPriceForPair(ctx context.Context, pairName string, date time.Time) (string, error)
	// GetKlineTrends retrieves candlestick trend data for the specified pairs and time range.
	GetKlineTrends(pairNames []string, timeFrame string, from time.Time, to time.Time) ([]candle.CandleTrend, error)
	// GetHighAndLowPriceFromDateForPairByPairName returns the highest and lowest prices
	// for a pair since the given date.
	GetHighAndLowPriceFromDateForPairByPairName(pairName string, fromDate time.Time) (candle.HighAndLowPrice, error)
	// CreateKlineSync creates a new kline synchronization task with the given parameters.
	CreateKlineSync(params CreateKlineSyncParams) error
	// GetActiveTimeFrames returns the list of supported candlestick time frames.
	GetActiveTimeFrames() []string
	// GetKlineSyncsByStatusAndLimit returns up to limit sync tasks matching the given status.
	GetKlineSyncsByStatusAndLimit(status string, limit int) []KlineSync
	// UpdateKlineSync persists changes to an existing kline sync task.
	UpdateKlineSync(klineSync *KlineSync) error
	// GetAveragePriceOfPairsUsingTime returns average prices for pairs over the given time range.
	GetAveragePriceOfPairsUsingTime(pairNames []string, startTime time.Time, endTime time.Time) ([]candle.AveragePairPrice, error)
}

type klineService struct {
	klineSyncRepository      KlineSyncRepository
	liveDataService          livedata.Service
	candleGRPCClient         candle.CandleGRPCClient
	resolutionToTimeFrameMap map[string]string
}

func (s *klineService) GetLastPriceForPair(ctx context.Context, pairName string, date time.Time) (string, error) {
	return s.candleGRPCClient.GetLastPriceForPair(pairName, date.Unix())
}

func (s *klineService) GetHighAndLowPriceFromDateForPairByPairName(pairName string, fromDate time.Time) (candle.HighAndLowPrice, error) {
	return s.candleGRPCClient.GetHighAndLowPriceForPairFromDate(pairName, fromDate.Unix())
}

func (s *klineService) CreateKlineSync(params CreateKlineSyncParams) error {
	startTime, err := time.Parse("2006-01-02 15:04:05", params.StartTime)
	if err != nil {
		return err
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", params.EndTime)
	if err != nil {
		return err
	}

	ks := &KlineSync{
		PairID:     params.PairID,
		TimeFrame:  sql.NullString{String: params.TimeFrame, Valid: true},
		StartTime:  sql.NullTime{Time: startTime, Valid: true},
		EndTime:    sql.NullTime{Time: endTime, Valid: true},
		Status:     SyncStatusCreated,
		Type:       params.Type,
		WithUpdate: params.WithUpdate,
	}

	return s.klineSyncRepository.Create(ks)
}

func (s *klineService) GetActiveTimeFrames() []string {
	return []string{
		Timeframe1minute,
		Timeframe5minutes,
		Timeframe1hour,
		Timeframe1day,
	}
}

func (s *klineService) GetKlineSyncsByStatusAndLimit(status string, limit int) []KlineSync {
	return s.klineSyncRepository.GetKlineSyncsByStatusAndLimit(status, limit)
}

func (s *klineService) UpdateKlineSync(klineSync *KlineSync) error {
	return s.klineSyncRepository.Update(klineSync)
}

func (s *klineService) GetAveragePriceOfPairsUsingTime(pairNames []string, startTime time.Time, endTime time.Time) ([]candle.AveragePairPrice, error) {
	return s.candleGRPCClient.GetAveragePriceOfPairs(pairNames, startTime.Unix(), endTime.Unix())
}

func (s *klineService) GetKlineTrends(pairNames []string, timeFrame string, from time.Time, to time.Time) ([]candle.CandleTrend, error) {
	return s.candleGRPCClient.GetCandleTrends(pairNames, timeFrame, from.Unix(), to.Unix())
}

func NewKlineService(klineSyncRepository KlineSyncRepository, liveDataService livedata.Service,
	candleGRPCClient candle.CandleGRPCClient) KlineService {

	resolutionTimeFrameMap := map[string]string{
		"1":   Timeframe1minute,
		"3":   Timeframe3minutes,
		"5":   Timeframe5minutes,
		"15":  Timeframe15minutes,
		"30":  Timeframe30minutes,
		"45":  Timeframe45minutes,
		"60":  Timeframe1hour,
		"120": Timeframe2hours,
		"180": Timeframe3hours,
		"240": Timeframe4hours,
		"1D":  Timeframe1day,
		"1W":  Timeframe1week,
		"1M":  Timeframe1month,
	}

	return &klineService{
		klineSyncRepository:      klineSyncRepository,
		liveDataService:          liveDataService,
		candleGRPCClient:         candleGRPCClient,
		resolutionToTimeFrameMap: resolutionTimeFrameMap,
	}
}
