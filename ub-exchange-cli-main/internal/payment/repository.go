package payment

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/user"
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID                int64
	UserID            int
	User              user.User     `gorm:"foreignKey:UserID"`
	CoinID            int64         `gorm:"column:currency_id"`
	Coin              currency.Coin `gorm:"foreignKey:CoinID"`
	Type              string
	Status            string
	Code              string `gorm:"column:money_currency"`
	FromAddress       sql.NullString
	ToAddress         sql.NullString
	TxID              sql.NullString
	AdminStatus       sql.NullString
	WithdrawType      sql.NullString
	BlockchainNetwork sql.NullString
	Amount            sql.NullString
	FeeAmount         sql.NullString
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (Payment) TableName() string {
	return "crypto_payments"
}

// Repository provides data access for payment (deposit/withdrawal) records with
// filtering and transaction support.
type Repository interface {
	// GetUserPayments returns payments matching the specified filters.
	GetUserPayments(filters GetPaymentFilters) []Payment
	// GetPaymentDetailByID returns detailed query fields for a single payment.
	GetPaymentDetailByID(paymentID int64) (DetailQueryFields, error)
	// GetPaymentByID looks up a payment by its database ID.
	GetPaymentByID(id int64, payment *Payment) error
	// GetPaymentByIDUsingTx looks up a payment by ID within an existing database transaction.
	GetPaymentByIDUsingTx(tx *gorm.DB, id int64, payment *Payment) error
	// GetExtraInfoByPaymentID retrieves extra metadata for a payment.
	GetExtraInfoByPaymentID(paymentID int64, extraInfo *ExtraInfo) error
	// GetExtraInfoByPaymentIDUsingTx retrieves extra metadata within a database transaction.
	GetExtraInfoByPaymentIDUsingTx(tx *gorm.DB, paymentID int64, extraInfo *ExtraInfo) error
	// GetInProgressWithdrawalsInExternalExchange returns withdrawals that are still
	// in progress on the external exchange and need status updates.
	GetInProgressWithdrawalsInExternalExchange() []ExternalWithdrawalUpdateDataNeeded
	// GetPaymentByCoinIDAndTxIDAndTypeUsingTx finds a payment by coin, transaction ID,
	// and type within a database transaction, used for duplicate detection.
	GetPaymentByCoinIDAndTxIDAndTypeUsingTx(tx *gorm.DB, coinID int64, txID string, txType string, payment *Payment) error
}

type ExtraInfo struct {
	ID                           int64
	PaymentID                    int64 `gorm:"column:crypto_payment_id"`
	LastHandledID                sql.NullInt64
	Tag                          sql.NullString
	NetworkFee                   sql.NullString
	UserMessage                  sql.NullString
	RejectionReason              sql.NullString
	IP                           sql.NullString
	AutoTransfer                 sql.NullBool
	AutoExchangeOrderID          sql.NullInt64
	AutoExchangeFailureType      sql.NullString
	AutoExchangeFailureReason    sql.NullString
	Price                        sql.NullString
	BtcPrice                     sql.NullString
	ExternalExchangeID           sql.NullInt64
	ExternalExchangeWithdrawID   sql.NullString
	ExternalExchangeWithdrawInfo sql.NullString
}

func (ExtraInfo) TableName() string {
	return "crypto_payment_extra_info"
}

type InternalTransfer struct {
	ID              int64
	FromBalanceID   int64         `gorm:"column:crypto_balance_from_id"`
	ToBalanceID     sql.NullInt64 `gorm:"column:crypto_balance_to_id"`
	Amount          string
	TxID            sql.NullString
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Metadata        sql.NullString
	Network         string `gorm:"column:blockchain_network"`
	ToCustomAddress sql.NullString
}

func (InternalTransfer) TableName() string {
	return "crypto_internal_transfer"
}

// InternalTransferRepository provides data access for internal fund transfer records
// between exchange wallets (hot, cold, external).
type InternalTransferRepository interface {
	// GetFromExternalInProgressTransfers returns all internal transfers from external
	// wallets that are currently in progress.
	GetFromExternalInProgressTransfers() []InternalTransfer
	// GetInternalTransferByID looks up an internal transfer by its database ID.
	GetInternalTransferByID(id int64, internalTransfer *InternalTransfer) error
	// Update persists changes to an existing internal transfer record.
	Update(internalTransfer *InternalTransfer) error
}
