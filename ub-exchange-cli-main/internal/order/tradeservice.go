package order

import "time"

// TradeService provides query methods for retrieving trade records.
type TradeService interface {
	// GetTradesOfUserBetweenTimes returns all trades for a user within the specified time range.
	GetTradesOfUserBetweenTimes(userID int, startTime, endTime string) []Trade
	// GetBotTradesByIDAndCreatedAtGreaterThan returns bot trades for a pair created after
	// the given trade ID and timestamp.
	GetBotTradesByIDAndCreatedAtGreaterThan(pairID int64, tradeID int64, createdAt time.Time) []Trade
}

type tradeService struct {
	repo TradeRepository
}

func (s *tradeService) GetTradesOfUserBetweenTimes(userID int, startTime, endTime string) []Trade {
	return s.repo.GetTradesOfUserBetweenTimes(userID, startTime, endTime)
}

func (s *tradeService) GetBotTradesByIDAndCreatedAtGreaterThan(pairID int64, tradeID int64, createdAt time.Time) []Trade {
	return s.repo.GetBotTradesByIDAndCreatedAtGreaterThan(pairID, tradeID, createdAt)
}

func NewTradeService(repo TradeRepository) TradeService {
	return &tradeService{repo}
}
