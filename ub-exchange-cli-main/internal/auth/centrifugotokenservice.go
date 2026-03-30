package auth

import (
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CentrifugoTokenService generates JWT tokens for Centrifugo client authentication.
type CentrifugoTokenService interface {
	// GenerateConnectionToken creates a connection JWT with sub (user ID) and exp claims.
	GenerateConnectionToken(userID int) (apiResponse response.APIResponse, statusCode int)
	// GenerateSubscriptionToken creates a subscription JWT for a private channel.
	GenerateSubscriptionToken(userID int, channel string) (apiResponse response.APIResponse, statusCode int)
}

type centrifugoTokenService struct {
	hmacSecret string
	configs    platform.Configs
}

func (s *centrifugoTokenService) GenerateConnectionToken(userID int) (apiResponse response.APIResponse, statusCode int) {
	token, err := s.generateToken(userID, "")
	if err != nil {
		return response.Error(fmt.Sprintf("failed to generate token: %v", err), 500, nil)
	}
	return response.Success(map[string]string{"token": token}, "")
}

func (s *centrifugoTokenService) GenerateSubscriptionToken(userID int, channel string) (apiResponse response.APIResponse, statusCode int) {
	token, err := s.generateToken(userID, channel)
	if err != nil {
		return response.Error(fmt.Sprintf("failed to generate token: %v", err), 500, nil)
	}
	return response.Success(map[string]string{"token": token}, "")
}

func (s *centrifugoTokenService) generateToken(userID int, channel string) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	if channel != "" {
		claims["channel"] = channel
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.hmacSecret))
}

func NewCentrifugoTokenService(configs platform.Configs) CentrifugoTokenService {
	return &centrifugoTokenService{
		hmacSecret: configs.GetString("centrifugo.token_hmac_secret"),
		configs:    configs,
	}
}
