package binance

import (
	"exchange-go/internal/platform"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

const (
	rateLimitRequest rateLimitType = iota
	rateLimitOrder
)

type rateLimitType int

type rateLimit struct {
	rateLimitType string
	interval      string
	intervalNum   int
	limit         int
	usage         int
}

type rateLimits map[rateLimitType]*rateLimit

type rateLimitHandler struct {
	logger platform.Logger
	limits rateLimits
}

func (r *rateLimitHandler) canPlaceOrder() bool {
	if !r.canRequest(NewOrderReuqestWeight) {
		return false
	}
	return r.limits[rateLimitOrder].usage < r.limits[rateLimitOrder].limit
}

func (r *rateLimitHandler) canRequest(weight int) bool {
	return r.limits[rateLimitRequest].usage+weight <= r.limits[rateLimitRequest].limit
}

func (r *rateLimitHandler) updateUsingHeader(header http.Header) {
	if orderCount := header.Get("x-mbx-order-count-10s"); orderCount != "" {
		//it means this is for order
		weight, err := strconv.Atoi(orderCount)
		if err == nil {
			r.limits[rateLimitOrder].usage = weight
		} else {
			r.logger.Warn("rate limit header does not exists or can not be converted to int",
				zap.String("service", "binanceRateLimitHandler"),
				zap.String("method", "updateUsingHeader"),
				zap.String("header", "x-mbx-order-count-10s"),
				zap.Error(err),
			)
			//we would never reach here but in case we cosider we have used 1 weight for orders
			r.limits[rateLimitOrder].usage += 1
		}
	}

	requestWeight := header.Get("x-mbx-used-weight-1m")
	if requestWeight != "" {
		weight, err := strconv.Atoi(requestWeight)
		if err == nil {
			r.limits[rateLimitRequest].usage = weight
			return
		}
		r.logger.Warn("rate limit header does not exists or can not be converted to int",
			zap.String("service", "binanceRateLimitHandler"),
			zap.String("method", "updateUsingHeader"),
			zap.String("header", "x-mbx-used-weight-1m"),
			zap.Error(err),
		)
		//we would never reach here but in case we cosider we have used 10 weight
		r.limits[rateLimitRequest].usage += 10
	}
}

func newRateLimitHandler(logger platform.Logger) *rateLimitHandler {
	rateLimits := make(rateLimits)
	rateLimits[rateLimitRequest] = &rateLimit{
		rateLimitType: "REQUEST_WEIGHT",
		interval:      "MINUTE",
		intervalNum:   1,
		limit:         1200,
		usage:         0,
	}
	rateLimits[rateLimitOrder] = &rateLimit{
		rateLimitType: "ORDERS",
		interval:      "SECOND",
		intervalNum:   10,
		limit:         50,
		usage:         0,
	}

	return &rateLimitHandler{
		logger: logger,
		limits: rateLimits,
	}
}
