// Package user_test tests the login history service. Covers:
//   - CreateLoginHistory: persisting a new login record with device, IP, email, and type
//   - GetLastLoginHistoryByUserID: retrieving the most recent login entry for a user
//
// Test data: testify mocks for LoginHistoryRepository with sample
// LoginHistory records containing SQL nullable fields.
package user_test

import (
	"database/sql"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHistoryService_CreateLoginHistory(t *testing.T) {
	loginHistoryRepository := new(mocks.LoginHistoryRepository)
	loginHistoryRepository.On("Create", mock.Anything).Once().Return(nil)
	loginHistoryService := user.NewLoginHistoryService(loginHistoryRepository)
	lh := &user.LoginHistory{
		UserID:   sql.NullInt64{Int64: 1, Valid: true},
		Device:   sql.NullString{String: "WEB", Valid: true},
		IP:       sql.NullString{String: "127.0.0.1", Valid: true},
		Email:    sql.NullString{String: "test@test.com", Valid: true},
		Type:     "SUCCESSFUL",
		Password: sql.NullString{String: "", Valid: false},
	}
	err := loginHistoryService.CreateLoginHistory(lh)
	assert.Nil(t, err)
	loginHistoryRepository.AssertExpectations(t)
}

func TestLoginHistoryService_GetLastLoginHistoryByUserId(t *testing.T) {
	loginHistoryRepository := new(mocks.LoginHistoryRepository)
	loginHistoryRepository.On("GetLastLoginHistoryByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		lh := args.Get(1).(*user.LoginHistory)
		lh.ID = 1

	})

	loginHistoryService := user.NewLoginHistoryService(loginHistoryRepository)
	lh := &user.LoginHistory{
		UserID:   sql.NullInt64{Int64: 1, Valid: true},
		Device:   sql.NullString{String: "WEB", Valid: true},
		IP:       sql.NullString{String: "127.0.0.1", Valid: true},
		Email:    sql.NullString{String: "test@test.com", Valid: true},
		Type:     "SUCCESSFUL",
		Password: sql.NullString{String: "", Valid: false},
	}
	err := loginHistoryService.GetLastLoginHistoryByUserID(1, lh)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), lh.ID)
	loginHistoryRepository.AssertExpectations(t)
}
