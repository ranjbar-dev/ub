// Package binance tests the Binance rate limit handler (white-box). Covers:
//   - Checking request rate limits (under and over the 1200 request/minute threshold)
//   - Checking order placement rate limits (under and over the 50 order/10s threshold)
//   - Composite rate limit checks (order blocked when request limit is exceeded)
//   - Updating rate limit usage counters from Binance API response headers
//
// Test data: in-memory rate limit handler initialized via newRateLimitHandler
// with manipulated usage counters and HTTP response headers.
package binance

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitHandler_canRequest(t *testing.T) {
	rateLimitHandler := newRateLimitHandler(nil)
	canRequest := rateLimitHandler.canRequest(1)
	assert.True(t, canRequest)

	rateLimitHandler.limits[rateLimitRequest].usage = 1200
	canRequest = rateLimitHandler.canRequest(1)
	assert.False(t, canRequest)
}

func TestRateLimitHandler_canPlaceOrder(t *testing.T) {
	rateLimitHandler := newRateLimitHandler(nil)
	canPlaceOrder := rateLimitHandler.canPlaceOrder()
	assert.True(t, canPlaceOrder)

	rateLimitHandler.limits[rateLimitOrder].usage = 50
	canPlaceOrder = rateLimitHandler.canPlaceOrder()
	assert.False(t, canPlaceOrder)

	rateLimitHandler.limits[rateLimitOrder].usage = 0
	rateLimitHandler.limits[rateLimitRequest].usage = 1200
	canPlaceOrder = rateLimitHandler.canPlaceOrder()
	assert.False(t, canPlaceOrder)
}

func TestRateLimitHandler_updateUsingHeader(t *testing.T) {
	rateLimitHandler := newRateLimitHandler(nil)
	header := http.Header{}
	header.Set("x-mbx-order-count-10s", "10")
	header.Set("x-mbx-used-weight-1m", "20")
	rateLimitHandler.updateUsingHeader(header)
	orderRateUsage := rateLimitHandler.limits[rateLimitOrder].usage
	requestRateUsage := rateLimitHandler.limits[rateLimitRequest].usage
	assert.Equal(t, 10, orderRateUsage)
	assert.Equal(t, 20, requestRateUsage)
}
