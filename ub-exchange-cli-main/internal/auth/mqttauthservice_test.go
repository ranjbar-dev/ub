// Package auth_test tests the MQTT authentication service. Covers:
//   - MQTT login always succeeds (open authentication)
//   - ACL checks for public topics (order-book, trade-book, ticker, kline)
//   - ACL checks for user private topics (success when channel matches, fail on mismatch)
//   - SuperUser validation (success, wrong username, wrong client ID)
//
// Test data: testify mocks for AuthService, Configs, and Logger with
// configured MQTT credentials and user private channel names.
package auth_test

import (
	"exchange-go/internal/auth"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMqttAuthService_Login(t *testing.T) {
	authService := new(mocks.AuthService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)

	params := auth.MqttLoginParams{}
	res, statusCode := mqttAuthService.Login(params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
}

func TestMqttAuthService_ACL_PublicTopics(t *testing.T) {

	authService := new(mocks.AuthService)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)

	publicTopics := []string{
		"order-book",
		"trade-book",
		"ticker",
		"kline",
	}

	for _, topic := range publicTopics {
		params := auth.MqttACLParams{
			Access:   1,
			Username: "test",
			ClientID: "test",
			Topic:    "main/trade/" + topic,
		}
		res, statusCode := mqttAuthService.ACL(params)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, true, res.Status)
		assert.Equal(t, "", res.Message)
	}
}

func TestMqttAuthService_ACL_UserPrivateTopic_Successful(t *testing.T) {
	u := &user.User{
		PrivateChannelName: "someprivatechannel",
	}
	authService := new(mocks.AuthService)
	authService.On("GetUser", "testToken").Twice().Return(u, nil)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)
	params := auth.MqttACLParams{
		Access:   1,
		Username: "testToken",
		ClientID: "test",
		Topic:    "main/trade/user/someprivatechannel/open-orders",
	}

	res, statusCode := mqttAuthService.ACL(params)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	params = auth.MqttACLParams{
		Access:   1,
		Username: "testToken",
		ClientID: "test",
		Topic:    "main/trade/user/someprivatechannel/crypto-payments",
	}

	res, statusCode = mqttAuthService.ACL(params)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	authService.AssertExpectations(t)
}

func TestMqttAuthService_ACL_UserPrivateTopic_Fail(t *testing.T) {
	u := &user.User{
		PrivateChannelName: "otherprivatechannel",
	}
	authService := new(mocks.AuthService)
	authService.On("GetUser", "testToken").Twice().Return(u, nil)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)
	params := auth.MqttACLParams{
		Access:   1,
		Username: "testToken",
		ClientID: "test",
		Topic:    "main/trade/user/someprivatechannel/open-orders",
	}

	res, statusCode := mqttAuthService.ACL(params)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "", res.Message)

	params = auth.MqttACLParams{
		Access:   1,
		Username: "testToken",
		ClientID: "test",
		Topic:    "main/trade/user/someprivatechannel/crypto-payments",
	}

	res, statusCode = mqttAuthService.ACL(params)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "", res.Message)

	authService.AssertExpectations(t)
}

func TestMqttAuthService_SuperUser_Successful(t *testing.T) {
	authService := new(mocks.AuthService)
	configs := new(mocks.Configs)
	configs.On("GetString", "mqtt.username").Once().Return("testPublisher")
	configs.On("GetString", "mqtt.clientid").Once().Return("testPublisher")

	logger := new(mocks.Logger)

	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)
	params := auth.MqttSuperUserParams{
		Username: "testPublisher",
		Password: "",
		ClientID: "testPublisher2",
	}

	res, statusCode := mqttAuthService.SuperUser(params)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	configs.AssertExpectations(t)

}

func TestMqttAuthService_SuperUser_FailWrongUsername(t *testing.T) {
	authService := new(mocks.AuthService)
	configs := new(mocks.Configs)
	configs.On("GetString", "mqtt.username").Once().Return("testPublisher")
	configs.On("GetString", "mqtt.clientid").Once().Return("testPublisher")

	logger := new(mocks.Logger)

	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)
	params := auth.MqttSuperUserParams{
		Username: "testPublisher2",
		Password: "",
		ClientID: "testPublisher2",
	}

	res, statusCode := mqttAuthService.SuperUser(params)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "", res.Message)

	configs.AssertExpectations(t)

}

func TestMqttAuthService_SuperUser_FailWrongClientId(t *testing.T) {
	authService := new(mocks.AuthService)
	configs := new(mocks.Configs)
	configs.On("GetString", "mqtt.username").Once().Return("testPublisher")
	configs.On("GetString", "mqtt.clientid").Once().Return("somethingElse")

	logger := new(mocks.Logger)

	mqttAuthService := auth.NewMqttAuthService(authService, configs, logger)
	params := auth.MqttSuperUserParams{
		Username: "testPublisher",
		Password: "",
		ClientID: "testPublisher2",
	}

	res, statusCode := mqttAuthService.SuperUser(params)

	assert.Equal(t, http.StatusUnauthorized, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "", res.Message)

	configs.AssertExpectations(t)

}
