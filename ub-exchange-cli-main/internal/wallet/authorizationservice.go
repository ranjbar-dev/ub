package wallet

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	authURI        = "/api/v1/auth/login"
	RedisKey       = "wallet:auth"
	ExpireDuration = time.Duration(time.Hour * 5)
)

type authRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponseBody struct {
	Token string `json:"token"`
}

// AuthorizationService manages authentication tokens for the external wallet
// microservice, caching them in Redis with automatic refresh on expiration.
type AuthorizationService interface {
	// GetToken returns a valid authentication token for the wallet service,
	// using a cached token from Redis or requesting a new one if expired.
	GetToken(ctx context.Context) (string, error)
}

type authorizationService struct {
	rc         platform.RedisClient
	logger     platform.Logger
	httpClient platform.HTTPClient
	host       string
	username   string
	password   string
}

func (as *authorizationService) GetToken(ctx context.Context) (string, error) {
	res, err := as.rc.HGetAll(ctx, RedisKey)
	if err != nil {
		return "", err
	}
	if isStillValid(res) {
		return res["token"], nil
	}

	httpResp, err := as.getTokenFromServer(ctx, as.username, as.password)
	if err != nil {
		return "", err
	}

	respBody := authResponseBody{}
	err = json.Unmarshal([]byte(httpResp), &respBody)
	if err != nil {
		return "", err
	}

	err = as.setTokenInRedis(ctx, respBody.Token)
	if err != nil {
		as.logger.Error2("can not set token in redis", err,
			zap.String("service", "walletAuthorizationService"),
			zap.String("method", "GetToken"),
		)
	}
	return respBody.Token, nil
}

func (as *authorizationService) getTokenFromServer(ctx context.Context, username string, password string) (string, error) {
	body := getAuthRequestBody(username, password)
	headers := getAuthRequestHeaders()
	authURL := as.host + authURI
	resp, _, statusCode, err := as.httpClient.HTTPPost(ctx, authURL, body, headers)
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("status code is not 200 it is %d", statusCode)
	}
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func isStillValid(data map[string]string) bool {
	if len(data) < 1 {
		return false
	}

	expiredAt, _ := data["expiredAt"]
	expiredAtTime, _ := time.Parse("2006-01-02 15:04:05", expiredAt)
	if expiredAtTime.Sub(time.Now()).Seconds() > 10 {
		return true
	}
	return false
}

func (as *authorizationService) setTokenInRedis(ctx context.Context, token string) error {
	now := time.Now()
	expiredAt := now.Add(ExpireDuration)
	err := as.rc.HSet(ctx, RedisKey, "token", token, "expiredAt", expiredAt.Format("2006-01-02 15:04:05"))
	return err
}

func getAuthRequestBody(username string, password string) authRequestBody {
	arb := authRequestBody{
		Username: username,
		Password: password,
	}

	return arb

}

func getAuthRequestHeaders() map[string]string {
	rh := map[string]string{
		"Content-Type": "application/json",
	}
	return rh
}

func NewAuthorizationService(rc platform.RedisClient, logger platform.Logger, httpClient platform.HTTPClient, configs platform.Configs) AuthorizationService {
	host := configs.GetString(HostEnvKey)
	username := configs.GetString(UsernameEnvKey)
	password := configs.GetString(PasswordEnvKey)
	return &authorizationService{rc, logger, httpClient, host, username, password}
}
