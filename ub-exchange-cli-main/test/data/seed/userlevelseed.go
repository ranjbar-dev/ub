package seed

import (
	"exchange-go/internal/user"

	"gorm.io/gorm"
)

func UserLevelSeed(db *gorm.DB) {
	l1 := user.Level{
		ID:                        1,
		Name:                      "vip1",
		ExchangeNumberLimit:       100,
		ExchangeVolumeLimitAmount: "10",
		TakerFeePercentage:        1.0,
		MakerFeePercentage:        1.0,
		MinExchangeVolume:         0.01,
		MaxExchangeVolume:         0.1,
	}

	l2 := user.Level{
		ID:                        2,
		Name:                      "vip2",
		ExchangeNumberLimit:       100,
		ExchangeVolumeLimitAmount: "10",
		TakerFeePercentage:        1.0,
		MakerFeePercentage:        1.0,
		MinExchangeVolume:         0.1,
		MaxExchangeVolume:         0.3,
	}

	db.Create(&l1)
	db.Create(&l2)

}
