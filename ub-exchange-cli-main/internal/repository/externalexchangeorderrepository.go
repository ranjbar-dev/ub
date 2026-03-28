package repository

import (
	"exchange-go/internal/externalexchange"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type externalExchangeOrderRepository struct {
	db *gorm.DB
}

func (r *externalExchangeOrderRepository) Create(o *externalexchange.Order) error {
	return r.db.Omit(clause.Associations).Create(o).Error
}

func (r *externalExchangeOrderRepository) Update(o *externalexchange.Order) error {
	return r.db.Omit(clause.Associations).Save(o).Error
}

func (r *externalExchangeOrderRepository) GetExternalExchangeOrdersLastTradeIds() []externalexchange.LastTradeIDAndPair {
	var result []externalexchange.LastTradeIDAndPair
	r.db.Raw("SELECT pair_currency_id AS PairID,MAX(last_trade_id) as TradeID FROM external_exchange_orders GROUP BY PairId").Scan(&result)
	return result
}

func NewExternalExchangeOrderRepository(db *gorm.DB) externalexchange.OrderRepository {
	return &externalExchangeOrderRepository{db}
}
