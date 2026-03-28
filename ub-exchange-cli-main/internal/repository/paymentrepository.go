package repository

import (
	"exchange-go/internal/payment"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentRepository struct {
	db *gorm.DB
}

func (pr *paymentRepository) GetUserPayments(filters payment.GetPaymentFilters) []payment.Payment {
	var payments []payment.Payment
	q := pr.db.Joins("Coin").
		Where("crypto_payments.user_id =?", filters.UserID)

	if filters.Type != "" {
		q.Where("crypto_payments.type =?", filters.Type)
	}

	if filters.Coin != "" {
		q.Where("Coin.code = ?", filters.Coin)
	}

	if filters.StartDate != "" {
		q.Where("crypto_payments.created_at >= ?", filters.StartDate)
	}

	if filters.EndDate != "" {
		q.Where("crypto_payments.created_at <= ?", filters.EndDate)
	}

	q.Order("crypto_payments.updated_at desc").Offset(int(filters.Page) * filters.PageSize).Limit(filters.PageSize)

	q.Find(&payments)
	return payments
}

func (pr *paymentRepository) GetPaymentDetailByID(paymentID int64) (payment.DetailQueryFields, error) {
	result := payment.DetailQueryFields{}

	err := pr.db.Table("crypto_payments as cp").
		Joins("join currencies as c on cp.currency_id = c.id").
		Joins("left join crypto_payment_extra_info as cpei on cp.id = cpei.crypto_payment_id").
		Select("cp.id as ID,cp.to_address as Address, cp.user_id as UserID,cp.tx_id as TxID,cp.blockchain_network as Network,c.code as Code,cpei.rejection_reason as RejectionReason").
		Where("cp.id =?", paymentID).Limit(1).Scan(&result).Error

	return result, err

}

func (pr *paymentRepository) GetPaymentByID(id int64, p *payment.Payment) error {
	return pr.db.Joins("User").Where(payment.Payment{ID: id}).First(p).Error
}

func (pr *paymentRepository) GetPaymentByIDUsingTx(tx *gorm.DB, id int64, p *payment.Payment) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(payment.Payment{ID: id}).First(p).Error
}

func (pr *paymentRepository) GetExtraInfoByPaymentID(paymentID int64, ei *payment.ExtraInfo) error {
	return pr.db.Where(payment.ExtraInfo{PaymentID: paymentID}).First(ei).Error
}

func (pr *paymentRepository) GetExtraInfoByPaymentIDUsingTx(tx *gorm.DB, paymentID int64, ei *payment.ExtraInfo) error {
	return tx.Where(payment.ExtraInfo{PaymentID: paymentID}).First(ei).Error
}

func (pr *paymentRepository) GetInProgressWithdrawalsInExternalExchange() []payment.ExternalWithdrawalUpdateDataNeeded {
	var results []payment.ExternalWithdrawalUpdateDataNeeded
	pr.db.Table("crypto_payments as cp").
		Joins("left join crypto_payment_extra_info as cpei on cp.id = cpei.crypto_payment_id").
		Select(""+
			"cp.id as PaymentID,"+
			"cp.updated_at as UpdatedAt,"+
			"cpei.id as PaymentExtraInfoID,"+
			"cpei.external_exchange_withdraw_id as ExternalExchangeWithdrawID").
		Where("cp.status = ? and cpei.external_exchange_withdraw_id IS NOT NULL", payment.StatusInProgress).Limit(100).Scan(&results)

	return results
}
func (pr *paymentRepository) GetPaymentByCoinIDAndTxIDAndTypeUsingTx(tx *gorm.DB, coinID int64, txID string, txType string, p *payment.Payment) error {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("currency_id = ? AND tx_id = ? AND type = ? ", coinID, txID, txType).First(p).Error
}

func NewPaymentRepository(db *gorm.DB) payment.Repository {
	return &paymentRepository{db}
}
