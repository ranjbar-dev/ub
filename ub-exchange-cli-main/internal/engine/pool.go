package engine

import "sync/atomic"

type pool struct {
	workerCount                 int
	collector                   chan *work
	workers                     []*worker
	obp                         OrderbookProvider
	cbm                         *callBackManager
	shouldCallPostOrderMatching *atomic.Bool
	logger                      Logger
}

func (p *pool) run() {
	for i := 0; i < p.workerCount; i++ {
		worker := newWorker(p.collector, i, p.cbm, p.obp, p.shouldCallPostOrderMatching, p.logger)
		p.workers = append(p.workers, worker)
		go worker.start()
	}
}

func (p *pool) addWork(work *work) {
	select {
	case p.collector <- work:
	default:
		p.logger.Warn("engine pool at capacity, blocking until worker available")
		p.collector <- work
	}
}

// Stop stops background workers
func (p *pool) stop() {
	for i := range p.workers {
		p.workers[i].stop()
	}
}

func newPool(workerCount int, obp OrderbookProvider, cbm *callBackManager, shouldCall *atomic.Bool, logger Logger) *pool {
	collector := make(chan *work, 1000)
	return &pool{
		workerCount:                 workerCount,
		collector:                   collector,
		obp:                         obp,
		cbm:                         cbm,
		shouldCallPostOrderMatching: shouldCall,
		logger:                      logger,
	}
}
