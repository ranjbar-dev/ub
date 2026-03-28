package binance

import (
	"context"
	"encoding/json"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchangews/types"
	"exchange-go/internal/platform"
	"exchange-go/internal/processor"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	StreamURL     = "wss://stream.binance.com:9443/stream"
	ExchangeName  = "binance"
	TradeStream   = "trade"
	DepthStream   = "depth"
	TickerStream  = "ticker"
	KlineStream1m = "kline_1m"
	KlineStream5m = "kline_5m"
	KlineStream1h = "kline_1h"
	KlineStream1d = "kline_1d"
)

type binanceWs struct {
	client      platform.WsClient
	processor   processor.Processor
	logger      platform.Logger
	activePairs []currency.Pair
	PairNameMap map[string]string
}

func (ws *binanceWs) Run(ctx context.Context, streams []string) {
	c, err := ws.client.Dial(ctx, StreamURL, nil)
	if err != nil {
		ws.logger.Error2("can not dial to binace websocket", err,
			zap.String("service", "binanceWebSocket"),
			zap.String("method", "Run"),
			zap.String("stramUrl", StreamURL),
		)
		os.Exit(1)
	}
	defer c.Close()

	c.SetPingHandler(func(appData string) error {
		err = c.WriteMessage(websocket.PongMessage, []byte{})
		if err != nil {
			ws.logger.Error2("can not write message in ping handler", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "Run"),
			)
			os.Exit(1)
		}
		return err
	})

	message, err := ws.createWsSubscribeRequest(streams)
	if err != nil {
		ws.logger.Error2("can not create ws subscribe request", err,
			zap.String("service", "binanceWebSocket"),
			zap.String("method", "Run"),
		)
		//Send 1 error to OS to restart again the container
		os.Exit(1)
	}

	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		ws.logger.Error2("can not write message", err,
			zap.String("service", "binanceWebSocket"),
			zap.String("method", "Run"),
		)
		os.Exit(1)
	}

	ws.readMessage(ctx, c)
}

func (ws *binanceWs) readMessage(ctx context.Context, c platform.WsConnection) {
	//done := make(chan bool)
	//go ws.keepAlive(ctx, c, done)
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			//since we are in the infinite loop for test purposes we have this code here to break the loop
			if err.Error() == "envTest" {
				break
			}
			ws.logger.Error2("can not read message", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "readMessage"),
			)
			err := c.Close()
			if err != nil {
				ws.logger.Error2("can not close socket ", err,
					zap.String("service", "binanceWebSocket"),
					zap.String("method", "readMessage"),
				)
				os.Exit(1)
			}
			//done <- true
			os.Exit(1)
		}

		if messageType == websocket.TextMessage {
			res := wsMessage{}
			err = json.Unmarshal(message, &res)
			if err != nil {
				ws.logger.Error2("can not marshaling ws message", err,
					zap.String("service", "binanceWebSocket"),
					zap.String("method", "readMessage"),
					zap.String("message", string(message)),
				)
				//return
				continue
			}
			if res.Data != nil {
				ws.processMessages(ctx, res)
			}

			continue
		}

		//if messageType == websocket.PongMessage {
		//	err = c.WriteMessage(websocket.PongMessage, []byte(""))
		//	ws.logger.Error("error in write pong message in pong", err)
		//	continue
		//}
		//
		//if messageType == websocket.PingMessage {
		//	err = c.WriteMessage(websocket.PongMessage, []byte(""))
		//	ws.logger.Error("error in write pong message in ping", err)
		//	continue
		//}

	}

}

//func (ws *binanceWs) keepAlive(ctx context.Context, c platform.WsConnection, done chan bool) {
//ticker := time.NewTicker(3 * time.Minute)
//for {
//select {
//case <-done:
//ticker.Stop()
//return
//case <-ticker.C:
//err := c.WriteMessage(websocket.PongMessage, []byte(""))
//if err != nil {
//ws.logger.Error("error in write pong message in pong", err)
//}
//err = c.WriteMessage(websocket.PingMessage, []byte(""))
//if err != nil {
//ws.logger.Error("error in write pong message in pong", err)
//}
//}
//}
//}

func (ws *binanceWs) createWsSubscribeRequest(streams []string) ([]byte, error) {
	var params []string
	//streams := getStreamNames()
	for _, binancePairName := range ws.PairNameMap {
		for _, stream := range streams {
			finalStreamName := strings.ToLower(binancePairName) + "@" + stream
			params = append(params, finalStreamName)
		}
	}

	r := wsSubscribeRequest{
		ID:     1,
		Method: "SUBSCRIBE",
		Params: params,
	}
	rBytes, err := json.Marshal(r)
	if err != nil {
		return []byte(""), err
	}
	return rBytes, nil
}

func (ws *binanceWs) processMessages(ctx context.Context, message wsMessage) {
	parts := strings.Split(message.Stream, "@")

	streamName := parts[1]
	switch streamName {
	case TradeStream:
		bts := tradeStream{}
		err := mapstructure.Decode(message.Data, &bts)
		if err != nil {
			ws.logger.Error2("can not process message of trade stream", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "pocessMessage"),
				zap.String("messageStream", message.Stream),
			)
			return
		}
		go ws.handleTradeStream(ctx, bts)
		break
	case DepthStream:
		bds := depthStream{}
		err := mapstructure.Decode(message.Data, &bds)
		if err != nil {
			ws.logger.Error2("can not process message of depth stream", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "pocessMessage"),
				zap.String("messageStream", message.Stream),
			)
			return
		}

		ws.handleDepthStream(ctx, bds)
		break
	case TickerStream:
		ts := tickerStream{}
		err := mapstructure.Decode(message.Data, &ts)
		if err != nil {
			ws.logger.Error2("can not process message of ticker stream", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "pocessMessage"),
				zap.String("messageStream", message.Stream),
			)
			return
		}
		go ws.handleTickerStream(ctx, ts)
		break
	default:
		ks := klineStream{}
		err := mapstructure.Decode(message.Data, &ks)
		if err != nil {
			ws.logger.Error2("can not process message of kline stream", err,
				zap.String("service", "binanceWebSocket"),
				zap.String("method", "pocessMessage"),
				zap.String("messageStream", message.Stream),
			)
			return
		}
		ws.handleKlineStream(ctx, ks)
		break

	}

}

func (ws *binanceWs) handleTradeStream(ctx context.Context, ts tradeStream) {
	tm := time.Unix(0, ts.Tt*int64(1000000))
	trade := processor.Trade{
		Pair:      getPairNameFromBinancePair(ws.PairNameMap, ts.S),
		Price:     ts.P,
		Amount:    ts.Q,
		CreatedAt: tm.Format("2006-01-02 15:04:05"),
		IsMaker:   ts.M,
		Ignore:    ts.I,
	}

	ws.processor.ProcessTrade(ctx, trade)

}

func (ws *binanceWs) handleDepthStream(ctx context.Context, ds depthStream) {
	// here we should using goroutine but since the rate of binance publish is high we encounter data race so for now
	// just handling it without goroutine
	data, err := json.Marshal(ds)
	if err != nil {
		ws.logger.Error2("can not marshal depth stream", err,
			zap.String("service", "binanceWebSocket"),
			zap.String("method", "handleDepthStream"),
		)
	}
	pairName := getPairNameFromBinancePair(ws.PairNameMap, ds.S)
	ws.processor.ProcessDepth(ctx, ExchangeName, pairName, ds.S, data)
}

func (ws *binanceWs) handleTickerStream(ctx context.Context, ts tickerStream) {
	pair := getPairNameFromBinancePair(ws.PairNameMap, ts.S)
	ticker := processor.Ticker{
		Pair:            pair,
		Price:           ts.C,
		Percentage:      ts.Pb,
		ID:              0,  //this field would be set in data processor
		EquivalentPrice: "", //this field would be set in data processor
		Volume:          ts.Q,
		High:            ts.H,
		Low:             ts.L,
	}

	ws.processor.ProcessTicker(ctx, ticker)

}

func (ws *binanceWs) handleKlineStream(ctx context.Context, ks klineStream) {
	pairName := getPairNameFromBinancePair(ws.PairNameMap, ks.K.S)
	tf := getTimeFrameFromBinanceTimeFrame(ks.K.I)
	startTime := time.Unix(0, ks.K.T*int64(1000000))
	closeTime := time.Unix(0, ks.K.Tb*int64(1000000))

	kline := processor.Kline{
		Pair:                pairName,
		TimeFrame:           tf,
		KlineStartTime:      startTime.Format("2006-01-02 15:04:05"),
		KlineCloseTime:      closeTime.Format("2006-01-02 15:04:05"),
		OpenPrice:           ks.K.O,
		ClosePrice:          ks.K.C,
		HighPrice:           ks.K.H,
		LowPrice:            ks.K.L,
		BaseVolume:          ks.K.V,
		QuoteVolume:         ks.K.Q,
		TakerBuyBaseVolume:  ks.K.Vb,
		TakerBuyQuoteVolume: ks.K.Qb,
	}

	ws.processor.ProcessKline(ctx, kline)

}

func createPairNameMap(activePairs []currency.Pair) map[string]string {

	pairNameMap := make(map[string]string)
	for _, pair := range activePairs {
		pairName := pair.Name
		binancePairName := strings.Replace(pairName, "-", "", -1)
		pairNameMap[pairName] = binancePairName
	}
	return pairNameMap
}

//get BTCUSDT return BTC-USDT
func getPairNameFromBinancePair(pairNameMap map[string]string, value string) string {
	for k, v := range pairNameMap {
		if v == value {
			return k
		}
	}
	return ""
}

func getStreamNames() []string {
	return []string{
		TradeStream,
		DepthStream,
		TickerStream,
		KlineStream1m,
		KlineStream5m,
		KlineStream1h,
		KlineStream1d,
	}
}

func getTimeFrameFromBinanceTimeFrame(binanceTimeFrame string) string {
	for k, v := range getTimeFrameMap() {
		if v == binanceTimeFrame {
			return k

		}
	}
	return ""
}

func getTimeFrameMap() map[string]string {
	return map[string]string{
		"1minute":  "1m",
		"5minutes": "5m",
		"1hour":    "1h",
		"1day":     "1d",
	}
}

func NewWs(wsClient platform.WsClient, processor processor.Processor, logger platform.Logger, activePairs []currency.Pair) types.ExternalWs {
	PairNameMap := createPairNameMap(activePairs)
	return &binanceWs{wsClient, processor, logger, activePairs, PairNameMap}
}
