// Package user_test tests the user config service. Covers:
//   - GetUserConfig: fetching a user's configuration by user ID
//
// Test data: testify mock for UserConfigRepository returning a Config
// entity populated via Run callback.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigService_GetUserConfig(t *testing.T) {
	userConfigRepo := new(mocks.UserConfigRepository)
	userConfigRepo.On("GetUserConfigByUserID", 1, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		userConfig := args.Get(1).(*user.Config)
		userConfig.ID = 1
	})
	configService := user.NewUserConfigService(userConfigRepo)
	config, err := configService.GetUserConfig(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), config.ID)
	userConfigRepo.AssertExpectations(t)
}
