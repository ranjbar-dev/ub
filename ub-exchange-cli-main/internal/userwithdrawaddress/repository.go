package userwithdrawaddress

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/user"
	"time"
)

type UserWithdrawAddress struct {
	ID         int64
	UserID     int
	User       user.User
	CoinID     int64         `gorm:"column:currency_id"`
	Coin       currency.Coin `gorm:"foreignKey:CoinID"`
	Address    string
	Label      sql.NullString
	IsDeleted  sql.NullBool
	IsFavorite sql.NullBool
	Network    sql.NullString
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (UserWithdrawAddress) TableName() string {
	return "user_withdraw_address"
}

// Repository provides data access for user withdrawal address records, including
// creation and lookups by address, ID, coin, and filters.
type Repository interface {
	// Create inserts a new withdrawal address for the user.
	Create(uwa *UserWithdrawAddress) error
	// GetUserWithdrawAddressesByAddress returns addresses matching the given user, coin, and address string.
	GetUserWithdrawAddressesByAddress(userID int, coinID int64, address string) []UserWithdrawAddress
	// GetUserWithdrawAddressByID looks up a single withdrawal address by its database ID.
	GetUserWithdrawAddressByID(id int64, uwa *UserWithdrawAddress) error
	// GetUserWithdrawAddressesByIds returns withdrawal addresses matching the given user and list of IDs.
	GetUserWithdrawAddressesByIds(userID int, ids []int64) []UserWithdrawAddress
	// GetUserWithdrawAddressesByCoinID returns all withdrawal addresses for a user and coin.
	GetUserWithdrawAddressesByCoinID(userID int, coinID int64) []UserWithdrawAddress
	// GetUserWithdrawAddresses returns withdrawal addresses matching the provided filters.
	GetUserWithdrawAddresses(filters GetWithdrawAddressesFilters) []UserWithdrawAddress
}
