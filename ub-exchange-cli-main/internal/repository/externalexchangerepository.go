package repository

import (
	"exchange-go/internal/externalexchange"
	"gorm.io/gorm"
)

type externalExchangeRepository struct {
	db *gorm.DB
}

func (er *externalExchangeRepository) GetEnabledPrivateExternalExchange(ee *externalexchange.ExternalExchange) error {
	filters := externalexchange.ExternalExchange{Status: externalexchange.StatusEnabled, Type: externalexchange.TypePrivate}
	return er.db.Where(filters).First(&ee).Error
}

func NewExternalExchangeRepository(db *gorm.DB) externalexchange.Repository {
	return &externalExchangeRepository{db}
}
