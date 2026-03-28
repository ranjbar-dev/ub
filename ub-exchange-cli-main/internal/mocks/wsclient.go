package mocks

import (
	"context"
	"exchange-go/internal/platform"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"
)

//Using testify mocks
type WsConnectionMocked struct {
	mock.Mock
}

var count int

func (m *WsConnectionMocked) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}
func (m *WsConnectionMocked) ReadMessage() (messageType int, p []byte, err error) {
	_ = m.Called()
	return m.CustomReadMessage()
	//return args.Int(0), args.Get(1).([]byte), args.Error(2)

}
func (m *WsConnectionMocked) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *WsConnectionMocked) SetPongHandler(h func(appData string) error) {
	_ = m.Called(h)
}

func (m *WsConnectionMocked) SetPingHandler(h func(appData string) error) {
	_ = m.Called(h)
}

type WsClient struct {
	mock.Mock
}

func (m *WsClient) Dial(ctx context.Context, urlStr string, requestHeader http.Header) (platform.WsConnection, error) {
	args := m.Called(ctx, urlStr, requestHeader)
	return args.Get(0).(platform.WsConnection), args.Error(1)
}

func (m *WsConnectionMocked) CustomReadMessage() (messageType int, p []byte, err error) {
	count++

	if count == 1 {
		tradeStream := "{" +
			"\"stream\":\"BTCUSDT@trade\"," +
			"\"data\":{" +
			"\"e\":\"trade\"," +
			"\"E\":1611490962000," +
			"\"s\":\"BTCUSDT\"," +
			"\"t\":45421," +
			"\"p\":\"32000\"," +
			"\"q\":\"0.01\"," +
			"\"b\":12345," +
			"\"a\":12345," +
			"\"T\":1611490962000," +
			"\"m\":true," +
			"\"M\":true" +
			"}}"
		return 1, []byte(tradeStream), nil
	}

	if count == 2 {
		depthStream := "{" +
			"\"stream\":\"BTCUSDT@depth\"," +
			"\"data\":{" +
			"\"e\":\"depthUpdate\"," +
			"\"E\":1611490962000," +
			"\"s\":\"BTCUSDT\"," +
			"\"U\":1245555," +
			"\"u\":1245552," +
			"\"b\":[[\"32725.19000000\",\"11.3120000\"]]," +
			"\"a\":[[\"32723.57000000\",\"1.86000000\"]]" +
			"}}"
		return 1, []byte(depthStream), nil
	}

	if count == 3 {
		klineStream := "{" +
			"\"stream\":\"BTCUSDT@kline_1m\"," +
			"\"data\":{" +
			"\"e\": \"kline\"," +
			"\"E\": 123456789," +
			"\"s\": \"BTCUSDT\"," +
			"\"k\":{" +
			"\"t\":123400000," +
			"\"T\":123460000," +
			"\"s\":\"BTCUSDT\"," +
			"\"i\":\"1m\"," +
			"\"f\":20254," +
			"\"L\":30254," +
			"\"o\":\"32200\"," +
			"\"c\":\"32200\"," +
			"\"h\":\"32200\"," +
			"\"l\":\"32200\"," +
			"\"v\":\"1000\"," +
			"\"n\":100," +
			"\"q\":\"1.0000\"," +
			"\"V\":\"500\"," +
			"\"Q\":\"0.500\"," +
			"\"B\":\"123456\"" +
			"}}}"

		return 1, []byte(klineStream), nil
	}

	if count == 4 {
		tickerStream := "{" +
			"\"stream\":\"BTCUSDT@ticker\"," +
			"\"data\":{" +
			"\"e\":\"24hrTicker\"," +
			"\"E\":123456789," +
			"\"s\":\"BTCUSDT\"," +
			"\"p\":\"120\"," +
			"\"P\":\"25\"," +
			"\"w\":\"120\"," +
			"\"x\":\"32100\"," +
			"\"c\":\"32300\"," +
			"\"Q\":\"10\"," +
			"\"b\":\"32100\"," +
			"\"B\":\"10\"," +
			"\"a\":\"32300\"," +
			"\"A\":\"10\"," +
			"\"o\":\"32100\"," +
			"\"h\":\"32400\"," +
			"\"l\":\"3200\"," +
			"\"v\":\"1000\"," +
			"\"q\":\"18\"," +
			"\"O\":0," +
			"\"C\":86400000," +
			"\"F\":0," +
			"\"L\":18150," +
			"\"n\":18151" +
			"}}"

		return 1, []byte(tickerStream), nil
	}

	return 1, []byte(""), fmt.Errorf("envTest")
}
