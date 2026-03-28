package consumer

import (
	"log"
	"ub-communicator/pkg/messaging"
)

// Work represents a unit of work dispatched from the consumer to a worker.
type Work struct {
	ID      int64
	Message messaging.Message
}

// Worker is a goroutine that processes Work items by calling messaging.Service.Send().
type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	Channel       chan Work
	End           chan bool
	Ms            messaging.Service
}

// Start launches the worker goroutine. It registers itself as available
// on the WorkerChannel, then waits for Work items on its private Channel.
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case work := <-w.Channel:
				if err := w.Ms.Send(work.Message); err != nil {
					// Error is already logged inside Send(); this captures any propagated error.
					log.Printf("worker [%d] send error: %v", w.ID, err)
				}
			case <-w.End:
				return
			}
		}
	}()
}

// Stop signals the worker goroutine to terminate.
func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}
