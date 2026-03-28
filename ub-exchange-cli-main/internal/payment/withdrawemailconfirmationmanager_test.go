// Package payment_test tests the WithdrawEmailConfirmationManager which
// handles email-based confirmation codes for withdrawal requests. Covers:
//   - CheckCode: Redis key does not exist, code has expired, code does not match
//   - IsAllowedToSendEmail: Redis key does not exist (allowed), code has expired
//     (allowed), and rate-limited when less than one minute since last send
//   - CreateAndSendWithdrawEmailConfirmationCode: stores confirmation data in
//     Redis with expiration and dispatches the email via CommunicationService
//   - RemoveConfirmationCodeFromRedis: deletes the confirmation key after use
//
// Test data: mocked RedisClient and CommunicationService; Redis hash keys
// follow the pattern "withdraw-confirmation:{userID}" with fields for userId,
// amount, coin, address, expiredAt, and code.
package payment_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/payment"
	"exchange-go/internal/user"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWithdrawEmailConfirmationManager_CheckCode_RedisKeyDoesNotExist(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(map[string]string{}, redis.Nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	code := "123456"
	isValid, err := service.CheckCode(u, code)

	assert.Nil(t, err)
	assert.Equal(t, false, isValid)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_CheckCode_HasExpired(t *testing.T) {
	now := time.Now().Unix()
	expiredAt := now - 2*60 //2 minutes ago
	expiredAtString := strconv.FormatInt(expiredAt, 10)
	rc := new(mocks.RedisClient)
	data := map[string]string{
		"userId":    "1",
		"amount":    "0.01",
		"coin":      "BTC",
		"address":   "someAddress",
		"expiredAt": expiredAtString,
		"code":      "123456",
	}
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(data, nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	code := "123456"
	isValid, err := service.CheckCode(u, code)

	assert.Nil(t, err)
	assert.Equal(t, false, isValid)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_CheckCode_NoEqualCode(t *testing.T) {
	now := time.Now().Unix()
	expiredAt := now + 30*60 //30 minutes later
	expiredAtString := strconv.FormatInt(expiredAt, 10)
	rc := new(mocks.RedisClient)
	data := map[string]string{
		"userId":    "1",
		"amount":    "0.01",
		"coin":      "BTC",
		"address":   "someAddress",
		"expiredAt": expiredAtString,
		"code":      "123457",
	}
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(data, nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	code := "123456"
	isValid, err := service.CheckCode(u, code)

	assert.Nil(t, err)
	assert.Equal(t, false, isValid)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_IsAllowedToSendEmail_RedisKeyDoesNotExist(t *testing.T) {
	rc := new(mocks.RedisClient)
	data := map[string]string{}
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(data, redis.Nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	isAllowed, err := service.IsAllowedToSendEmail(u, "BTC", "0.01", "someAddress")

	assert.Nil(t, err)
	assert.Equal(t, true, isAllowed)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_IsAllowedToSendEmail_HasExpired(t *testing.T) {
	now := time.Now().Unix()
	expiredAt := now - 2*60 //2 minutes ago
	expiredAtString := strconv.FormatInt(expiredAt, 10)

	rc := new(mocks.RedisClient)
	data := map[string]string{
		"userId":    "1",
		"amount":    "0.01",
		"coin":      "BTC",
		"address":   "someAddress",
		"expiredAt": expiredAtString,
		"code":      "123457",
	}
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(data, nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	isAllowed, err := service.IsAllowedToSendEmail(u, "BTC", "0.01", "someAddress")

	assert.Nil(t, err)
	assert.Equal(t, true, isAllowed)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_IsAllowedToSendEmail_LessThanOneMinute(t *testing.T) {
	now := time.Now().Unix()
	expiredAt := now + 3*60*60 //20 seconds ago
	expiredAtString := strconv.FormatInt(expiredAt, 10)

	rc := new(mocks.RedisClient)
	data := map[string]string{
		"userId":    "1",
		"amount":    "0.01",
		"coin":      "BTC",
		"address":   "someAddress",
		"expiredAt": expiredAtString,
		"code":      "123457",
	}
	rc.On("HGetAll", mock.Anything, "withdraw-confirmation:1").Once().Return(data, nil)
	communicationService := new(mocks.CommunicationService)
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	isAllowed, err := service.IsAllowedToSendEmail(u, "BTC", "0.01", "someAddress")
	assert.Nil(t, err)
	assert.Equal(t, false, isAllowed)
	rc.AssertExpectations(t)
}

func TestWithdrawEmailConfirmationManager_CreateAndSendWithdrawEmailConfirmationCode(t *testing.T) {
	rc := new(mocks.RedisClient)

	//the commented code should work but since map does not have order every time different data sent to redis
	rc.On("HSet", mock.Anything, "withdraw-confirmation:1", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Once().Return(nil)

	//rc.On("HSet", mock.Anything, "withdraw-confirmation:1", "userId", "1", "amount", "0.01",
	//	"coin", "BTC", "address", "someAddress", "expiredAt", mock.Anything, "code",
	//	mock.Anything).Once().Return( nil)

	rc.On("Expire", mock.Anything, "withdraw-confirmation:1", mock.Anything).Once().Return(true, nil)

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendWithdrawConfirmationEmail", mock.Anything, "BTC", "0.01", "someAddress", mock.Anything).Once().Return()
	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	err := service.CreateAndSendWithdrawEmailConfirmationCode(u, "BTC", "0.01", "someAddress")

	assert.Nil(t, err)
	rc.AssertExpectations(t)
	communicationService.AssertExpectations(t)

}

func TestWithdrawEmailConfirmationManager_RemoveConfirmationCodeFromRedis(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("Del", mock.Anything, "withdraw-confirmation:1").Once().Return(int64(1), nil)

	communicationService := new(mocks.CommunicationService)

	service := payment.NewWithdrawEmailConfirmationManager(rc, communicationService)
	u := user.User{
		ID: 1,
	}
	err := service.RemoveConfirmationCodeFromRedis(u, "BTC", "0.01", "someAddress")

	assert.Nil(t, err)
	rc.AssertExpectations(t)
}
