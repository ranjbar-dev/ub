package platform

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

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
	if r.isConnected && r.connection != nil && !r.connection.IsClosed() {
		return nil
	}
	r.isConnected = false
	rabbitMqDsn := r.configs.GetString("rabbitmq.dsn")

	var conn *amqp.Connection
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			r.logger.Warn("retrying rabbitmq connection", zap.Int("attempt", attempt+1), zap.Duration("backoff", backoff))
			time.Sleep(backoff)
		}
		cfg := amqp.Config{
			Dial: amqp.DefaultDial(10 * time.Second),
		}
		if r.configs.GetBool("rabbitmq.tls") {
			cfg.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			}
		}
		conn, err = amqp.DialConfig(rabbitMqDsn, cfg)
		if err == nil {
			break
		}
		r.logger.Error("rabbitmq dial failed", zap.Int("attempt", attempt+1), zap.Error(err))
	}
	if err != nil {
		return fmt.Errorf("failed to dial rabbitmq after %d attempts: %w", maxRetries, err)
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
