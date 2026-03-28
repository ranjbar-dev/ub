// Package binance_test tests the Binance WebSocket client integration. Covers:
//   - Running the WebSocket connection with depth stream subscription
//   - Processing trade, depth, kline, and ticker messages from the WebSocket feed
//   - Verifying proper connection lifecycle (dial, write subscribe message, close)
//
// Test data: mock WebSocket client and connection, mock processor for trade/depth/kline/ticker
// events, and a single active BTC-USDT pair fixture.
package binance_test

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchangews/binance"
	"exchange-go/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestBinanceWs_Run(t *testing.T) {
	wsConnectionMocked := new(mocks.WsConnectionMocked)
	wsConnectionMocked.On("ReadMessage", mock.Anything).Return()

	wsConnectionMocked.On("WriteMessage", 1, mock.Anything).Once().Return(nil)
	wsConnectionMocked.On("SetPingHandler", mock.Anything).Once().Return(nil)
	wsConnectionMocked.On("Close").Once().Return(nil)

	wsClient := new(mocks.WsClient)
	wsClient.On("Dial", mock.Anything, mock.Anything, mock.Anything).Once().Return(wsConnectionMocked, nil)
	processor := new(mocks.Processor)
	processor.On("ProcessTrade", mock.Anything, mock.Anything).Once().Return()
	processor.On("ProcessDepth", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()
	processor.On("ProcessKline", mock.Anything, mock.Anything).Once().Return()
	processor.On("ProcessTicker", mock.Anything, mock.Anything).Once().Return()

	logger := new(mocks.Logger)
	activePairs := []currency.Pair{
		{
			ID:              1,
			Name:            "BTC-USDT",
			IsActive:        true,
			Spread:          1.5,
			ShowDigits:      6,
			BasisCoinID:     1,
			DependentCoinID: 2,
		},
	}
	binanceWs := binance.NewWs(wsClient, processor, logger, activePairs)
	depthStream := []string{binance.DepthStream}
	binanceWs.Run(context.Background(), depthStream)
	time.Sleep(50 * time.Millisecond)
	processor.AssertExpectations(t)
}
