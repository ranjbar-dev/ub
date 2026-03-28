// Package user_test tests the forgot-password manager. Covers:
//   - GenerateForgotPasswordAndSendEmail: code generation, Redis storage, and async email dispatch
//   - IsCodeCorrect: validating a reset code against the stored Redis entry
//   - DeleteKey: removing a forgot-password entry from Redis after use
//
// Test data: testify mocks for RedisClient, CommunicationService, and Configs
// with Redis hash data containing expiredAt and code fields.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestForgotPasswordManager_GenerateForgotPasswordAndSendEmail(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HSet", mock.Anything, "forgot-password:1", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Once().Return(nil)
	redisClient.On("Expire", mock.Anything, "forgot-password:1", mock.Anything).Once().Return(true, nil)

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendUserForgotPasswordEmail", mock.Anything, mock.Anything).Once().Return()

	configs := new(mocks.Configs)
	configs.On("GetDomain").Once().Return("localhost")
	forgotPasswordManager := user.NewForgotPasswordManager(redisClient, communicationService, configs)
	u := user.User{
		ID: 1,
	}
	err := forgotPasswordManager.GenerateForgotPasswordAndSendEmail(u, "web", "127.0.0.1")
	assert.Nil(t, err)
	time.Sleep(20 * time.Millisecond)
	redisClient.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestForgotPasswordManager_IsCodeCorrect(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	expiredAt := strconv.FormatInt(time.Now().Add(1*time.Hour).Unix(), 10)
	data := map[string]string{
		"expiredAt": expiredAt,
		"code":      "someCode",
	}
	redisClient.On("HGetAll", mock.Anything, "forgot-password:1").Once().Return(data, nil)
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)

	forgotPasswordManager := user.NewForgotPasswordManager(redisClient, communicationService, configs)
	u := user.User{
		ID: 1,
	}
	isCodeCorrect := forgotPasswordManager.IsCodeCorrect(u, "someCode")
	assert.Equal(t, true, isCodeCorrect)

	redisClient.AssertExpectations(t)
}

func TestForgotPasswordManager_DeleteKey(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("Del", mock.Anything, "forgot-password:1").Once().Return(int64(1), nil)
	communicationService := new(mocks.CommunicationService)

	configs := new(mocks.Configs)

	forgotPasswordManager := user.NewForgotPasswordManager(redisClient, communicationService, configs)
	u := user.User{
		ID: 1,
	}
	err := forgotPasswordManager.DeleteKey(u)
	assert.Nil(t, err)

	redisClient.AssertExpectations(t)

}
