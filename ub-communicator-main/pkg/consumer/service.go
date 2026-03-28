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

	exchange := s.configs.GetString("rabbitmq.exchange")
	if exchange == "" {
		exchange = "messages"
	}
	queueName := s.configs.GetString("rabbitmq.queue_name")
	if queueName == "" {
		queueName = "messages.command.send.consumer"
	}
	binding := s.configs.GetString("rabbitmq.binding")
	if binding == "" {
		binding = "messages.command.send"
	}

	// WHY: Using topic exchange (not direct/fanout) to allow future multi-consumer
	// routing based on message type patterns (e.g., messages.command.send.email).
	err = ch.ExchangeDeclare(exchange, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %q: %w", exchange, err)
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
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
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
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
				continue
			}
			collector.Work <- Work{Message: message}
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
