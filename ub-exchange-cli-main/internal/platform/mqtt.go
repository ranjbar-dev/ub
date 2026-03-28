package platform

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MqttClient provides MQTT publishing for real-time data push to client
// applications via the EMQX broker. Messages are published to topics that
// clients subscribe to for live tickers, order updates, and trade notifications.
type MqttClient interface {
	// Publish sends a message to the specified MQTT topic with the given QoS level
	// and retention flag. The connection is lazily established if not already active.
	// Publishing is a no-op in the test environment.
	Publish(topic string, qos byte, retained bool, payload interface{}) (token mqtt.Token)
}

type client struct {
	mqttClient mqtt.Client
	logger     Logger
	env        string
}

func (c *client) Publish(topic string, qos byte, retained bool, payload interface{}) (token mqtt.Token) {

	if c.env == EnvTest {
		return
	}
	if c.isConnected() {
		t := c.mqttClient.Publish(topic, qos, retained, payload)
		t.Wait()
		if t.Error() != nil {
			c.logger.Warn("error in mqtt publish",
				zap.Error(t.Error()),
				zap.String("service", "mqttClient"),
				zap.String("method", "connect"),
			)
		}

		return t
	}

	c.connect()
	return c.mqttClient.Publish(topic, qos, retained, payload)
}

func (c *client) isConnected() bool {
	return c.mqttClient.IsConnected()
}

func (c *client) connect() {
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		c.logger.Warn("can not connect to mqtt ",
			zap.Error(token.Error()),
			zap.String("service", "mqttClient"),
			zap.String("method", "connect"),
		)
	}

}

func NewMqttClient(configs Configs, logger Logger) MqttClient {
	dsn := configs.GetString("mqtt.dsn")
	uuid := uuid.NewString()
	clientID := configs.GetString("mqtt.clientid") + uuid
	username := configs.GetString("mqtt.username")
	password := configs.GetString("mqtt.password")
	env := configs.GetEnv()
	opts := mqtt.NewClientOptions()
	opts.AddBroker(dsn)
	opts.ClientID = clientID
	opts.Username = username
	opts.Password = password
	opts.Order = false
	opts.OnConnectionLost = func(i mqtt.Client, e error) {
		logger.Warn("error in new MqttClient",
			zap.Error(e),
			zap.String("service", "mqttClient"),
			zap.String("method", "NewMqttClient"),
		)
	}
	mqttClient := mqtt.NewClient(opts)
	client := &client{mqttClient, logger, env}
	if env != EnvTest {
		client.connect()
	}
	return client
}
