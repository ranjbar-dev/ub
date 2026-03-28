package engine

type pool struct {
	workerCount int
	collector   chan *work
	workers     []*worker
}

func (p *pool) run() {
	for i := 0; i < p.workerCount; i++ {
		worker := newWorker(p.collector, i, cbm)
		p.workers = append(p.workers, worker)
		go worker.start()
	}
	//p.stopChan = make(chan bool)
	//<-p.stopChan
}

func (p *pool) addWork(work *work) {
	p.collector <- work
}

// Stop stops background workers
func (p *pool) stop() {
	for i := range p.workers {
		p.workers[i].stop()
	}
	//p.stopChan <- true
}

func newPool(workerCount int) *pool {
	collector := make(chan *work, 1000)
	//p.stopChan = make(chan bool)
	return &pool{
		workerCount: workerCount,
		collector:   collector,
	}
}
