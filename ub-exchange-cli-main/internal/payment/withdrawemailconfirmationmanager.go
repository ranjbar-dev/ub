package payment

import (
	"context"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	WithdrawEmailConfirmationHashPrefix = "withdraw-confirmation:"
	ExpirationSeconds                   = 3 * 60 * 60 //3 hour
)

// WithdrawEmailConfirmationManager handles email-based confirmation codes for
// withdrawal security, including code generation, validation, and expiration.
type WithdrawEmailConfirmationManager interface {
	// CheckCode verifies whether the provided confirmation code is valid for the user.
	CheckCode(u user.User, code string) (bool, error)
	// IsAllowedToSendEmail checks whether the user is allowed to request a new
	// confirmation email for the given withdrawal parameters.
	IsAllowedToSendEmail(u user.User, coin string, amount string, address string) (bool, error)
	// CreateAndSendWithdrawEmailConfirmationCode generates a new confirmation code,
	// stores it in Redis, and sends it to the user via email.
	CreateAndSendWithdrawEmailConfirmationCode(u user.User, coin string, amount string, address string) error
	// RemoveConfirmationCodeFromRedis deletes the confirmation code from Redis after
	// a successful withdrawal or expiration.
	RemoveConfirmationCodeFromRedis(u user.User, coin string, amount string, address string) error
}

type withdrawEmailConfirmationManager struct {
	redisClient          platform.RedisClient
	communicationService communication.Service
}

func (m *withdrawEmailConfirmationManager) CheckCode(u user.User, code string) (bool, error) {
	key := m.getWithdrawEmailConfirmationKey(u.ID)
	ctx := context.Background()
	data, err := m.redisClient.HGetAll(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}

	if err == redis.Nil {
		return false, nil
	}

	expiredAtTimestamp, _ := strconv.ParseInt(data["expiredAt"], 10, 64)
	nowTimestamp := time.Now().Unix()
	if expiredAtTimestamp < nowTimestamp {
		return false, nil
	}

	if data["code"] == code {
		return true, nil
	}

	return false, nil
}

func (m *withdrawEmailConfirmationManager) getWithdrawEmailConfirmationKey(userID int) string {
	userIDString := strconv.Itoa(userID)
	return WithdrawEmailConfirmationHashPrefix + userIDString
}

func (m *withdrawEmailConfirmationManager) IsAllowedToSendEmail(u user.User, coin string, amount string, address string) (bool, error) {
	key := m.getWithdrawEmailConfirmationKey(u.ID)
	ctx := context.Background()

	data, err := m.redisClient.HGetAll(ctx, key)
	if err == redis.Nil {
		return true, nil
	}

	nowTimestamp := time.Now().Unix()
	expiredAtTimestamp, _ := strconv.ParseInt(data["expiredAt"], 10, 64)
	if expiredAtTimestamp < nowTimestamp {
		return true, nil
	}

	createdAtTimestamp := expiredAtTimestamp - ExpirationSeconds

	//only one email per minute is possible
	if nowTimestamp-createdAtTimestamp > 60 {
		return true, nil
	}

	return false, nil
}

func (m *withdrawEmailConfirmationManager) CreateAndSendWithdrawEmailConfirmationCode(u user.User, coin string, amount string, address string) error {
	//generate code
	//todo although this works but for security reasons better to use crypto/rand package instead of math/rand package
	min := 111111
	max := 999999
	code := rand.Int63n(int64(max-min+1)) + int64(min)
	codeString := strconv.FormatInt(code, 10)

	data := make(map[string]string, 6)

	now := time.Now().Unix()
	expiredAt := now + ExpirationSeconds

	expiredAtString := strconv.FormatInt(expiredAt, 10)

	userIDString := strconv.Itoa(u.ID)

	data["userId"] = userIDString
	data["amount"] = amount
	data["coin"] = coin
	data["address"] = address
	data["expiredAt"] = expiredAtString
	data["code"] = codeString

	var values []interface{}
	for k, v := range data {
		values = append(values, k, v)
	}

	ctx := context.Background()
	key := m.getWithdrawEmailConfirmationKey(u.ID)

	err := m.redisClient.HSet(ctx, key, values...)
	if err != nil {
		return err
	}

	//errors ignored on purpose because the data being set in redis would be enough for us
	_, _ = m.redisClient.Expire(ctx, key, time.Duration(ExpirationSeconds*time.Second))

	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}
	m.communicationService.SendWithdrawConfirmationEmail(cu, coin, amount, address, codeString)
	return nil
}

func (m *withdrawEmailConfirmationManager) RemoveConfirmationCodeFromRedis(u user.User, coin string, amount string, address string) error {
	ctx := context.Background()
	key := m.getWithdrawEmailConfirmationKey(u.ID)
	_, err := m.redisClient.Del(ctx, key)
	return err
}

func NewWithdrawEmailConfirmationManager(redisClient platform.RedisClient, communicationService communication.Service) WithdrawEmailConfirmationManager {
	return &withdrawEmailConfirmationManager{
		redisClient:          redisClient,
		communicationService: communicationService,
	}
}
