package transaction

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID          int64
	UserID      int
	CoinID      int64 `gorm:"column:currency_id"`
	OrderID     sql.NullInt64
	Type        string
	MoneyAmount sql.NullString
	Amount      sql.NullString
	CoinName    string        `gorm:"column:money_currency"`
	PaymentID   sql.NullInt64 `gorm:"column:crypto_payment_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository is a placeholder interface for transaction persistence operations.
type Repository interface {
}
