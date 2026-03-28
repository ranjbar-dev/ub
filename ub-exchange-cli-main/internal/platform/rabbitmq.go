package platform

import (
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMqClient provides message publishing to RabbitMQ for asynchronous
// communication with downstream services such as ub-communicator (email/SMS).
// The connection is lazily established and protected by a mutex for safe concurrent use.
type RabbitMqClient interface {
	// Enqueue publishes serialized data to a RabbitMQ topic exchange derived from
	// the routing key. For example, a key of "kline.BTC.USDT" publishes to the
	// "kline" exchange. Enqueue is a no-op in the test environment.
	Enqueue(key string, data []byte) error
}

type rabbitMqClient struct {
	connection  *amqp.Connection
	configs     Configs
	logger      Logger
	mutex       *sync.Mutex
	isConnected bool
}

func (r *rabbitMqClient) Enqueue(key string, data []byte) error {
	if r.configs.GetEnv() == EnvTest {
		return nil
	}
	ch, err := r.openChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	topic := r.getTopic(key)
	err = ch.ExchangeDeclare(topic, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return ch.Publish(topic, key, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         data,
	})
}

//for key like kline.something.otherthing the topic would be kline
func (r *rabbitMqClient) getTopic(key string) string {
	if pos := strings.IndexByte(key, '.'); pos > 0 {
		return key[:pos]
	}
	return key
}

func (r *rabbitMqClient) openChannel() (*amqp.Channel, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if err := r.connect(); err != nil {
		return nil, err
	}
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err

	}
	return ch, nil
}

func (r *rabbitMqClient) connect() error {
	if r.isConnected {
		return nil
	}
	rabbitMqDsn := r.configs.GetString("rabbitmq.dsn")
	conn, err := amqp.Dial(rabbitMqDsn)
	if err != nil {
		r.logger.Error2("error in dialing rabbitmq", err,
			zap.String("service", "rabbitMqClient"),
			zap.String("method", "NewRabbitMqClient"),
			zap.String("dsn", rabbitMqDsn),
		)
		return err
	}
	r.connection = conn
	r.isConnected = true
	return nil
}

func NewRabbitMqClient(c Configs, logger Logger) RabbitMqClient {
	return &rabbitMqClient{
		connection:  nil,
		configs:     c,
		logger:      logger,
		mutex:       &sync.Mutex{},
		isConnected: false,
	}
}
