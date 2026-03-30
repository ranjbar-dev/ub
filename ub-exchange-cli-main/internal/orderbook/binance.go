package orderbook

import (
	"context"
	"encoding/json"
	"exchange-go/internal/livedata"
	"exchange-go/internal/platform"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const DepthURL = "https://www.binance.com/api/v3/depth?limit=100&symbol="
const FetchingPeriod = 1800 // 30 minutes

type binanceOrderBook struct {
	httpClient          platform.HTTPClient
	liveDataService     livedata.Service
	lastBinanceAPIFetch int64
	mu                  sync.Mutex
}

//this is data from websocket
type binanceDepth struct {
	FirstUpdatedID int64       `json:"Ub"` // First update ID in event
	LastUpdatedID  int64       `json:"U"`  // Final update ID in event
	Bids           [][2]string `json:"B"`  // Bids to be updated
	Asks           [][2]string `json:"A"`  // Asks to be updated
}

//this is data from api
type depthSnapshot struct {
	LastUpdatedID int64       `json:"lastUpdateId"`
	Bids          [][2]string `json:"bids"`
	Asks          [][2]string `json:"asks"`
}

func (bo *binanceOrderBook) updateOrderBook(ctx context.Context, pairName string, externalExchangePairName string, depth binanceDepth) (RawOrderBook, error) {
	r := RawOrderBook{}
	firstUpdatedID := depth.FirstUpdatedID
	redisDepthSnapshot, err := bo.getDepthSnapshotFromRedis(ctx, pairName)

	if err != nil && err != redis.Nil {
		return r, err
	}

	handleInvalidData(&redisDepthSnapshot)

	if firstUpdatedID-redisDepthSnapshot.LastUpdatedID != 1 && len(redisDepthSnapshot.Asks) > 0 {
		apiDepthSnapshot, err := bo.fetchDepthFromBinance(ctx, externalExchangePairName)
		if err != nil {
			return r, err
		}

		cleanOutOfDateOrders(redisDepthSnapshot, apiDepthSnapshot.LastUpdatedID)
		attachAPIDepthSnapshot(&redisDepthSnapshot, apiDepthSnapshot)
		sortSnapshot(&redisDepthSnapshot)
		err = bo.setDepthSnapshot(ctx, pairName, redisDepthSnapshot)
		if err != nil {
			return r, err
		}
	}

	updateRedisDepthSnapshot(&redisDepthSnapshot, depth)
	sortSnapshot(&redisDepthSnapshot)
	err = bo.setDepthSnapshot(ctx, pairName, redisDepthSnapshot)
	if err != nil {
		return r, nil
	}

	r = RawOrderBook{
		Bids: redisDepthSnapshot.Bids,
		Asks: redisDepthSnapshot.Asks,
	}
	return r, nil

}

func (bo *binanceOrderBook) fetchDepthFromBinance(ctx context.Context, externalExchangePairName string) (depthSnapshot, error) {
	depthSnapshot := depthSnapshot{}
	//fetch from binance again since the fetching period duration is passed
	//if 2000 millisecond or 2 second has passed from last fetch then we fetch
	bo.mu.Lock()
	defer bo.mu.Unlock()
	lastFetch := bo.lastBinanceAPIFetch
	now := time.Now().Unix()
	if lastFetch == 0 || now-lastFetch >= 2 {
		url := DepthURL + externalExchangePairName
		bo.lastBinanceAPIFetch = now
		headers := make(map[string]string)
		body, _, statusCode, err := bo.httpClient.HTTPGet(ctx, url, headers)
		if err != nil {
			return depthSnapshot, err
		}

		if statusCode != http.StatusOK {
			return depthSnapshot, fmt.Errorf("status code is not 200 it is %d with body %s", statusCode, string(body))
		}
		err = json.Unmarshal([]byte(body), &depthSnapshot)
		if err != nil {
			return depthSnapshot, nil
		}
	}

	return depthSnapshot, nil
}

func (bo *binanceOrderBook) setDepthSnapshot(ctx context.Context, pairName string, depthSnapshot livedata.RedisDepthSnapshot) error {
	bids := depthSnapshot.Bids
	if len(bids) > 1000 {
		bids = bids[:1000]
	}
	asks := depthSnapshot.Asks

	if len(asks) > 1000 {
		asks = asks[:1000]
	}

	dp := livedata.RedisDepthSnapshot{
		LastUpdatedID: depthSnapshot.LastUpdatedID,
		UpdatedAt:     time.Now().Unix(),
		Bids:          bids,
		Asks:          asks,
	}
	return bo.liveDataService.SetDepthSnapshot(ctx, pairName, dp)
}

func (bo *binanceOrderBook) getDepthSnapshotFromRedis(ctx context.Context, pairName string) (livedata.RedisDepthSnapshot, error) {
	return bo.liveDataService.GetDepthSnapshot(ctx, pairName)
}

/**
removing the data of bids and asks which their lastUpdatedId is less than the receiving one from depth stream
*/
func cleanOutOfDateOrders(redisDepthSnapshot livedata.RedisDepthSnapshot, lastUpdatedID int64) {
	if lastUpdatedID == 0 {
		return
	}
	for i := 0; i < len(redisDepthSnapshot.Bids); i++ {
		existingLastUpdatedID, _ := strconv.ParseInt(redisDepthSnapshot.Bids[i][2], 10, 64)
		if existingLastUpdatedID < lastUpdatedID {
			redisDepthSnapshot.Bids = append(redisDepthSnapshot.Bids[:i], redisDepthSnapshot.Bids[i+1:]...)
			i-- // Since we just deleted redisDepthSnapshot.Bids[i], we must redo that index
		}
	}

	for i := 0; i < len(redisDepthSnapshot.Asks); i++ {
		existingLastUpdatedID, _ := strconv.ParseInt(redisDepthSnapshot.Asks[i][2], 10, 64)
		if existingLastUpdatedID < lastUpdatedID {
			redisDepthSnapshot.Asks = append(redisDepthSnapshot.Asks[:i], redisDepthSnapshot.Asks[i+1:]...)
			i-- // Since we just deleted redisDepthSnapshot.Asks[i], we must redo that index
		}
	}
}

/*
	if the price receiving from websocket is already exists we update the amount ,
	if the amount is 0 we remove the price level,if it does not exists we attach it
*/
func attachAPIDepthSnapshot(redisDepthSnapshot *livedata.RedisDepthSnapshot, apiDepthSnapshot depthSnapshot) {
	lastUpdatedIDString := strconv.FormatInt(apiDepthSnapshot.LastUpdatedID, 10)
	for _, bid := range apiDepthSnapshot.Bids {
		exists := false
		for i, currentBid := range redisDepthSnapshot.Bids {
			if currentBid[0] == bid[0] {
				exists = true
				redisDepthSnapshot.Bids[i][1] = bid[1]
				redisDepthSnapshot.Bids[i][2] = lastUpdatedIDString
				if bid[1] == "0.00000000" {
					redisDepthSnapshot.Bids[i] = redisDepthSnapshot.Bids[len(redisDepthSnapshot.Bids)-1]
					redisDepthSnapshot.Bids = redisDepthSnapshot.Bids[:len(redisDepthSnapshot.Bids)-1]
				}
				break
			}
		}
		if !exists {
			redisDepthSnapshot.Bids = append(redisDepthSnapshot.Bids, [3]string{bid[0], bid[1], lastUpdatedIDString})
		}

	}

	for _, ask := range apiDepthSnapshot.Asks {
		exists := false
		for i, currentAsk := range redisDepthSnapshot.Asks {
			if currentAsk[0] == ask[0] {
				exists = true
				redisDepthSnapshot.Asks[i][1] = ask[1]
				redisDepthSnapshot.Asks[i][2] = lastUpdatedIDString
				if ask[1] == "0.00000000" {
					redisDepthSnapshot.Asks[i] = redisDepthSnapshot.Asks[len(redisDepthSnapshot.Asks)-1]
					redisDepthSnapshot.Asks = redisDepthSnapshot.Asks[:len(redisDepthSnapshot.Asks)-1]
				}
				break
			}
		}

		if !exists {
			redisDepthSnapshot.Asks = append(redisDepthSnapshot.Asks, [3]string{ask[0], ask[1], lastUpdatedIDString})
		}
	}
	redisDepthSnapshot.LastUpdatedID = apiDepthSnapshot.LastUpdatedID

}

func updateRedisDepthSnapshot(redisDepthSnapshot *livedata.RedisDepthSnapshot, depth binanceDepth) {
	lastUpdatedID := depth.LastUpdatedID
	lastUpdatedIDString := strconv.FormatInt(lastUpdatedID, 10)

	for _, bid := range depth.Bids {
		exists := false
		for i, currentBid := range redisDepthSnapshot.Bids {
			if currentBid[0] == bid[0] {
				exists = true
				redisDepthSnapshot.Bids[i][1] = bid[1]
				redisDepthSnapshot.Bids[i][2] = lastUpdatedIDString
				if bid[1] == "0.00000000" {
					redisDepthSnapshot.Bids[i] = redisDepthSnapshot.Bids[len(redisDepthSnapshot.Bids)-1]
					redisDepthSnapshot.Bids = redisDepthSnapshot.Bids[:len(redisDepthSnapshot.Bids)-1]
				}
				break
			}
		}
		if !exists {
			if bid[1] != "0.00000000" {
				redisDepthSnapshot.Bids = append(redisDepthSnapshot.Bids, [3]string{bid[0], bid[1], lastUpdatedIDString})
			}
		}
	}

	for _, ask := range depth.Asks {
		exists := false
		for i, currentAsk := range redisDepthSnapshot.Asks {
			if currentAsk[0] == ask[0] {
				exists = true
				redisDepthSnapshot.Asks[i][1] = ask[1]
				redisDepthSnapshot.Asks[i][2] = lastUpdatedIDString
				if ask[1] == "0.00000000" {
					redisDepthSnapshot.Asks[i] = redisDepthSnapshot.Asks[len(redisDepthSnapshot.Asks)-1]
					redisDepthSnapshot.Asks = redisDepthSnapshot.Asks[:len(redisDepthSnapshot.Asks)-1]
				}
				break
			}
		}
		if !exists {
			if ask[1] != "0.00000000" {
				redisDepthSnapshot.Asks = append(redisDepthSnapshot.Asks, [3]string{ask[0], ask[1], lastUpdatedIDString})
			}
		}

	}

	redisDepthSnapshot.LastUpdatedID = lastUpdatedID

}

func sortSnapshot(redisDepthSnapshot *livedata.RedisDepthSnapshot) {

	sort.Slice(redisDepthSnapshot.Bids, func(i, j int) bool {
		first, _ := strconv.ParseFloat(redisDepthSnapshot.Bids[i][0], 64)
		second, _ := strconv.ParseFloat(redisDepthSnapshot.Bids[j][0], 64)
		return first > second
	})

	sort.SliceStable(redisDepthSnapshot.Asks, func(i, j int) bool {
		first, _ := strconv.ParseFloat(redisDepthSnapshot.Asks[i][0], 64)
		second, _ := strconv.ParseFloat(redisDepthSnapshot.Asks[j][0], 64)
		return first < second
	})
}

func handleInvalidData(redisDepthSnapshot *livedata.RedisDepthSnapshot) {
	//when the process is stopped and we have data in redis new incoming data are invalid
	//here we try to empty the depth
	now := time.Now().Unix()
	if redisDepthSnapshot.LastUpdatedID != 0 && now-redisDepthSnapshot.UpdatedAt > 2 {
		if len(redisDepthSnapshot.Bids) > 50 && len(redisDepthSnapshot.Asks) > 50 {
			redisDepthSnapshot.LastUpdatedID = 0
			redisDepthSnapshot.Bids = [][3]string{}
			redisDepthSnapshot.Asks = [][3]string{}
		}
	}
}

func getBinanceOrderBook(httpClient platform.HTTPClient, liveDataService livedata.Service) externalExchangeOrderBook {
	return &binanceOrderBook{
		httpClient:          httpClient,
		liveDataService:     liveDataService,
		lastBinanceAPIFetch: 0,
	}
}
