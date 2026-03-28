package test

import (
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MqttAuthTests struct {
	*suite.Suite
	httpServer http.Handler
	db         *gorm.DB
	userActor  *userActor
}

func (t *MqttAuthTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.db = getDb()
	t.userActor = getUserActor()

}

func (t *MqttAuthTests) SetupTest() {}

func (t *MqttAuthTests) TearDownTest() {}

func (t *MqttAuthTests) TearDownSuite() {}

func (t *MqttAuthTests) TestLogin() {
	res := httptest.NewRecorder()
	reader := strings.NewReader("")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/login", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//token := "Bearer " + t.userActor.Token
	//req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)
}

func (t *MqttAuthTests) TestACL_PublicTopics() {
	publicTopics := []string{
		"main/trade/order-book/BTC-USDT",
		"main/trade/trade-book/BTC-USDT",
		"main/trade/ticker",
		"main/trade/kline/BTC-USDT",
	}
	for _, topic := range publicTopics {
		res := httptest.NewRecorder()
		reader := strings.NewReader("access=1&username=username&clientid=clientid&topic=" + topic)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/acl", reader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//token := "Bearer " + t.userActor.Token
		//req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		assert.Equal(t.T(), http.StatusOK, res.Code)
	}
}

func (t *MqttAuthTests) TestACL_PrivateTopic() {
	privateTopics := []string{
		"main/trade/user/userActorPrivateChannel/open-orders", // userActorPrivateChannel is set in main_test file
		"main/trade/user/userActorPrivateChannel/crypto-payments",
	}

	//successful scenarios
	for _, topic := range privateTopics {
		res := httptest.NewRecorder()
		reader := strings.NewReader("access=1&username=" + t.userActor.Token + "&clientid=clientid&topic=" + topic)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/acl", reader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//token := "Bearer " + t.userActor.Token
		//req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		assert.Equal(t.T(), http.StatusOK, res.Code)
	}

	//failed scenarios which we do not add token
	for _, topic := range privateTopics {
		res := httptest.NewRecorder()
		reader := strings.NewReader("access=1&username=username&clientid=clientid&topic=" + topic)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/acl", reader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//token := "Bearer " + t.userActor.Token
		//req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		assert.Equal(t.T(), http.StatusUnauthorized, res.Code)
	}

	//failed scenarios which with wrong channel name for user
	failedPrivateTopics := []string{
		"main/trade/user/failedChannel/open-orders",
		"main/trade/user/failedChannel/crypto-payments",
	}
	for _, topic := range failedPrivateTopics {
		res := httptest.NewRecorder()
		reader := strings.NewReader("access=1&username=" + t.userActor.Token + "&clientid=clientid&topic=" + topic)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/acl", reader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		//token := "Bearer " + t.userActor.Token
		//req.Header.Set("Authorization", token)
		t.httpServer.ServeHTTP(res, req)
		assert.Equal(t.T(), http.StatusUnauthorized, res.Code)
	}

}

func (t *MqttAuthTests) TestACL_SuperUser() {
	//successful scenario
	res := httptest.NewRecorder()
	params := url.Values{}
	params.Set("clientid", "mqtt_abbas2") // is set in config_test.yaml
	params.Set("username", "mqtt_abbas")  // is set in config_test.yaml
	params.Set("password", "mqtt_abbas")  // is set in config_test.yaml

	reader := strings.NewReader(params.Encode())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/superuser", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//token := "Bearer " + t.userActor.Token
	//req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusOK, res.Code)

	//failed scenario wrong username
	res = httptest.NewRecorder()
	params = url.Values{}
	params.Set("username", "wrongUsername")
	params.Set("password", "mqtt_abbas")  // is set in config_test.yaml
	params.Set("clientid", "mqtt_abbas2") // is set in config_test.yaml

	reader = strings.NewReader(params.Encode())
	req = httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/superuser", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//token := "Bearer " + t.userActor.Token
	//req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusUnauthorized, res.Code)

	//failed scenario wrong clientid
	res = httptest.NewRecorder()
	params = url.Values{}
	params.Set("username", "mqtt_abbas2") // is set in config_test.yaml
	params.Set("password", "mqtt_abbas")  // is set in config_test.yaml
	params.Set("clientid", "wrongone")

	reader = strings.NewReader(params.Encode())
	req = httptest.NewRequest(http.MethodPost, "/api/v1/emqtt/superuser", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//token := "Bearer " + t.userActor.Token
	//req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)
	assert.Equal(t.T(), http.StatusUnauthorized, res.Code)

}

func TestMqttAuth(t *testing.T) {
	suite.Run(t, &MqttAuthTests{
		Suite: new(suite.Suite),
	})
}
