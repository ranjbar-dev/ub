package user

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"exchange-go/internal/communication"
	"exchange-go/internal/platform"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	ForgotPasswordHashPrefix        = "forgot-password:"
	ForgotPasswordExpirationSeconds = 3 * 60 * 60 //3 hour
)

// ForgotPasswordManager handles the password reset flow including code generation,
// validation, and email delivery.
type ForgotPasswordManager interface {
	// GenerateForgotPasswordAndSendEmail creates a password reset code, stores it in Redis,
	// and sends a reset link to the user's email.
	GenerateForgotPasswordAndSendEmail(u User, device string, ip string) error
	// IsCodeCorrect validates whether the provided reset code matches the stored code
	// and has not expired.
	IsCodeCorrect(u User, code string) bool
	// DeleteKey removes the password reset code from Redis after it has been used.
	DeleteKey(u User) error
}

type forgotPasswordManager struct {
	redisClient          platform.RedisClient
	communicationService communication.Service
	configs              platform.Configs
	logger               platform.Logger
}

func (s *forgotPasswordManager) GenerateForgotPasswordAndSendEmail(u User, device string, ip string) error {
	data := make(map[string]string, 6)

	userIDString := strconv.Itoa(u.ID)
	data["userId"] = userIDString

	code := s.generateNewForgotPasswordCode()
	data["code"] = code

	data["mobileCode"] = ""

	now := time.Now().Unix()
	expiredAt := now + ForgotPasswordExpirationSeconds
	expiredAtString := strconv.FormatInt(expiredAt, 10)
	data["expiredAt"] = expiredAtString
	data["isCodeConfirmed"] = "0" // equal to false

	link := s.generateLinkForForgotPassword(u, code)

	//if device == userdevice.DeviceWeb {

	//}
	//else {
	//	min := 111111
	//	max := 999999
	//	rand.Seed(time.Now().UnixNano())
	//	mobileCode := rand.Int63n(int64(max-min+1)) + int64(min)
	//	code := strconv.FormatInt(mobileCode, 10)
	//	data["mobileCode"] = code
	//}

	var values []interface{}
	for k, v := range data {
		values = append(values, k, v)
	}

	key := s.getForgotPasswordKey(u.ID)
	err := s.redisClient.HSet(context.Background(), key, values...)
	if err != nil {
		return err
	}

	//errors ignored on purpose because the data being set in redis would be enough for us
	_, _ = s.redisClient.Expire(context.Background(), key, time.Duration(ForgotPasswordExpirationSeconds*time.Second))

	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}

	params := communication.ForgotPasswordEmailParams{
		Link:          link,
		CurrentIP:     ip,
		CurrentDevice: device,
		RequestDate:   time.Now().Format("2006-01-02 15:04:05"),
	}
	platform.SafeGo(s.logger, "user.SendUserForgotPasswordEmail", func() {
		s.communicationService.SendUserForgotPasswordEmail(cu, params)
	})
	return nil
}

func (s *forgotPasswordManager) IsCodeCorrect(u User, code string) bool {
	key := s.getForgotPasswordKey(u.ID)
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

	if data["code"] == code {
		return true
	}

	return false
}

func (s *forgotPasswordManager) DeleteKey(u User) error {
	key := s.getForgotPasswordKey(u.ID)
	_, err := s.redisClient.Del(context.Background(), key)
	return err

}

func (s *forgotPasswordManager) generateLinkForForgotPassword(u User, code string) string {
	domain := s.configs.GetDomain()
	link := ""
	if strings.HasSuffix(domain, "/") {
		link = domain + "auth/forgot-password/update?code=" + code + "&email=" + u.Email
	} else {
		link = domain + "/auth/forgot-password/update?code=" + code + "&email=" + u.Email
	}
	return link
}

func (s *forgotPasswordManager) generateNewForgotPasswordCode() string {
	uuidCode := uuid.NewString()
	sha := sha1.New()
	_, _ = io.WriteString(sha, uuidCode)
	shaResult := hex.EncodeToString(sha.Sum(nil))

	h := md5.New()
	_, _ = io.WriteString(h, shaResult)
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

func (s *forgotPasswordManager) getForgotPasswordKey(userID int) string {
	userIDString := strconv.Itoa(userID)
	return ForgotPasswordHashPrefix + userIDString
}

func NewForgotPasswordManager(redisClient platform.RedisClient, communicationService communication.Service,
	configs platform.Configs, logger platform.Logger) ForgotPasswordManager {
	return &forgotPasswordManager{
		redisClient:          redisClient,
		communicationService: communicationService,
		configs:              configs,
		logger:               logger,
	}

}
