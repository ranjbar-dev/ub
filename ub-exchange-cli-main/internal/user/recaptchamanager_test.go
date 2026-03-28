// Package user_test tests the reCAPTCHA manager. Covers:
//   - CheckRecaptcha: verifying a successful reCAPTCHA response from the Google API
//
// Test data: testify mocks for HttpClient (returning a JSON success response),
// Configs (secret key, environment), Logger, and UbCaptchaManager.
package user_test

import (
	"bytes"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecaptchaManager_CheckRecaptcha(t *testing.T) {

	httpClient := new(mocks.HttpClient)
	body := []byte("{" +
		"\"success\":true," +
		"\"score\":0.1," +
		"\"action\":\"\"," +
		"\"challenge_ts\":\"124564\"" +
		"}")

	r := bytes.NewReader(body)

	responseBody := io.NopCloser(r)

	response := http.Response{
		StatusCode: http.StatusOK,
		Body:       responseBody,
	}
	//response.Body =
	httpClient.On("PostForm", mock.Anything, mock.Anything, mock.Anything).Once().Return(&response, nil)

	configs := new(mocks.Configs)
	configs.On("GetRecaptchaSecretKey").Once().Return("secretKey")
	configs.On("GetEnv").Once().Return("prod")

	logger := new(mocks.Logger)
	ubCaptchaManager := new(mocks.UbCaptchaManager)
	recaptchaManager := user.NewRecaptchaManager(httpClient, configs, logger, ubCaptchaManager)
	recaptchaString := "somerecaptcha"
	device := "WEB"
	ip := "127.0.0.1"
	isSuccessful, err := recaptchaManager.CheckRecaptcha(recaptchaString, device, ip)
	assert.Nil(t, err)
	assert.Equal(t, true, isSuccessful)
	configs.AssertExpectations(t)
	httpClient.AssertExpectations(t)
	logger.AssertExpectations(t)
	ubCaptchaManager.AssertExpectations(t)
}
