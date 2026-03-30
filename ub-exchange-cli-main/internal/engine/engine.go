package engine

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
)

const engineRedisTimeout = 5 * time.Second

var logHandler Logger

// ResultHandler is a callback interface for processing matched trades.
type ResultHandler interface {
	// CallBack receives completed orders and an optional partial fill, and returns the settlement result.
	CallBack(doneOrders []Order, partialOrder *Order) MatchingResult
}

type MatchingResult struct {
	Err                   error
	RemainingPartialOrder *Order
	RemovingDoneOrderIds  []int64
}

// Logger is a structured logging interface for engine operations.
type Logger interface {
	// Warn logs a message at the warn level with optional structured fields.
	Warn(msg string, fields ...zap.Field)
	// Info logs a message at the info level with optional structured fields.
	Info(msg string, fields ...zap.Field)
}

// Engine controls the order matching engine lifecycle.
type Engine interface {
	// Run starts the worker pool with the given number of workers and optionally launches the order dispatcher.
	Run(workerCount int, shouldStartDispatcher bool)
	// Stop gracefully shuts down the engine, stopping the dispatcher and draining the worker pool.
	Stop()
	// DispatchManually triggers a single manual order dispatch. Intended for use in tests only.
	DispatchManually() error
	// SetPostOrderMatchingCall enables or disables post-match callbacks.
	SetPostOrderMatchingCall(shouldCall bool) error
	// SubmitOrder adds an order to the end of the processing queue.
	SubmitOrder(order Order) error
	// RemoveOrder removes an order from both the processing queue and the order book.
	RemoveOrder(order Order) error
	// HandleInQueueOrders processes stop/limit orders in the order book that are triggered by a price change on the given pair.
	HandleInQueueOrders(pair string, price string) error
	// RetrieveOrder restores a missing order back into the queue or order book, giving it priority via left-push.
	RetrieveOrder(order Order) error
}

type engine struct {
	pool                        *pool
	queue                       *queue
	env                         string
	quit                        chan bool
	obp                         OrderbookProvider
	cbm                         *callBackManager
	shouldCallPostOrderMatching atomic.Bool
}

func (e *engine) SetPostOrderMatchingCall(shouldCall bool) error {
	e.shouldCallPostOrderMatching.Store(shouldCall)
	return nil
}

func (e *engine) Run(workerCount int, shouldStartDispatcher bool) {
	e.pool = newPool(workerCount, e.obp, e.cbm, &e.shouldCallPostOrderMatching)
	e.pool.run()
	if shouldStartDispatcher {
		go e.dispatchOrder()
	}
}

func (e *engine) Stop() {
	e.quit <- true
	e.pool.stop()
}

func (e *engine) DispatchManually() error {
	if e.env != "test" {
		return fmt.Errorf("DispatchManually: must only be called in test environment")
	}
	order, err := e.queue.lPop(context.Background())
	if err != nil && err != redis.Nil {
		logHandler.Warn("error in engine:dispatchOrder",
			zap.Error(err),
		)
		return nil
	}
	if err == redis.Nil {
		return nil
	}
	work := work{order: order}
	e.pool.addWork(&work)
	return nil
}

func (e *engine) dispatchOrder() {
	ctx := context.Background()
	backoff := time.Duration(0)
	maxBackoff := 30 * time.Second
	for {
		select {
		case <-e.quit:
			return
		default:
			select {
			case <-e.quit:
				return
			default:
			}

			order, err := e.queue.blPop(ctx, 1*time.Second)
			if err != nil && err != redis.Nil {
				logHandler.Warn("error in engine:dispatchOrder",
					zap.Error(err),
				)
				if backoff == 0 {
					backoff = 100 * time.Millisecond
				} else {
					backoff = backoff * 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
				time.Sleep(backoff)
				continue
			}
			backoff = 0
			if err == redis.Nil {
				continue
			}
			work := work{order: order}
			e.pool.addWork(&work)
		}
	}
}

func (e *engine) SubmitOrder(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), engineRedisTimeout)
	defer cancel()
	return e.queue.rPush(ctx, order)
}

func (e *engine) RemoveOrder(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), engineRedisTimeout)
	defer cancel()
	ob := newOrderBook(order.Pair, e.obp)
	err := e.queue.remove(ctx, order)
	if err != nil {
		logHandler.Warn("error in engine:RemoveOrder",
			zap.Error(err),
			zap.String("orderId", order.ID),
		)
		return err
	}
	err = ob.removeOrder(order)
	if err != nil {
		logHandler.Warn("error in engine:RemoveOrder",
			zap.Error(err),
			zap.String("orderId", order.ID),
		)
	}

	return err
}

func (e *engine) HandleInQueueOrders(pair string, price string) error {
	ob := newOrderBook(pair, e.obp)

	// H5: Use PopOrders (atomic read+remove) instead of GetOrders to prevent
	// race condition where workers match against orders being processed here.
	orders, err := ob.popInQueueOrders(price)
	if err != nil {
		return err
	}
	go func() {
		for i, o := range orders {
			var emptyDoneOrders []Order
			o.IsAlreadyInOrderBook = true
			engineMatchingResult := e.cbm.callBack(emptyDoneOrders, &orders[i])
			if engineMatchingResult.Err != nil {
				logHandler.Warn("error is not nil",
					zap.Error(engineMatchingResult.Err),
					zap.String("service", "engine"),
					zap.String("method", "HandleInQueueOrders"),
					zap.String("orderId", o.ID),
				)
				// Put the order back since callback failed
				ob.rewriteOrderBook(nil, &orders[i])
				continue
			}
			if engineMatchingResult.RemainingPartialOrder != nil {
				// Order not fully filled — put remaining partial back into the orderbook
				ob.rewriteOrderBook(nil, engineMatchingResult.RemainingPartialOrder)
			}
			// If RemainingPartialOrder is nil, order is fully filled — already removed by PopOrders
		}
	}()
	return nil
}

/**
 * In any case that orders in orderbook or queue got removed (for example for redis server failure),we call this method
 * to retrieve order to queue and orderbook, but since the order in our command are sorted from oldest to newest we lPush
 * order to queue to give them priority since we have lPop in our order dispatch
 */
func (e *engine) RetrieveOrder(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), engineRedisTimeout)
	defer cancel()
	if order.Price != "" {
		//it means limit order
		ob := newOrderBook(order.Pair, e.obp)
		exists, err := ob.orderExists(ctx, order)
		if err != nil {
			return err
		}
		if !exists {
			existsInQueue, err := e.queue.exists(ctx, order)
			if err != nil {
				return err
			}

			if !existsInQueue {
				err := e.queue.lPush(ctx, order)
				return err
			}
			return nil
		}

		return nil
	}

	//it means market order
	exists, err := e.queue.exists(ctx, order)
	if err != nil {
		return err
	}

	if !exists {
		err := e.queue.lPush(ctx, order)
		return err
	}

	return nil
}

func NewEngine(qh QueueHandler, obp OrderbookProvider, rh ResultHandler, logger Logger, env string) Engine {
	logHandler = logger
	cbm := getCallbackManager(rh)
	queue := newQueue(qh)
	quit := make(chan bool, 1)
	e := &engine{
		queue: queue,
		env:   env,
		quit:  quit,
		obp:   obp,
		cbm:   cbm,
	}
	e.shouldCallPostOrderMatching.Store(true)
	return e
}
