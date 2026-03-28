// Package user_test tests the phone confirmation manager. Covers:
//   - IsAllowedToSendSms: no prior record, expired code, and active (non-expired) code
//   - GeneratePhoneConfirmationCodeAndSendSms: code generation, Redis storage, and async SMS dispatch
//   - IsCodeCorrect: valid code, expired code, mismatched code, and mismatched phone number
//   - DeleteKey: removing a phone confirmation entry from Redis
//
// Test data: testify mocks for RedisClient and CommunicationService with
// Redis hash data containing userId, code, expiredAt, and phone fields.
package user_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPhoneConfirmationManager_IsAllowedToSendSms(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:1").Once().Return(map[string]string{}, redis.Nil)

	nowTimestamp := time.Now().Unix()
	expiredAtTimestamp := nowTimestamp - (20 * 60) //20 minutes ago
	expiredAtTimestampString := strconv.FormatInt(expiredAtTimestamp, 10)
	data := map[string]string{
		"userId":    "2",
		"code":      "",
		"expiredAt": expiredAtTimestampString,
		"phone":     "",
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:2").Once().Return(data, nil)

	expiredAtTimestamp2 := nowTimestamp + user.SmsExpirationSeconds
	expiredAtTimestampString2 := strconv.FormatInt(expiredAtTimestamp2, 10)

	data = map[string]string{
		"userId":    "3",
		"code":      "",
		"expiredAt": expiredAtTimestampString2,
		"phone":     "",
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:3").Once().Return(data, nil)

	communicationService := new(mocks.CommunicationService)
	phoneConfirmationManager := user.NewPhoneConfirmationManager(redisClient, communicationService)
	u1 := user.User{
		ID: 1,
	}
	isAllowed := phoneConfirmationManager.IsAllowedToSendSms(u1)
	assert.Equal(t, true, isAllowed)

	//expiredAt is before now
	u2 := user.User{
		ID: 2,
	}
	isAllowed = phoneConfirmationManager.IsAllowedToSendSms(u2)
	assert.Equal(t, true, isAllowed)

	//less than one minute
	u3 := user.User{
		ID: 3,
	}
	isAllowed = phoneConfirmationManager.IsAllowedToSendSms(u3)
	assert.Equal(t, false, isAllowed)

	redisClient.AssertExpectations(t)
}

func TestPhoneConfirmationManager_GeneratePhoneConfirmationCodeAndSendSms(t *testing.T) {
	redisClient := new(mocks.RedisClient)

	//the commented code should work but since map does not have order every time different data sent to redis
	redisClient.On("HSet", mock.Anything, "phone-confirmation:1", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	redisClient.On("Expire", mock.Anything, "phone-confirmation:1", mock.Anything).Once().Return(true, nil)

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendUserPhoneConfirmationSms", mock.Anything, mock.Anything).Once().Return()
	phoneConfirmationManager := user.NewPhoneConfirmationManager(redisClient, communicationService)
	u := user.User{
		ID: 1,
	}
	err := phoneConfirmationManager.GeneratePhoneConfirmationCodeAndSendSms(u, "+989121234567")
	assert.Nil(t, err)

	redisClient.AssertExpectations(t)
	time.Sleep(20 * time.Millisecond)
	communicationService.AssertExpectations(t)

}

func TestPhoneConfirmationManager_IsCodeCorrect(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	expiredAtTimestamp := time.Now().Unix() + (20 * 60) //20 minutes ago
	expiredAtTimestampString := strconv.FormatInt(expiredAtTimestamp, 10)

	data := map[string]string{
		"userId":    "1",
		"code":      "123456",
		"expiredAt": expiredAtTimestampString,
		"phone":     "+989121234567",
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:1").Once().Return(data, nil)

	expiredAtTimestamp2 := time.Now().Unix() - (20 * 60) //20 minutes ago
	expiredAtTimestampString2 := strconv.FormatInt(expiredAtTimestamp2, 10)
	data2 := map[string]string{
		"userId":    "2",
		"code":      "123456",
		"expiredAt": expiredAtTimestampString2,
		"phone":     "+989121234567",
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:2").Once().Return(data2, nil)

	data3 := map[string]string{
		"userId":    "3",
		"code":      "1234567",
		"expiredAt": expiredAtTimestampString,
		"phone":     "+989121234567",
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:3").Once().Return(data3, nil)

	data4 := map[string]string{
		"userId":    "3",
		"code":      "123456",
		"expiredAt": expiredAtTimestampString,
		"phone":     "+989121234568", //different last digit
	}
	redisClient.On("HGetAll", mock.Anything, "phone-confirmation:4").Once().Return(data4, nil)

	communicationService := new(mocks.CommunicationService)
	phoneConfirmationManager := user.NewPhoneConfirmationManager(redisClient, communicationService)
	u := user.User{
		ID: 1,
	}
	isCorrect := phoneConfirmationManager.IsCodeCorrect(u, "+989121234567", "123456")
	assert.True(t, isCorrect)

	//expiredAt has passed
	u2 := user.User{
		ID: 2,
	}
	isCorrect = phoneConfirmationManager.IsCodeCorrect(u2, "+989121234567", "123456")
	assert.False(t, isCorrect)

	//code is not the same as redis
	u3 := user.User{
		ID: 3,
	}
	isCorrect = phoneConfirmationManager.IsCodeCorrect(u3, "+989121234567", "123456")
	assert.False(t, isCorrect)

	//phone is not the same as redis
	u4 := user.User{
		ID: 4,
	}
	isCorrect = phoneConfirmationManager.IsCodeCorrect(u4, "+989121234567", "123456")
	assert.False(t, isCorrect)

	redisClient.AssertExpectations(t)
}

func TestPhoneConfirmationManager_DeleteKey(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("Del", mock.Anything, "phone-confirmation:1").Once().Return(int64(1), nil)

	communicationService := new(mocks.CommunicationService)
	phoneConfirmationManager := user.NewPhoneConfirmationManager(redisClient, communicationService)
	u := user.User{
		ID: 1,
	}
	err := phoneConfirmationManager.DeleteKey(u)
	assert.Nil(t, err)

	redisClient.AssertExpectations(t)
}
