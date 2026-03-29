package platform

import (
	"crypto/tls"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMqClient manages a connection to RabbitMQ with mutex-protected lazy initialization.
// Connections are reused across GetChannel calls. If the connection drops,
// the next GetChannel call will attempt to reconnect.
type RabbitMqClient interface {
	// GetChannel returns an AMQP channel from the managed connection.
	// Thread-safe: uses a mutex to protect connection state.
	GetChannel() (*amqp.Channel, error)
}

type rabbitMqClient struct {
	connection  *amqp.Connection
	configs     Configs
	logger      Logger
	mutex       *sync.Mutex
	isConnected bool
}

func (r *rabbitMqClient) GetChannel() (*amqp.Channel, error) {
	// WHY: Mutex protects lazy connection creation. Multiple goroutines calling
	// GetChannel() concurrently during startup could create duplicate connections
	// without this synchronization.
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if err := r.connect(); err != nil {
		return nil, fmt.Errorf("rabbitmq connection failed: %w", err)
	}
	ch, err := r.connection.Channel()
	if err != nil {
		// Connection might have gone stale between connect() and Channel()
		// Mark as disconnected so next call attempts reconnection
		r.isConnected = false
		return nil, fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}
	return ch, nil
}

func (r *rabbitMqClient) connect() error {
	// Check if existing connection is still alive
	if r.isConnected && r.connection != nil && !r.connection.IsClosed() {
		return nil
	}

	// Connection is dead or doesn't exist — (re)connect
	r.isConnected = false
	rabbitMqDsn := r.configs.GetString("rabbitmq.dsn")

	var conn *amqp.Connection
	var err error
	if r.configs.GetBool("rabbitmq.tls") {
		tlsCfg := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		}
		conn, err = amqp.DialTLS(rabbitMqDsn, tlsCfg)
	} else {
		conn, err = amqp.Dial(rabbitMqDsn)
	}
	if err != nil {
		r.logger.Error("failed to dial rabbitmq", zap.Error(err))
		return fmt.Errorf("failed to dial rabbitmq: %w", err)
	}
	r.connection = conn
	r.isConnected = true
	return nil
}

// NewRabbitMqClient creates a RabbitMQ client. The actual connection
// is deferred until the first GetChannel() call.
func NewRabbitMqClient(c Configs, logger Logger) RabbitMqClient {

	return &rabbitMqClient{
		connection:  nil,
		configs:     c,
		logger:      logger,
		mutex:       &sync.Mutex{},
		isConnected: false,
	}
}
