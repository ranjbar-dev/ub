package repository

import (
	"exchange-go/internal/payment"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type internalTransferRepository struct {
	db *gorm.DB
}

func (r *internalTransferRepository) GetFromExternalInProgressTransfers() []payment.InternalTransfer {
	var result []payment.InternalTransfer
	r.db.Table("crypto_internal_transfer").
		Joins("left join crypto_balance as from_crypto_balance on crypto_internal_transfer.crypto_balance_from_id = from_crypto_balance.id").
		Joins("left join crypto_balance as to_crypto_balance on crypto_internal_transfer.crypto_balance_to_id = to_crypto_balance.id").
		Where("crypto_internal_transfer.status = ?", payment.InternalTransferStatusInProgress).
		Where("from_crypto_balance.type = ?", payment.BalanceTypeExternal).
		Order("crypto_internal_transfer.id asc").Find(&result)
	return result
}

func (r *internalTransferRepository) GetInternalTransferByID(id int64, internalTransfer *payment.InternalTransfer) error {
	return r.db.Where(payment.InternalTransfer{ID: id}).First(internalTransfer).Error
}

func (r *internalTransferRepository) Update(internalTransfer *payment.InternalTransfer) error {
	return r.db.Omit(clause.Associations).Save(internalTransfer).Error
}

func NewInternalTransferRepository(db *gorm.DB) payment.InternalTransferRepository {
	return &internalTransferRepository{db}
}
