package repository

import (
	"exchange-go/internal/userbalance"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userBalanceRepository struct {
	db *gorm.DB
}

func (ubr *userBalanceRepository) GetBalancesOfUsersForCoins(userIds []int, coinIds []int64) []userbalance.UserBalance {
	var balances []userbalance.UserBalance
	ubr.db.Where("user_id IN ?", userIds).Where("currency_id IN ?", coinIds).Find(&balances)
	return balances
}

func (ubr *userBalanceRepository) GetBalancesOfUsersForCoinsUsingTx(tx *gorm.DB, userIds []int, coinIds []int64) []userbalance.UserBalance {
	var balances []userbalance.UserBalance
	tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id IN ?", userIds).Where("currency_id IN ?", coinIds).Find(&balances)
	return balances
}

func (ubr *userBalanceRepository) GetBalanceOfUserByCoinID(userID int, coinID int64, ub *userbalance.UserBalance) error {
	//todo this join should be here but it does not work in tests try to debug this later
	//return pr.db.Joins("Coin").Where(userbalance.UserBalance{UserId: userId, CoinId: coinId}).First(ub).Error
	return ubr.db.Where(userbalance.UserBalance{UserID: userID, CoinID: coinID}).First(ub).Error
}

func (ubr *userBalanceRepository) GetBalanceOfUserByCoinIDUsingTx(tx *gorm.DB, userID int, coinID int64, ub *userbalance.UserBalance) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(userbalance.UserBalance{UserID: userID, CoinID: coinID}).First(ub).Error
}

func (ubr *userBalanceRepository) GetUserBalancesForCoins(userID int, coinIds []int64) []userbalance.UserBalance {
	var balances []userbalance.UserBalance
	ubr.db.Joins("Coin").Where("user_id = ?", userID).Where("currency_id IN ?", coinIds).Find(&balances)
	return balances
}

func (ubr *userBalanceRepository) GetUserAllBalances(userID int, filters userbalance.AllBalancesFilters) []userbalance.UserBalance {
	var balances []userbalance.UserBalance
	q := ubr.db.Joins("Coin").Where("user_id = ?", userID)

	if filters.CoinName != "" {
		q.Where("Coin.name like ?", "%"+filters.CoinName+"%")
	}

	if filters.CoinCode != "" {
		q.Where("Coin.code like ?", "%"+filters.CoinCode+"%")
	}

	q.Order("Coin.priority desc")

	q.Find(&balances)
	return balances
}

func (ubr *userBalanceRepository) GetBalancesWithoutAddresses(filters map[string]interface{}) []userbalance.UserBalance {
	var balances []userbalance.UserBalance
	q := ubr.db.Joins("Coin").Joins("User").Where("address IS NULL OR address= ?", "")
	if userIDString, ok := filters["userId"]; ok {
		userID, _ := strconv.ParseInt(userIDString.(string), 10, 64)
		q.Where("user_id = ?", userID)
	}

	page := filters["page"].(int)
	pageSize := filters["pageSize"].(int)
	q.Offset(page * pageSize).Limit(pageSize).Order("id asc").Find(&balances)
	return balances
}

func (u *userBalanceRepository) GetUserBalanceByCoinAndAddressUsingTx(tx *gorm.DB, coinID int64, address string, ub *userbalance.UserBalance) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("currency_id = ? AND (address = ? OR JSON_CONTAINS(other_addresses,JSON_OBJECT(\"address\",?)) = 1)", coinID, address, address).First(ub).Error

}
func (u *userBalanceRepository) GetUserBalanceByIDUsingTx(tx *gorm.DB, id int64, ub *userbalance.UserBalance) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Joins("Coin").Where("user_balances.id = ?", id).First(ub).Error
}

func NewUserBalanceRepository(db *gorm.DB) userbalance.Repository {
	return &userBalanceRepository{db}
}
