package repository

import (
	"exchange-go/internal/currency"
	"gorm.io/gorm"
)

type currencyRepository struct {
	db *gorm.DB
}

func (cr *currencyRepository) GetCoinsAlphabetically() []currency.Coin {
	var coins []currency.Coin
	cr.db.Where(currency.Coin{IsActive: true}).Order("code asc").Find(&coins)
	return coins
}

func (cr *currencyRepository) GetActiveCoins() []currency.Coin {
	var coins []currency.Coin
	cr.db.Where(currency.Coin{IsActive: true}).Order("priority desc").Find(&coins)
	return coins
}

func (cr *currencyRepository) GetCoinByCode(code string, coin *currency.Coin) error {
	return cr.db.Where("code = ? ", code).First(coin).Error
}

func NewCurrencyRepository(db *gorm.DB) currency.Repository {
	return &currencyRepository{db}
}
