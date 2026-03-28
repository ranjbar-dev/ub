// Package user_test tests the user level service. Covers:
//   - GetLevelByID: fetching a level by its database ID
//   - RecalculateUserLevel: determining the correct level based on exchange volume and KYC tier
//   - GetLevelByCode: fetching a level by its code identifier
//
// Test data: testify mocks for UserLevelRepository with four levels
// defined by MinExchangeVolume, MaxExchangeVolume, and MinKycLevel thresholds.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLevelService_GetLevelById(t *testing.T) {
	userLevelRepo := new(mocks.UserLevelRepository)
	userLevelRepo.On("GetLevelByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		userLevel := args.Get(1).(*user.Level)
		userLevel.ID = 1
	})
	userLevelService := user.NewUserLevelService(userLevelRepo)
	level, err := userLevelService.GetLevelByID(1)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), level.ID)
	userLevelRepo.AssertExpectations(t)
}

func TestLevelService_RecalculateUserLevel(t *testing.T) {
	userLevelRepo := new(mocks.UserLevelRepository)
	allLevels := []user.Level{
		{
			ID:                1,
			MinExchangeVolume: 0.1,
			MaxExchangeVolume: 1.0,
			MinKycLevel:       0,
		},
		{
			ID:                2,
			MinExchangeVolume: 1.0,
			MaxExchangeVolume: 2.0,
			MinKycLevel:       1,
		},
		{
			ID:                3,
			MinExchangeVolume: 2.0,
			MaxExchangeVolume: 3.0,
			MinKycLevel:       1,
		},
		{
			ID:                4,
			MinExchangeVolume: 3.0,
			MaxExchangeVolume: 4.0,
			MinKycLevel:       1,
		},
	}
	userLevelRepo.On("GetAllLevels").Once().Return(allLevels)
	userLevelService := user.NewUserLevelService(userLevelRepo)
	u := user.User{
		Kyc: 1,
	}

	level := userLevelService.RecalculateUserLevel(u, "3.0")
	assert.Equal(t, int64(4), level.ID)
	userLevelRepo.AssertExpectations(t)
}

func TestLevelService_GetLevelByCode(t *testing.T) {
	userLevelRepo := new(mocks.UserLevelRepository)
	userLevelRepo.On("GetLevelByCode", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		userLevel := args.Get(1).(*user.Level)
		userLevel.ID = 1
	})
	userLevelService := user.NewUserLevelService(userLevelRepo)
	level, err := userLevelService.GetLevelByCode(1)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), level.ID)
	userLevelRepo.AssertExpectations(t)
}
