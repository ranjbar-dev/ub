package communication

import (
	"exchange-go/internal/platform"

	"go.uber.org/zap"
)

const (
	rkKLineCreated = "livedata.event.kline-created"

	rkMessageSend = "messages.command.send"
)

// QueueManager publishes messages to RabbitMQ queues for asynchronous processing
// by downstream consumers.
type QueueManager interface {
	// PublishEmailOrSms sends notification data to the email/SMS processing queue for delivery.
	PublishEmailOrSms(data []byte)
	// PublishKline sends kline (candlestick) data to the kline processing queue.
	PublishKline(data []byte)
}

type queueManager struct {
	rabbitmqClient platform.RabbitMqClient
	logger         platform.Logger
}

func (q *queueManager) PublishEmailOrSms(data []byte) {
	err := q.rabbitmqClient.Enqueue(rkMessageSend, data)
	if err != nil {
		q.logger.Error2("can not enqueue kline", err,
			zap.String("service", "QueueManager"),
			zap.String("method", "PublishEmailOrSms"),
			zap.String("data", string(data)),
		)
	}
}

func (q *queueManager) PublishKline(data []byte) {
	err := q.rabbitmqClient.Enqueue(rkKLineCreated, data)
	if err != nil {
		q.logger.Error2("can not enqueue kline", err,
			zap.String("service", "QueueManager"),
			zap.String("method", "PublishKline"),
			zap.String("data", string(data)),
		)
	}
}

func NewQueueManager(rabbitmqClient platform.RabbitMqClient, logger platform.Logger) QueueManager {
	return &queueManager{
		rabbitmqClient: rabbitmqClient,
		logger:         logger,
	}
}
