package repository

import (
	"exchange-go/internal/currency"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type klineSyncRepository struct {
	db *gorm.DB
}

func (r *klineSyncRepository) Create(klineSync *currency.KlineSync) error {
	return r.db.Omit(clause.Associations).Create(klineSync).Error
}

func (r *klineSyncRepository) Update(klineSync *currency.KlineSync) error {
	return r.db.Omit(clause.Associations).Save(klineSync).Error
}

func (r *klineSyncRepository) GetKlineSyncsByStatusAndLimit(status string, limit int) []currency.KlineSync {
	var klineSyncs []currency.KlineSync
	r.db.Where("status = ?", status).Order("id asc").Offset(0).Limit(limit).Find(&klineSyncs)
	return klineSyncs
}

func NewKlineSyncRepository(db *gorm.DB) currency.KlineSyncRepository {
	return &klineSyncRepository{db}
}
