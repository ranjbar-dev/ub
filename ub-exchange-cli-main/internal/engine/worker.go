package engine

import (
	"fmt"
	"sync/atomic"

	"go.uber.org/zap"
)

const (
	WorkTypeMatching = "matching"
	WorkTypeInQueue  = "inqueue"
)

type work struct {
	order Order
}

type worker struct {
	ID                          int
	workChan                    chan *work
	quit                        chan bool
	callBackManager             *callBackManager
	obp                         OrderbookProvider
	shouldCallPostOrderMatching *atomic.Bool
	logger                      Logger
}

func (w *worker) start() {
	for {
		select {
		case <-w.quit:
			return
		case work := <-w.workChan:
			select {
			case <-w.quit:
				return
			default:
			}
			w.processOrder(work.order)
		}

	}
}

func (w *worker) processOrder(o Order) {
	ob := newOrderBook(o.Pair, w.obp, w.logger)
	doneOrders, partialOrder, err := ob.processOrder(o)
	if err != nil {
		w.logger.Warn("error in engine:ProcessOrder",
			zap.Error(err),
			zap.String("orderId", o.ID),
		)
		return
	}

	if w.shouldCallPostOrderMatching.Load() {
		engineMatchingResult := w.callBackManager.callBack(doneOrders, partialOrder)
		var removingDoneOrders []Order
		for _, id := range engineMatchingResult.RemovingDoneOrderIds {
			idString := fmt.Sprintf("%011d", id)
			for _, doneOrder := range doneOrders {
				if idString == doneOrder.ID {
					removingDoneOrders = append(removingDoneOrders, doneOrder)
				}
			}
		}
		ob.rewriteOrderBook(removingDoneOrders, engineMatchingResult.RemainingPartialOrder)
	} else {
		//this is only for test env and would not happen in production
		ob.rewriteOrderBook(doneOrders, partialOrder)
	}

}

func (w *worker) stop() {
	w.quit <- true
}

func newWorker(workChan chan *work, ID int, callBackManager *callBackManager, obp OrderbookProvider, shouldCall *atomic.Bool, logger Logger) *worker {
	return &worker{
		ID:                          ID,
		workChan:                    workChan,
		quit:                        make(chan bool, 1),
		callBackManager:             callBackManager,
		obp:                         obp,
		shouldCallPostOrderMatching: shouldCall,
		logger:                      logger,
	}
}
