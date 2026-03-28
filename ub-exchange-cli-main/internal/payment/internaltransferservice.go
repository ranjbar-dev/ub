package payment

import (
	"strings"
)

const (
	InternalTransferStatusCreated    = "CREATED"
	InternalTransferStatusCompleted  = "COMPLETED"
	InternalTransferStatusInProgress = "IN_PROGRESS"
	InternalTransferStatusFailed     = "FAILED"
	InternalTransferStatusCanceled   = "CANCELED"
	InternalTransferStatusRejected   = "REJECTED"
	BalanceTypeHot                   = "HOT"
	BalanceTypeExternal              = "EXTERNAL"
	BalanceTypeCold                  = "COLD"
	BalanceTypeInternal              = "INTERNAL"
)

// InternalTransferService manages internal fund transfers between exchange wallets
// (e.g. hot, cold, external) and their status updates.
type InternalTransferService interface {
	// GetFromExternalInProgressTransfers returns all in-progress transfers from external wallets.
	GetFromExternalInProgressTransfers() []InternalTransfer
	// Update persists changes to an existing internal transfer record.
	Update(internalTransfer *InternalTransfer) error
	// UpdateStatus changes the status of an internal transfer by its ID.
	UpdateStatus(id int64, status string) error
}

type internalTransferService struct {
	internalTransferRepository InternalTransferRepository
}

func (s *internalTransferService) GetFromExternalInProgressTransfers() []InternalTransfer {
	return s.internalTransferRepository.GetFromExternalInProgressTransfers()
}

func (s *internalTransferService) Update(internalTransfer *InternalTransfer) error {
	return s.internalTransferRepository.Update(internalTransfer)
}

func (s *internalTransferService) UpdateStatus(id int64, status string) error {
	internalTransfer := &InternalTransfer{}
	err := s.internalTransferRepository.GetInternalTransferByID(id, internalTransfer)
	if err != nil {
		return err
	}
	internalTransfer.Status = strings.ToUpper(status)
	err = s.Update(internalTransfer)
	return err
}

func NewInternalTransferService(internalTransferRepository InternalTransferRepository) InternalTransferService {
	return &internalTransferService{
		internalTransferRepository: internalTransferRepository,
	}
}
