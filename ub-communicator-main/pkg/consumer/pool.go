package consumer

import (
	"fmt"
	"ub-communicator/pkg/messaging"
)

// Collector holds the channels used to send work to the pool and signal shutdown.
type Collector struct {
	Work chan Work // Send Work items here for dispatch to workers.
	End  chan bool // Send true to stop all workers and the dispatcher.
}

// Pool manages a set of worker goroutines that process messages in parallel.
type Pool interface {
	// StartDispatcher launches workerCount goroutines and returns a Collector
	// for sending work and signaling shutdown.
	StartDispatcher(workerCount int) Collector
}

type pool struct {
	ms            messaging.Service
	workerChannel chan chan Work
}

func (p *pool) StartDispatcher(workerCount int) Collector {
	var workers []Worker
	input := make(chan Work, 100) // buffered to prevent blocking on bursts
	end := make(chan bool)
	collector := Collector{Work: input, End: end}

	for i := 1; i <= workerCount; i++ {
		fmt.Printf("starting worker: %d\n", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan Work),
			WorkerChannel: p.workerChannel,
			End:           make(chan bool),
			Ms:            p.ms,
		}
		worker.Start()
		workers = append(workers, worker)
	}

	// Dispatcher goroutine: routes incoming work to available workers.
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop()
				}
				return
			case work := <-input:
				// WHY: Worker availability is managed via a "channel of channels" pattern.
				// Each worker sends its private work channel onto workerChannel when idle.
				// The dispatcher picks the first available worker by reading from workerChannel.
				// This ensures work is only dispatched to workers that are ready to receive.
				workChan := <-p.workerChannel
				workChan <- work
			}
		}
	}()

	return collector
}

// NewPool creates a new worker pool that dispatches messages via the given service.
func NewPool(ms messaging.Service) Pool {
	return &pool{
		ms:            ms,
		workerChannel: make(chan chan Work),
	}
}
