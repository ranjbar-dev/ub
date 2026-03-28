// Package wallet_test tests the wallet authorization service. Covers:
//   - Acquiring authentication tokens from the wallet API when no cached token exists
//   - Caching tokens in Redis after successful retrieval
//
// Test data: mock Redis client, HTTP client, and config provider with
// JSON token response fixtures.
package wallet_test

import (
	"context"
	"exchange-go/internal/mocks"
	"exchange-go/internal/wallet"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthorizationService_GetToken(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("HGetAll", mock.Anything, mock.Anything).Once().Return(map[string]string{}, nil)
	rc.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	logger := new(mocks.Logger)
	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"token\": \"token\"" +
		"}")
	httpClient.On("HTTPPost", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(body, http.Header{}, http.StatusOK, nil)

	configs := new(mocks.Configs)
	configs.On("GetString", mock.Anything).Times(3).Return("")

	s := wallet.NewAuthorizationService(rc, logger, httpClient, configs)
	ctx := context.Background()
	token, err := s.GetToken(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "token", token)
}
