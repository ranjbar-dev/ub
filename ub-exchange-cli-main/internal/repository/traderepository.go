package repository

import (
	"exchange-go/internal/order"
	"time"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type tradeRepository struct {
	db *gorm.DB
}

func (tr *tradeRepository) GetTradesOfUserBetweenTimes(userID int, startTime, endTime string) []order.Trade {
	var trades []order.Trade
	q := tr.db.
		Joins("join pair_currencies AS pairs on trades.pair_currency_id = pairs.id").
		Joins("join orders  AS o on trades.buy_order_id = o.id OR trades.sell_order_id = o.id").
		Where("o.creator_user_id = ?", userID).
		Where("trades.created_at >= ?", startTime).
		Where("trades.created_at <= ?", endTime)
	q.Find(&trades)
	return trades
}

func (tr *tradeRepository) Create(t *order.Trade) error {
	return tr.db.Omit(clause.Associations).Create(t).Error
}

func (tr *tradeRepository) GetBotTradesByIDAndCreatedAtGreaterThan(pairID int64, tradeID int64, createdAt time.Time) []order.Trade {
	var trades []order.Trade
	tr.db.Where("id > ? and pair_currency_id = ? and created_at >= ?", tradeID, pairID, createdAt).Find(&trades)
	return trades
}

func NewTradeRepository(db *gorm.DB) order.TradeRepository {
	return &tradeRepository{db}
}
