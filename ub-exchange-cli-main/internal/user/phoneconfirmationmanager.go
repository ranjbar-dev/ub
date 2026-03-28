package user

import (
	"context"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	PhoneConfirmationHashPrefix = "phone-confirmation:"
	SmsExpirationSeconds        = 3 * 60 * 60 //3 hour
)

// PhoneConfirmationManager handles the SMS-based phone verification flow including
// rate limiting, code generation, validation, and cleanup.
type PhoneConfirmationManager interface {
	// IsAllowedToSendSms checks whether the user is allowed to request a new SMS code
	// based on rate limiting rules (minimum 60 seconds between requests).
	IsAllowedToSendSms(u User) bool
	// GeneratePhoneConfirmationCodeAndSendSms generates a 6-digit verification code,
	// stores it in Redis, and sends it to the provided phone number via SMS.
	GeneratePhoneConfirmationCodeAndSendSms(u User, phone string) error
	// IsCodeCorrect validates whether the provided code and phone number match the
	// stored verification data and have not expired.
	IsCodeCorrect(u User, phone string, code string) bool
	// DeleteKey removes the phone confirmation code from Redis after successful verification.
	DeleteKey(u User) error
}

type phoneConfirmationManager struct {
	redisClient          platform.RedisClient
	communicationService communication.Service
}

func (s *phoneConfirmationManager) IsAllowedToSendSms(u User) bool {
	key := s.getPhoneConfirmationKey(u.ID)
	ctx := context.Background()
	data, err := s.redisClient.HGetAll(ctx, key)
	if err == redis.Nil {
		return true
	}

	nowTimestamp := time.Now().Unix()

	expiredAtTimestamp, _ := strconv.ParseInt(data["expiredAt"], 10, 64)
	if expiredAtTimestamp < nowTimestamp {
		return true
	}

	createdAtTimestamp := expiredAtTimestamp - SmsExpirationSeconds

	if nowTimestamp-createdAtTimestamp > 60 {
		return true
	}

	return false
}

func (s *phoneConfirmationManager) GeneratePhoneConfirmationCodeAndSendSms(u User, phone string) error {
	//generate code
	//todo although this works but for security reasons better to use crypto/rand package instead of math/rand package
	min := 111111
	max := 999999
	code := rand.Int63n(int64(max-min+1)) + int64(min)
	codeString := strconv.FormatInt(code, 10)

	data := make(map[string]string, 6)

	now := time.Now().Unix()
	expiredAt := now + SmsExpirationSeconds
	expiredAtString := strconv.FormatInt(expiredAt, 10)

	userIDString := strconv.Itoa(u.ID)

	data["userId"] = userIDString
	data["code"] = codeString
	data["expiredAt"] = expiredAtString
	data["phone"] = phone

	var values []interface{}
	for k, v := range data {
		values = append(values, k, v)
	}

	ctx := context.Background()
	key := s.getPhoneConfirmationKey(u.ID)

	err := s.redisClient.HSet(ctx, key, values...)
	if err != nil {
		return err
	}

	//errors ignored on purpose because the data being set in redis would be enough for us
	_, _ = s.redisClient.Expire(ctx, key, time.Duration(SmsExpirationSeconds*time.Second))

	cu := communication.CommunicatingUser{
		Email: "",
		Phone: phone,
	}
	go s.communicationService.SendUserPhoneConfirmationSms(cu, codeString)
	return nil
}

func (s *phoneConfirmationManager) getPhoneConfirmationKey(userID int) string {
	userIDString := strconv.Itoa(userID)
	return PhoneConfirmationHashPrefix + userIDString
}

func (s *phoneConfirmationManager) IsCodeCorrect(u User, phone string, code string) bool {
	key := s.getPhoneConfirmationKey(u.ID)
	ctx := context.Background()
	data, err := s.redisClient.HGetAll(ctx, key)
	if err == redis.Nil {
		return false
	}

	nowTimestamp := time.Now().Unix()

	expiredAtTimestamp, _ := strconv.ParseInt(data["expiredAt"], 10, 64)
	if expiredAtTimestamp < nowTimestamp {
		return false
	}

	if data["phone"] == phone && data["code"] == code {
		return true
	}

	return false
}

func (s *phoneConfirmationManager) DeleteKey(u User) error {
	key := s.getPhoneConfirmationKey(u.ID)
	_, err := s.redisClient.Del(context.Background(), key)
	return err
}

func NewPhoneConfirmationManager(redisClient platform.RedisClient, communicationService communication.Service) PhoneConfirmationManager {
	return &phoneConfirmationManager{
		redisClient:          redisClient,
		communicationService: communicationService,
	}
}
