package repository

import (
	"exchange-go/internal/externalexchange"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderFromExternalRepository struct {
	db *gorm.DB
}

func (r *orderFromExternalRepository) GetLastOrderFromExternalByPairID(pairID int64, orderFromExternal *externalexchange.OrderFromExternal) error {
	return r.db.Where("order_from_external.pair_currency_id = ?", pairID).Order("timestamp desc").Offset(0).Limit(1).Find(orderFromExternal).Error
}

func (r *orderFromExternalRepository) Create(orderFromExternal *externalexchange.OrderFromExternal) error {
	return r.db.Omit(clause.Associations).Create(orderFromExternal).Error
}

func (r *orderFromExternalRepository) GetOrderByExternalOrderID(externalOrderID int64, orderFromExternal *externalexchange.OrderFromExternal) error {
	return r.db.Where("external_order_id = ?", externalOrderID).First(orderFromExternal).Error
}

func NewOrderFromExternalRepository(db *gorm.DB) externalexchange.OrderFromExternalRepository {
	return &orderFromExternalRepository{db}
}
