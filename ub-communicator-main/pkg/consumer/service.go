// Package consumer implements RabbitMQ message consuming with a worker pool pattern.
// It declares a topic exchange, binds a queue, and dispatches incoming messages
// to a configurable number of worker goroutines for parallel processing.
//
// Key types:
//   - Service: main consumer lifecycle (Consume blocks until context cancellation)
//   - Pool: manages worker goroutines for concurrent message processing
//   - Worker: individual goroutine that calls messaging.Service.Send()
//   - Collector: channels for sending work to the pool and signaling shutdown
package consumer

import (
	"context"
	"fmt"
	"ub-communicator/pkg/messaging"
	"ub-communicator/pkg/platform"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Service defines the RabbitMQ consumer lifecycle.
// Call Consume(ctx) to start consuming; it blocks until the context is cancelled.
type Service interface {
	Consume(ctx context.Context) error
}

type service struct {
	rc      platform.RabbitMqClient
	ms      messaging.Service
	pool    Pool
	logger  platform.Logger
	configs platform.Configs
}

func (s *service) Consume(ctx context.Context) error {
	ch, err := s.rc.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get rabbitmq channel: %w", err)
	}
	defer ch.Close()

	exchange := s.configs.GetString("rabbitmq.exchange")
	if exchange == "" {
		exchange = "email_exchange"
	}
	queueName := s.configs.GetString("rabbitmq.queue_name")
	if queueName == "" {
		queueName = "email_queue_1"
	}
	binding := s.configs.GetString("rabbitmq.binding")

	exchangeType := s.configs.GetString("rabbitmq.exchange_type")
	if exchangeType == "" {
		exchangeType = "direct"
	}
	var amqpExchangeType string
	switch exchangeType {
	case "topic":
		amqpExchangeType = amqp.ExchangeTopic
	case "fanout":
		amqpExchangeType = amqp.ExchangeFanout
	default:
		amqpExchangeType = amqp.ExchangeDirect
	}

	err = ch.ExchangeDeclare(exchange, amqpExchangeType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %q: %w", exchange, err)
	}

	// --- DLQ setup ---
	dlqExchange := s.configs.GetString("rabbitmq.dlq_exchange")
	if dlqExchange == "" {
		dlqExchange = "email_dlq_exchange"
	}
	dlqQueue := s.configs.GetString("rabbitmq.dlq_queue")
	if dlqQueue == "" {
		dlqQueue = "email_dlq"
	}
	dlqRoutingKey := s.configs.GetString("rabbitmq.dlq_routing_key")
	if dlqRoutingKey == "" {
		dlqRoutingKey = "email_dlq"
	}

	err = ch.ExchangeDeclare(dlqExchange, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ exchange %q: %w", dlqExchange, err)
	}
	dlq, err := ch.QueueDeclare(dlqQueue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ queue %q: %w", dlqQueue, err)
	}
	err = ch.QueueBind(dlq.Name, dlqRoutingKey, dlqExchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ queue %q: %w", dlq.Name, err)
	}
	// --- End DLQ setup ---

	mainQueueArgs := amqp.Table{
		"x-dead-letter-exchange":    dlqExchange,
		"x-dead-letter-routing-key": dlqRoutingKey,
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, mainQueueArgs)
	if err != nil {
		return fmt.Errorf("failed to declare queue %q: %w", queueName, err)
	}

	err = ch.QueueBind(q.Name, binding, exchange, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue %q to exchange %q: %w", q.Name, exchange, err)
	}

	// WHY: autoAck=true means messages are acknowledged immediately on delivery
	// from RabbitMQ, not after successful processing. This trades at-least-once
	// delivery for throughput. To prevent message loss on crash, change to
	// autoAck=false and call d.Ack(false) after successful Send().
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming from queue %q: %w", q.Name, err)
	}

	workerCount := s.configs.GetInt("consumer.worker_count")
	if workerCount <= 0 {
		workerCount = 5
	}
	collector := s.pool.StartDispatcher(workerCount)

	s.logger.Info("consumer started, waiting for messages")

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				collector.End <- true
				return fmt.Errorf("rabbitmq delivery channel closed unexpectedly")
			}
			message, err := s.ms.CreateMessage(d.Body)
			if err != nil {
				s.logger.Error("failed to parse message", zap.Error(err))
				// Nack with requeue=false — malformed body will never parse successfully.
				if nackErr := d.Nack(false, false); nackErr != nil {
					s.logger.Error("failed to nack unparseable message", zap.Error(nackErr))
				}
				continue
			}
			select {
			case collector.Work <- Work{Message: message, Delivery: d}:
			case <-ctx.Done():
				s.logger.Info("shutdown while queuing work, stopping consumer")
				collector.End <- true
				return ctx.Err()
			}
		case <-ctx.Done():
			s.logger.Info("shutdown signal received, stopping consumer")
			collector.End <- true
			return ctx.Err()
		}
	}
}

// NewConsumerService creates a consumer that reads from RabbitMQ and dispatches
// messages through the given worker pool using the messaging service.
func NewConsumerService(rc platform.RabbitMqClient, ms messaging.Service, pool Pool, logger platform.Logger, configs platform.Configs) Service {
	return &service{
		rc:      rc,
		ms:      ms,
		pool:    pool,
		logger:  logger,
		configs: configs,
	}
}
