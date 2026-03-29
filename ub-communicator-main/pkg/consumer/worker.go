package consumer

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"ub-communicator/pkg/messaging"
)

// Work represents a unit of work dispatched from the consumer to a worker.
type Work struct {
	ID       int64
	Message  messaging.Message
	Delivery amqp.Delivery
}

// Worker is a goroutine that processes Work items by calling messaging.Service.Send().
type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	Channel       chan Work
	Ms            messaging.Service
	ctx           context.Context
	cancel        context.CancelFunc
}

// Start launches the worker goroutine. It registers itself as available
// on the WorkerChannel, then waits for Work items on its private Channel.
func (w *Worker) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("worker [%d] panic recovered: %v", w.ID, r)
			}
		}()
		for {
			// First blocking point: register as available.
			// Must also listen on ctx.Done() or Stop() deadlocks here.
			select {
			case w.WorkerChannel <- w.Channel:
			case <-w.ctx.Done():
				return
			}
			// Second blocking point: wait for work or shutdown.
			select {
			case work := <-w.Channel:
				if err := w.Ms.Send(work.Message); err != nil {
					log.Printf("worker [%d] send error: %v", w.ID, err)
					// Nack with requeue=true so the broker retries another worker.
					if nackErr := work.Delivery.Nack(false, true); nackErr != nil {
						log.Printf("worker [%d] nack error: %v", w.ID, nackErr)
					}
				} else {
					if ackErr := work.Delivery.Ack(false); ackErr != nil {
						log.Printf("worker [%d] ack error: %v", w.ID, ackErr)
					}
				}
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// Stop signals the worker goroutine to terminate.
func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.cancel()
}
