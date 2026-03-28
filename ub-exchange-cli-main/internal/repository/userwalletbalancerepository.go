package repository

import (
	"exchange-go/internal/userbalance"

	"gorm.io/gorm"
)

type userWalletBalanceRepository struct {
	db *gorm.DB
}

func (r *userWalletBalanceRepository) FindUserWalletBalance(userID int, coinID int64, networkCoinID int64, userWalletBalance *userbalance.UserWalletBalance) error {

	tx := r.db.Where("user_id = ? AND currency_id = ?", userID, coinID)

	if networkCoinID != 0 {
		tx.Where("network_currency_id = ?", networkCoinID)
	}

	return tx.First(userWalletBalance).Error
}

func NewUserWalletBalanceRepository(db *gorm.DB) userbalance.UserWalletBalanceRepository {
	return &userWalletBalanceRepository{db}
}
