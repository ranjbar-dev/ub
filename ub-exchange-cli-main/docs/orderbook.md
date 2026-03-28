### orderbook
we provide the binance orderbook in our exchange.binance does not have a api or websocket for the orderbook and we should 
create the orderbook using depth stream in ws. all logic about creating orderbook is happening in orderbook package.
the following lines are the binance description to create order book.

1. Open a stream to wss://stream.binance.com:9443/ws/bnbbtc@depth.
2. Buffer the events you receive from the stream.
3. Get a depth snapshot from https://api.binance.com/api/v3/depth?symbol=BNBBTC&limit=1000 .
4. Drop any event where u is <= lastUpdateId in the snapshot.
5. The first processed event should have U <= lastUpdateId+1 AND u >= lastUpdateId+1.
6. While listening to the stream, each new event's U should be equal to the previous event's u+1.
7. The data in each event is the absolute quantity for a price level.
8. If the quantity is 0, remove the price level.
9. Receiving an event that removes a price level that is not in your local order book can happen and is normal.

##### how does orderbook service works
after receiving websocket data of depth ,we save it in redis, as stream data keep receiving
we try to update the orderbook data in redis , somehow some data would be out of date which should be removed.
the out of date data are the data that lastUpdatedId is less than the current receiving one , in this situation we get 
depth data from binance api and try to rewrite the orderbook.finally we sort the data and send it to clients using MQTT

