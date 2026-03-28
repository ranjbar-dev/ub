package seed

import (
	"exchange-go/internal/externalexchange"
	"gorm.io/gorm"
)

func ExternalExchangeSeed(db *gorm.DB) {
	ee := externalexchange.ExternalExchange{
		ID:       1,
		Name:     "binance",
		MetaData: "{}",
		Status:   "enabled",
		Type:     "private",
	}

	db.Create(&ee)

}
