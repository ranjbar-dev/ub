package auth

import (
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"strings"

	"go.uber.org/zap"
)

const (
	UserPrivateTopicPrefix    = "user"
	UserOpenOrderTopicPostfix = "open-orders"
	UserPaymentsTopicPostfix  = "crypto-payments"
)

type MqttLoginParams struct {
}

type MqttACLParams struct {
	Access   int    `form:"access"`
	Username string `form:"username"`
	ClientID string `form:"clientid"`
	Topic    string `form:"topic"`
}

type MqttSuperUserParams struct {
	Username string `form:"username"`
	Password string `form:"password"`
	ClientID string `form:"clientid"`
	//Ipaddress string `form:"ipaddress"`
}

// MqttAuthService handles authentication and authorization for the EMQX MQTT broker,
// including connection login, topic-level access control, and super-user verification.
type MqttAuthService interface {
	// Login authenticates an MQTT client connection using JWT credentials.
	Login(params MqttLoginParams) (apiResponse response.APIResponse, statusCode int)
	// ACL checks topic-level publish and subscribe permissions for an MQTT client.
	ACL(params MqttACLParams) (apiResponse response.APIResponse, statusCode int)
	// SuperUser checks whether the MQTT client has admin-level super-user access.
	SuperUser(params MqttSuperUserParams) (apiResponse response.APIResponse, statusCode int)
}

type mqttAuthService struct {
	authService Service
	configs     platform.Configs
	logger      platform.Logger
	publicTopic []string
}

func (s *mqttAuthService) Login(params MqttLoginParams) (apiResponse response.APIResponse, statusCode int) {
	//everyone can login so we return success
	return response.Success(nil, "")
}

func (s *mqttAuthService) ACL(params MqttACLParams) (apiResponse response.APIResponse, statusCode int) {
	if params.Access == 1 {
		if s.isTopicAllowed(params.Topic, params.Username, params.ClientID) {
			return response.Success(nil, "")
		}
		return response.Unauthorized(nil, "")
	}
	publisherUserName := s.configs.GetString("mqtt.username")
	publisherClientID := s.configs.GetString("mqtt.clientid")

	if params.Username == publisherUserName && strings.Contains(params.ClientID, publisherClientID) {
		return response.Success(nil, "")
	}
	return response.Unauthorized(nil, "")
}

func (s *mqttAuthService) isTopicAllowed(topic, username, clientID string) bool {
	topicParts := strings.Split(topic, "/")
	if len(topicParts) < 3 { //our topics are generated from 3 parts or more
		return false
	}

	if topicParts[0] == "main" && topicParts[1] == "trade" {
		for _, t := range s.publicTopic {
			if t == topicParts[2] {
				return true
			}
		}

		if topicParts[2] == UserPrivateTopicPrefix {
			if len(topicParts) > 4 {
				if topicParts[4] != UserOpenOrderTopicPostfix && topicParts[4] != UserPaymentsTopicPostfix {
					return false
				}
				if s.doesChannelBelongsToUser(username, topicParts[3]) {
					return true
				}
			}
		}
	}
	return false

}

func (s *mqttAuthService) doesChannelBelongsToUser(username string, channel string) bool {
	//user name of mqtt client is the jwt token we have for our auth
	u, err := s.authService.GetUser(username)
	if err != nil {
		s.logger.Error2("error in getting user", err,
			zap.String("service", "MqttAuthService"),
			zap.String("method", "doesChannelBelongsToUser"),
			zap.String("usermane", username),
		)
		return false
	}

	if u.PrivateChannelName == channel {
		return true
	}

	return false
}

func (s *mqttAuthService) SuperUser(params MqttSuperUserParams) (apiResponse response.APIResponse, statusCode int) {
	publisherUserName := s.configs.GetString("mqtt.username")
	//publisherPassword := s.configs.GetString("mqtt.password")
	publisherClientID := s.configs.GetString("mqtt.clientid")

	if params.Username == publisherUserName && strings.Contains(params.ClientID, publisherClientID) {
		return response.Success(nil, "")
	}
	return response.Unauthorized(nil, "")
}

func NewMqttAuthService(authService Service, configs platform.Configs, logger platform.Logger) MqttAuthService {
	s := &mqttAuthService{
		authService: authService,
		configs:     configs,
		logger:      logger,
	}

	s.publicTopic = []string{
		"order-book",
		"trade-book",
		"ticker",
		"kline",
	}

	return s
}
