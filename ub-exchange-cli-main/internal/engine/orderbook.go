package engine

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/shopspring/decimal"
)

const (
	SideBid = "bid"
	SideAsk = "ask"
	// redisTimeout is the maximum time to wait for a Redis operation before giving up.
	redisTimeout = 5 * time.Second
)

type OrderBookProviderParams struct {
	Pair     string
	Side     string
	Price    string
	MinPrice string
	MaxPrice string
}

// OrderbookProvider manages order book persistence in Redis sorted sets.
type OrderbookProvider interface {
	// GetOrders retrieves all orders at the specified price level and side from the order book.
	GetOrders(ctx context.Context, params OrderBookProviderParams) ([]Order, error)
	// RewriteOrderBook atomically updates the order book after matching by removing done orders and inserting a partial fill.
	RewriteOrderBook(ctx context.Context, doneOrders []Order, partialOrder *Order) error
	// RemoveOrder deletes a specific order from the order book.
	RemoveOrder(ctx context.Context, order Order) error
	// PopOrders removes and returns all orders matching the given parameters, used during order matching.
	PopOrders(ctx context.Context, params OrderBookProviderParams) ([]Order, error)
	// Exists checks whether a given order is present in the order book.
	Exists(ctx context.Context, order Order) (bool, error)
}

type orderBook struct {
	orderbookProvider OrderbookProvider
	pair              string
	asks              []Order
	bids              []Order
}

func (ob *orderBook) processOrder(o Order) (doneOrders []Order, partialOrder *Order, error error) {
	if o.Price != "" {
		return ob.processLimitOrder(o)
	}
	return ob.processMarketOrder(o)

}

func (ob *orderBook) processLimitOrder(o Order) (done []Order, partialOrder *Order, error error) {
	var partial *Order
	price, err := o.GetPrice()
	if err != nil {
		return nil, &o, err
	}
	quantityToTrade, err := o.GetQuantity()
	if err != nil {
		return nil, &o, err
	}
	sideToLoad := SideBid
	comparator := price.LessThanOrEqual
	if o.Side == SideBid {
		sideToLoad = SideAsk
		comparator = price.GreaterThanOrEqual
	}

	if err := ob.loadOrders(sideToLoad, o.Price, "", ""); err != nil {
		return nil, &o, err
	}
	bestOrder, bestPriceFound := ob.bestOrder(sideToLoad)
	bestOrderPrice, _ := bestOrder.GetPrice()

	partial = &o
	for quantityToTrade.Sign() > 0 && bestPriceFound && comparator(bestOrderPrice) {
		// Self-trade prevention: skip orders from the same user
		if o.UserID != "" && bestOrder.UserID != "" && o.UserID == bestOrder.UserID {
			bestOrder, bestPriceFound = ob.bestOrder(sideToLoad)
			if bestPriceFound {
				bestOrderPrice, _ = bestOrder.GetPrice()
			}
			continue
		}
		prevQuantity := quantityToTrade
		doneOrders, newPartial, remaining := ob.tradeOrders(*partial, bestOrder)
		done = append(done, doneOrders...)
		partial = newPartial
		quantityToTrade = remaining
		if quantityToTrade.Equal(prevQuantity) {
			break
		}
		bestOrder, bestPriceFound = ob.bestOrder(sideToLoad)
		bestOrderPrice, _ = bestOrder.GetPrice()
	}

	return done, partial, nil
}

func (ob *orderBook) tradeOrders(order Order, bestOrder Order) (doneOrders []Order, partialOrder *Order, remaining decimal.Decimal) {
	finalPrice := order.Price
	if finalPrice == "" {
		finalPrice = bestOrder.Price
	}

	orderQuantity, _ := order.GetQuantity()
	marketPrice, _ := order.GetMarketPrice()
	bestOrderQuantity, _ := bestOrder.GetQuantity()
	bestOrderPrice, _ := bestOrder.GetPrice()
	quantityToTrade := orderQuantity
	orderValue := orderQuantity.Mul(marketPrice)
	bestOrderValue := bestOrderQuantity.Mul(bestOrderPrice)

	//this is only valid for bid market orders
	if order.IsBidMarket() {
		if bestOrderValue.LessThanOrEqual(orderValue) {
			quantityToTrade = bestOrderQuantity
		} else {
			quantityToTrade = orderValue.Div(bestOrderPrice)
		}
	}

	if quantityToTrade.Equal(bestOrderQuantity) {
		//update  orders
		order.TradedWithOrderID = bestOrder.ID
		order.QuantityTraded = bestOrder.Quantity
		//order.TradePrice = order.Price
		order.TradePrice = finalPrice

		bestOrder.TradedWithOrderID = order.ID
		bestOrder.QuantityTraded = bestOrder.Quantity
		//bestOrder.TradePrice = order.Price
		bestOrder.TradePrice = finalPrice
		doneOrders := append(doneOrders, order, bestOrder)
		var partial *Order
		remaining := decimal.Zero
		if order.IsBidMarket() {
			remaining = orderValue.Sub(bestOrderValue).Div(marketPrice)
			//check if remaining is not zero then calculate partial
			if remaining.IsPositive() {
				quantity := remaining.StringFixed(16)
				partialOrder := newOrder(order.Pair, order.ID, order.Side, quantity, order.Price,
					order.MarketPrice, order.Timestamp, order.UserID)
				partial = &partialOrder
			}
		}
		return doneOrders, partial, remaining
	}

	if quantityToTrade.LessThan(bestOrderQuantity) {
		partialQuantity := bestOrderQuantity.Sub(quantityToTrade)
		partialQuantityString := partialQuantity.String()
		//update orders
		order.QuantityTraded = quantityToTrade.StringFixed(16)
		order.TradedWithOrderID = bestOrder.ID
		//order.TradePrice = order.Price
		order.TradePrice = finalPrice

		bestOrder.TradedWithOrderID = order.ID
		bestOrder.QuantityTraded = quantityToTrade.StringFixed(16)
		//bestOrder.TradePrice = order.Price
		bestOrder.TradePrice = finalPrice

		po := newOrder(ob.pair, bestOrder.ID, bestOrder.Side, partialQuantityString, bestOrder.Price,
			bestOrder.MarketPrice, bestOrder.Timestamp, bestOrder.UserID)
		po.IsAlreadyInOrderBook = true
		doneOrders := append(doneOrders, order, bestOrder)
		return doneOrders, &po, decimal.Zero
	}

	//here means else
	partialQuantity := quantityToTrade.Sub(bestOrderQuantity)
	partialQuantityString := partialQuantity.String()
	//update orders
	order.TradedWithOrderID = bestOrder.ID
	order.QuantityTraded = bestOrder.Quantity
	//order.TradePrice = order.Price
	order.TradePrice = finalPrice

	bestOrder.TradedWithOrderID = order.ID
	bestOrder.QuantityTraded = bestOrder.Quantity
	//bestOrder.TradePrice = order.Price
	bestOrder.TradePrice = finalPrice

	po := newOrder(ob.pair, order.ID, order.Side, partialQuantityString, order.Price, order.MarketPrice, order.Timestamp, order.UserID)
	doneOrders = append(doneOrders, order, bestOrder)
	return doneOrders, &po, partialQuantity
}

func (ob *orderBook) processMarketOrder(o Order) (doneOrders []Order, partialOrder *Order, error error) {
	var done []Order
	var partial *Order

	maxPrice, err := o.GetMaxPrice()
	if err != nil {
		return done, partialOrder, err
	}

	minPrice, err := o.GetMinPrice()
	if err != nil {
		return done, partialOrder, err
	}

	if minPrice.GreaterThan(maxPrice) {
		return done, &o, fmt.Errorf("minPrice (%s) must be <= maxPrice (%s)", minPrice, maxPrice)
	}

	quantityToTrade, err := o.GetQuantity()
	if err != nil {
		return done, &o, err
	}
	sideToLoad := SideBid
	comparator := minPrice.LessThanOrEqual

	if o.Side == SideBid {
		sideToLoad = SideAsk
		comparator = maxPrice.GreaterThanOrEqual
	}

	if err := ob.loadOrders(sideToLoad, o.Price, minPrice.String(), maxPrice.String()); err != nil {
		return done, &o, err
	}
	bestOrder, bestPriceFound := ob.bestOrder(sideToLoad)
	bestOrderPrice, _ := bestOrder.GetPrice()
	partial = &o
	for quantityToTrade.Sign() > 0 && bestPriceFound && comparator(bestOrderPrice) {
		// Self-trade prevention: skip orders from the same user
		if o.UserID != "" && bestOrder.UserID != "" && o.UserID == bestOrder.UserID {
			bestOrder, bestPriceFound = ob.bestOrder(sideToLoad)
			if bestPriceFound {
				bestOrderPrice, _ = bestOrder.GetPrice()
			}
			continue
		}
		prevQuantity := quantityToTrade
		doneOrders, newPartial, remaining := ob.tradeOrders(*partial, bestOrder)
		done = append(done, doneOrders...)
		partial = newPartial
		quantityToTrade = remaining
		if quantityToTrade.Equal(prevQuantity) {
			break
		}
		bestOrder, bestPriceFound = ob.bestOrder(sideToLoad)
		bestOrderPrice, _ = bestOrder.GetPrice()
	}

	// C6: Market orders are IOC (Immediate Or Cancel) — if no trades executed,
	// return nil partial to signal cancellation instead of returning the unfilled order
	if len(done) == 0 {
		return done, nil, nil
	}

	return done, partial, nil
}

func (ob *orderBook) bestOrder(sideToLoad string) (order Order, found bool) {
	if sideToLoad == SideAsk {
		if len(ob.asks) > 0 {
			best := ob.asks[0]
			ob.asks = ob.asks[1:]
			return best, true
		}
	} else {
		if len(ob.bids) > 0 {
			best := ob.bids[len(ob.bids)-1]
			ob.bids = ob.bids[:len(ob.bids)-1]
			return best, true
		}
	}

	return order, false
}

func (ob *orderBook) rewriteOrderBook(doneOrders []Order, partialOrder *Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	err := ob.orderbookProvider.RewriteOrderBook(ctx, doneOrders, partialOrder)
	if err != nil {
		logHandler.Warn("error in engine:rewriteOrderBook",
			zap.Error(err),
		)
	}
	return err
}

func (ob *orderBook) loadOrders(side string, price string, minPrice string, maxPrice string) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	params := OrderBookProviderParams{
		Pair:     ob.pair,
		Side:     side,
		Price:    price,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}
	orders, err := ob.orderbookProvider.GetOrders(ctx, params)
	if err != nil {
		logHandler.Warn("error in engine:loadOrders",
			zap.Error(err),
			zap.String("pair", params.Pair),
			zap.String("side", params.Side),
			zap.String("price", params.Price),
		)
		return err
	}

	if side == SideAsk {
		// Sort asks ascending by price (lowest first = best ask at index 0), FIFO tiebreak
		asksCopy := make([]Order, len(orders))
		copy(asksCopy, orders)
		sort.Slice(asksCopy, func(i, j int) bool {
			firstPrice, errI := decimal.NewFromString(asksCopy[i].Price)
			secondPrice, errJ := decimal.NewFromString(asksCopy[j].Price)
			if errI != nil || errJ != nil {
				return i < j
			}
			if firstPrice.Equal(secondPrice) {
				firstID, errI := strconv.ParseInt(asksCopy[i].ID, 10, 64)
				secondID, errJ := strconv.ParseInt(asksCopy[j].ID, 10, 64)
				if errI != nil || errJ != nil {
					return i < j
				}
				return firstID < secondID
			}
			return firstPrice.LessThan(secondPrice)
		})
		ob.asks = asksCopy
	} else {
		// Sort bids ascending by price (highest at end, popped first by bestOrder).
		// Same price: higher IDs first so oldest (lowest ID) is at end = FIFO via pop.
		bidsCopy := make([]Order, len(orders))
		copy(bidsCopy, orders)
		sort.Slice(bidsCopy, func(i, j int) bool {
			firstPrice, errI := decimal.NewFromString(bidsCopy[i].Price)
			secondPrice, errJ := decimal.NewFromString(bidsCopy[j].Price)
			if errI != nil || errJ != nil {
				return i < j
			}
			if firstPrice.Equal(secondPrice) {
				firstID, errI := strconv.ParseInt(bidsCopy[i].ID, 10, 64)
				secondID, errJ := strconv.ParseInt(bidsCopy[j].ID, 10, 64)
				if errI != nil || errJ != nil {
					return i < j
				}
				return firstID > secondID
			}
			return firstPrice.LessThan(secondPrice)
		})
		ob.bids = bidsCopy
	}

	return nil
}

func (ob *orderBook) removeOrder(o Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()
	return ob.orderbookProvider.RemoveOrder(ctx, o)
}

func (ob *orderBook) getInQueueOrder(price string) (orders []Order, error error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	bidParams := OrderBookProviderParams{
		Pair:     ob.pair,
		Side:     SideBid,
		Price:    price,
		MinPrice: "",
		MaxPrice: "",
	}

	bidOrders, err := ob.orderbookProvider.GetOrders(ctx, bidParams)
	if err != nil {
		return orders, err
	}

	orders = append(orders, bidOrders...)

	askParams := OrderBookProviderParams{
		Pair:     ob.pair,
		Side:     SideAsk,
		Price:    price,
		MinPrice: "",
		MaxPrice: "",
	}

	askOrders, err := ob.orderbookProvider.GetOrders(ctx, askParams)
	if err != nil {
		return orders, err
	}
	orders = append(orders, askOrders...)

	return orders, nil
}

// popInQueueOrders atomically reads and removes in-queue orders, preventing race conditions
// where workers could match against orders being processed by HandleInQueueOrders.
func (ob *orderBook) popInQueueOrders(price string) (orders []Order, error error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisTimeout)
	defer cancel()

	bidParams := OrderBookProviderParams{
		Pair:     ob.pair,
		Side:     SideBid,
		Price:    price,
		MinPrice: "",
		MaxPrice: "",
	}
	bidOrders, err := ob.orderbookProvider.PopOrders(ctx, bidParams)
	if err != nil {
		return orders, err
	}
	orders = append(orders, bidOrders...)

	askParams := OrderBookProviderParams{
		Pair:     ob.pair,
		Side:     SideAsk,
		Price:    price,
		MinPrice: "",
		MaxPrice: "",
	}
	askOrders, err := ob.orderbookProvider.PopOrders(ctx, askParams)
	if err != nil {
		return orders, err
	}
	orders = append(orders, askOrders...)

	return orders, nil
}

func (ob *orderBook) orderExists(ctx context.Context, o Order) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, redisTimeout)
	defer cancel()
	return ob.orderbookProvider.Exists(ctx, o)
}

func newOrderBook(pair string, obp OrderbookProvider) orderBook {
	orderBook := orderBook{
		pair:              pair,
		orderbookProvider: obp,
	}
	return orderBook
}
