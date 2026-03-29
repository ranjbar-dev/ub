// Package di provides a dependency injection container with lazy singleton
// initialization. Each service is created on first access and cached for
// subsequent calls. The container wires together all application dependencies
// from configuration through to the consumer service.
//
// Initialization order (lazy, triggered by GetConsumer()):
//   config → logger → rabbitMQ → mongoDB → repository → mail/sms services → messaging → pool → consumer
package di

import (
	"fmt"
	"ub-communicator/config"
	"ub-communicator/pkg/consumer"
	"ub-communicator/pkg/messaging"
	"ub-communicator/pkg/platform"
	"ub-communicator/pkg/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

// Container provides access to application services.
// All services are lazily initialized on first access.
type Container interface {
	GetConsumer() consumer.Service
}

type container struct {
	consumer          consumer.Service
	pool              consumer.Pool
	messagingService  messaging.Service
	mailService       messaging.MailService
	smsService        messaging.SmsService
	mailerClient      platform.MailerClient
	configs           platform.Configs
	logger            platform.Logger
	httpClient        platform.HttpClient
	db                *mongo.Client
	messageRepository messaging.Repository
	rabbitMq          platform.RabbitMqClient
}

func (c *container) GetConsumer() consumer.Service {
	if c.consumer == nil {
		rc := c.getRabbitMq()
		ms := c.getMessagingService()
		logger := c.getLogger()
		configs := c.getConfigs()
		pool := c.getPool()
		c.consumer = consumer.NewConsumerService(rc, ms, pool, logger, configs)
	}
	return c.consumer
}

func (c *container) getPool() consumer.Pool {
	if c.pool == nil {
		ms := c.getMessagingService()
		c.pool = consumer.NewPool(ms)
	}
	return c.pool
}

func (c *container) getMessagingService() messaging.Service {
	if c.messagingService == nil {
		messageRepository := c.getMessageRepository()
		mailService := c.getMailService()
		smsService := c.getSmsService()
		logger := c.getLogger()
		c.messagingService = messaging.NewMessagingService(messageRepository, mailService, smsService, logger)
	}
	return c.messagingService
}

func (c *container) getMailService() messaging.MailService {
	if c.mailService == nil {
		mc := c.getMailerClient()
		c.mailService = messaging.NewMailService(mc)
	}
	return c.mailService
}

func (c *container) getSmsService() messaging.SmsService {
	if c.smsService == nil {
		httpClient := c.getHttpClient()
		configs := c.getConfigs()
		c.smsService = messaging.NewSmsService(httpClient, configs)
	}
	return c.smsService
}

func (c *container) getMailerClient() platform.MailerClient {
	if c.mailerClient == nil {
		configs := c.getConfigs()
		logger := c.getLogger()
		mailerClient := platform.NewMailerClient(configs, logger)
		if mailerClient == nil {
			panic("failed to instantiate mailer client: check mailer_broker config value")
		}
		c.mailerClient = mailerClient
	}
	return c.mailerClient
}

func (c *container) getConfigs() platform.Configs {
	if c.configs == nil {
		viper, err := config.SetConfigs()
		if err != nil {
			// Cannot return error from this getter without cascading changes
			// Log error and return nil, which will be caught by NewContainer validation
			return nil
		}
		c.configs = platform.NewConfigs(viper)
	}
	return c.configs
}

func (c *container) getLogger() platform.Logger {
	if c.logger == nil {
		configs := c.getConfigs()
		c.logger = platform.NewLogger(configs)
	}
	return c.logger
}

func (c *container) getHttpClient() platform.HttpClient {
	if c.httpClient == nil {
		c.httpClient = platform.NewHttpClient()
	}
	return c.httpClient
}

func (c *container) getDb() (*mongo.Client, error) {
	if c.db == nil {
		configs := c.getConfigs()
		db, err := platform.NewDbClient(configs)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		c.db = db
	}
	return c.db, nil
}

func (c *container) getMessageRepository() messaging.Repository {
	if c.messageRepository == nil {
		configs := c.getConfigs()
		db, err := c.getDb()
		if err != nil {
			panic(fmt.Sprintf("failed to get database connection for message repository: %v", err))
		}
		repo := repository.NewMessageRepository(db, configs)
		if repo == nil {
			panic("NewMessageRepository returned nil: check database configuration")
		}
		c.messageRepository = repo
	}
	return c.messageRepository
}

func (c *container) getRabbitMq() platform.RabbitMqClient {
	if c.rabbitMq == nil {
		configs := c.getConfigs()
		logger := c.getLogger()
		c.rabbitMq = platform.NewRabbitMqClient(configs, logger)
	}
	return c.rabbitMq
}

// NewContainer creates and returns the application DI container.
// Returns (Container, error) so initialization failures are propagated
// instead of crashing the process via panic().
func NewContainer() (Container, error) {
	c := &container{}

	// Eagerly validate critical dependencies to fail fast with clear errors
	if c.getConfigs() == nil {
		return nil, fmt.Errorf("failed to initialize configuration")
	}
	if c.getLogger() == nil {
		return nil, fmt.Errorf("failed to initialize logger")
	}

	return c, nil
}
