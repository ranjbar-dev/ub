package userbalance

import (
	"database/sql"
	"encoding/json"
	"exchange-go/internal/currency"
	"exchange-go/internal/user"
	"time"

	"gorm.io/gorm"
)

type OtherAddress struct {
	Code    string `json:"code"`
	Address string `json:"address"`
}

type UserBalance struct {
	ID               int64
	UserID           int
	User             user.User     `gorm:"foreignKey:UserID"`
	CoinID           int64         `gorm:"column:currency_id"`
	Coin             currency.Coin `gorm:"foreignKey:CoinID"`
	FrozenBalance    string
	BalanceCoin      string `gorm:"column:balance_currency"`
	Status           string
	Address          sql.NullString
	AutoExchangeCoin sql.NullString `gorm:"column:auto_exchange_code"`
	Amount           string
	FrozenAmount     string
	OtherAddresses   sql.NullString
}

func (ub UserBalance) GetOtherAddresses() ([]OtherAddress, error) {
	var otherAddresses []OtherAddress
	if !ub.OtherAddresses.Valid {
		return otherAddresses, nil
	}
	err := json.Unmarshal([]byte(ub.OtherAddresses.String), &otherAddresses)
	return otherAddresses, err
}

// Repository provides data access for user balance records with filtering,
// transaction support, and address-based lookups.
type Repository interface {
	// GetBalanceOfUserByCoinID retrieves a user's balance for a specific coin.
	GetBalanceOfUserByCoinID(userID int, coinID int64, balance *UserBalance) error
	// GetBalanceOfUserByCoinIDUsingTx retrieves a user's balance within a database transaction.
	GetBalanceOfUserByCoinIDUsingTx(tx *gorm.DB, userID int, coinID int64, balance *UserBalance) error
	// GetBalancesOfUsersForCoins returns balances for multiple users and coins at once.
	GetBalancesOfUsersForCoins(userIds []int, coinIds []int64) []UserBalance
	// GetBalancesOfUsersForCoinsUsingTx returns balances for multiple users/coins within a transaction.
	GetBalancesOfUsersForCoinsUsingTx(tx *gorm.DB, userIds []int, coinIds []int64) []UserBalance
	// GetUserBalancesForCoins returns a single user's balances for the specified coins.
	GetUserBalancesForCoins(userID int, coinIds []int64) []UserBalance
	// GetUserAllBalances returns all balances for a user, optionally filtered.
	GetUserAllBalances(userID int, filters AllBalancesFilters) []UserBalance
	// GetBalancesWithoutAddresses returns balances that have no deposit address assigned.
	GetBalancesWithoutAddresses(filters map[string]interface{}) []UserBalance
	// GetUserBalanceByCoinAndAddressUsingTx looks up a balance by coin and deposit address
	// within a database transaction.
	GetUserBalanceByCoinAndAddressUsingTx(tx *gorm.DB, coinID int64, address string, ub *UserBalance) error
	// GetUserBalanceByIDUsingTx looks up a balance by ID within a database transaction.
	GetUserBalanceByIDUsingTx(tx *gorm.DB, id int64, ub *UserBalance) error
}

type UserWalletBalance struct {
	ID            int64
	UserID        int
	CoinID        int64         `gorm:"column:currency_id"`
	NetworkCoinID sql.NullInt64 `gorm:"column:network_currency_id"`
	Balance       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (UserWalletBalance) TableName() string {
	return "user_wallet_balance"
}

// UserWalletBalanceRepository provides data access for per-network wallet balances.
type UserWalletBalanceRepository interface {
	// FindUserWalletBalance looks up a user's wallet balance for a specific coin on a specific network.
	FindUserWalletBalance(userID int, coinID int64, networkCoinID int64, userBalance *UserWalletBalance) error
}
