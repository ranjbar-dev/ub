// Package payment_test tests the InternalTransferService which manages
// internal fund transfers between user balances. Covers:
//   - GetFromExternalInProgressTransfers: retrieves in-progress transfers
//     originating from external sources and verifies correct record count and IDs
//   - Update: persists changes to an InternalTransfer record via the repository
//
// Test data: mocked InternalTransferRepository returning InternalTransfer
// structs with minimal field population.
package payment_test

import (
	"database/sql"
	"exchange-go/internal/mocks"
	"exchange-go/internal/payment"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInternalTransferService_GetFromExternalInProgressTransfers(t *testing.T) {
	internalTransferRepository := new(mocks.InternalTransferRepository)
	data := []payment.InternalTransfer{
		{
			ID:              1,
			FromBalanceID:   0,
			ToBalanceID:     sql.NullInt64{},
			Amount:          "",
			TxID:            sql.NullString{},
			Status:          "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			Metadata:        sql.NullString{},
			Network:         "",
			ToCustomAddress: sql.NullString{},
		},
		{
			ID:              2,
			FromBalanceID:   0,
			ToBalanceID:     sql.NullInt64{},
			Amount:          "",
			TxID:            sql.NullString{},
			Status:          "",
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			Metadata:        sql.NullString{},
			Network:         "",
			ToCustomAddress: sql.NullString{},
		},
	}
	internalTransferRepository.On("GetFromExternalInProgressTransfers").Once().Return(data)

	service := payment.NewInternalTransferService(internalTransferRepository)
	internalTransfers := service.GetFromExternalInProgressTransfers()
	assert.Equal(t, 2, len(internalTransfers))
	assert.Equal(t, int64(1), internalTransfers[0].ID)
	assert.Equal(t, int64(2), internalTransfers[1].ID)
	internalTransferRepository.AssertExpectations(t)

}

func TestInternalTransferService_Update(t *testing.T) {
	internalTransferRepository := new(mocks.InternalTransferRepository)
	internalTransferRepository.On("Update", mock.Anything).Once().Return(nil)

	service := payment.NewInternalTransferService(internalTransferRepository)

	internalTransfer := &payment.InternalTransfer{}
	err := service.Update(internalTransfer)
	assert.Nil(t, err)
	internalTransferRepository.AssertExpectations(t)
}
