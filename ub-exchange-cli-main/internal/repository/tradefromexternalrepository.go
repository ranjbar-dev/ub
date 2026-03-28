package repository

import (
	"exchange-go/internal/externalexchange"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tradeFromExternalRepository struct {
	db *gorm.DB
}

func (r *tradeFromExternalRepository) GetLastTradeFromExternalByPairID(pairID int64, tradeFromExternal *externalexchange.TradeFromExternal) error {
	return r.db.Joins("JOIN order_from_external ON order_from_external.id = trade_from_external.order_id AND order_from_external.pair_currency_id = ?", pairID).Order("trade_from_external.timestamp DESC").Offset(0).Limit(1).Find(tradeFromExternal).Error
}

func (r *tradeFromExternalRepository) Create(tradeFromExternal *externalexchange.TradeFromExternal) error {
	return r.db.Omit(clause.Associations).Create(tradeFromExternal).Error
}

func NewTradeFromExternalRepository(db *gorm.DB) externalexchange.TradeFromExternalRepository {
	return &tradeFromExternalRepository{db}
}
