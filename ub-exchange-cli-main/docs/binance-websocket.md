## live data from binance 

the data we get includes kline,ticker,trade book, depth (order book)

the binance documentation could be found [here](https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md)

### how we handle the data from binance websocket 

since the number of connection to binance web socket is limited 
we only have 3 connection
 * one for depth(order book)
 * one for ticker and trade book
 * one for kline (1m 5m 1h 1d)
which all handled using supervisord service (the config file is in our ub-server project `.docker/supervisor/go-supervisord`)
 
#### how the code works 

first we connect to binance web socket for our supported pairs in file  `internal/externalexchangews/binance/ws`

after receiving data for each stream we should map and filter data to what is more suitable for us

mapping and filtering happens in  `internal/processor/dataprocessor`

for each stream we do different things and finally publish data to our clients using MQTT protocol

#### ticker stream
- for ticker stream we save data in redis including price, percentage and volume), the ticker also gives us the signal 
which one of our stop orders or limit orders that are in queue should be submitted, this happen in 
orderSubmissionManager for limit orders and in stopOrderSubmissionManager for stop orders.

data in ticker includes these fields

filed name | type | description
---------- | ---- | -----------
name |  string | which is pair name like "BTC-USDT" 
price | string | which is pair price like "35000.23"
percentage  | string  | which is pair percentage change comparison to last 24 hour price like "12.23"
id   | int64 | which is pair id in our database
equivalentPrice  | string | which is dollar equivalent of price like "32.43"
volume  | string | which is last 24 hour traded volume like "1260000.43"
high   | string | which is last 24 hour highest price
low    | string | which is last 24 lowest highest price


#### trade book stream
for trade book stream  we only save in redis 

data in trade book includes these fields

filed name | type | description
---------- | ---- | -----------
pair | string | which is pair name like "BTC-USDT"
price | string | which is trade price the orders are matched like "34000.32"
amount | string | which is trade amount like "0.032" 
createdAt | string | the time trade is created like "2020-12-12 20:20:20"
isMaker | bool | is the maker side the buy one
ignore | bool|


#### depth stream
for depth stream we create our order book which is the binance order book itself, since binance has not any stream for order book
we have to create the order book data from depth stream which all happens in orderbook package
please read the `orderbook.md` for details how we handle the orderbook

 

#### kline stream
for kline stream we save data in redis as pre_kline, then in any change of time frame we get data from redis and persist it in database in ohlc table
data in kline includes these fields

filed name | type | description
---------- | ---- | -----------
pair | string  | which is pair name like "BTC-USDT"
timeFrame | string | which includes 1minutes,5minutes,1hour,1day
startTime | string | the time candle is started like "2020-12-12 20:20:00"
closeTime | string | the time candle is closed like "2020-12-12 20:21:00"
openPrice | string | the price in which candle is opened 
closePrice | string | the price in which  candle is closed
highPrice | string | the high price in candle
lowPrice | string  | the low price in candle
baseVolume | string | the aggregated amount which are traded in base coin
quoteVolume | string | the aggregated amount which are traded in quote coin 
takerBuyBaseVolume | string | 
takerBuyQuoteVolume | string | 

considering some points about binance web socket
* we should send pong message in every ping of binance , if we do not the binance will drop the connection
* number of stream  is limited in each connection (maximum 1024)

