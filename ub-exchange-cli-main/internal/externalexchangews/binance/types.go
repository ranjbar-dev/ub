package binance

type wsSubscribeRequest struct {
	ID     int      `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type wsMessage struct {
	Stream string                 `json:"stream"`
	Data   map[string]interface{} `json:"data"`
}

type tradeStream struct {
	Et string `mapstructure:"e"` // Event type
	E  int64  `mapstructure:"E"` // Event time
	S  string `mapstructure:"s"` // Symbol
	T  int64  `mapstructure:"t"` // Trade ID
	P  string `mapstructure:"p"` // price
	Q  string `mapstructure:"q"` // Quantity
	B  int64  `mapstructure:"b"` // Buyer order ID
	A  int64  `mapstructure:"a"` // Seller order ID
	Tt int64  `mapstructure:"T"` // Trade time
	M  bool   `mapstructure:"m"` // Is the buyer the market maker?
	I  bool   `mapstructure:"M"` // Ignore
}

type depthStream struct {
	Et string     `mapstructure:"e"` // Event type
	E  int64      `mapstructure:"E"` // Event time
	S  string     `mapstructure:"s"` // Symbol
	Ub int64      `mapstructure:"U"` // First update ID in event
	U  int64      `mapstructure:"u"` // Final update ID in event
	B  [][]string `mapstructure:"b"` // Bids to be updated
	A  [][]string `mapstructure:"a"` // Asks to be updated
}

type tickerStream struct {
	Et string `mapstructure:"e"` // Event type
	E  int64  `mapstructure:"E"` // Event time
	S  string `mapstructure:"s"` // Symbol
	P  string `mapstructure:"p"` // price change
	Pb string `mapstructure:"P"` // price change percent
	W  string `mapstructure:"w"` // Weighted average price
	X  string `mapstructure:"x"` // First trade(F)-1 price (first trade before the 24hr rolling window)
	C  string `mapstructure:"c"` // Last price
	Qb string `mapstructure:"Q"` // Last quantity
	B  string `mapstructure:"b"` // Best bid price
	Bb string `mapstructure:"B"` // Best bid quantity
	A  string `mapstructure:"a"` // Best ask price
	Ab string `mapstructure:"A"` // Best ask quantity
	O  string `mapstructure:"o"` // Open price
	H  string `mapstructure:"h"` // High price
	L  string `mapstructure:"l"` // Low price
	V  string `mapstructure:"v"` // Total traded base asset volume
	Q  string `mapstructure:"q"` // Total traded quote asset volume
	Ob int64  `mapstructure:"O"` // Statistics open time
	Cb int64  `mapstructure:"C"` // Statistics close time
	F  int64  `mapstructure:"F"` // First trade ID
	Lb int64  `mapstructure:"L"` // Last trade Id
	N  int64  `mapstructure:"n"` // Total number of trades
}

type kline struct {
	T  int64  `mapstructure:"t"` // Kline start time
	Tb int64  `mapstructure:"T"` // Kline close time
	S  string `mapstructure:"s"` // Symbol
	I  string `mapstructure:"i"` // Interval
	F  int64  `mapstructure:"f"` // First trade ID
	Lb int64  `mapstructure:"L"` // Last trade ID
	O  string `mapstructure:"o"` // Open price
	C  string `mapstructure:"c"` // Close price
	H  string `mapstructure:"h"` // High price
	L  string `mapstructure:"l"` // Low price
	V  string `mapstructure:"v"` // Base asset volume
	N  int64  `mapstructure:"n"` // Number of trades
	//X  bool   `mapstructure:"x"` // Is this kline closed?
	Q  string `mapstructure:"q"` // Quote asset volume
	Vb string `mapstructure:"V"` // Taker buy base asset volume
	Qb string `mapstructure:"Q"` // Taker buy quote asset volume
	//B  string `mapstructure:"B"` // Ignore
}

type klineStream struct {
	Et string `mapstructure:"e"` // Event type
	E  int64  `mapstructure:"E"` // Event time
	S  string `mapstructure:"s"` // Symbol
	K  kline  `mapstructure:"k"` //kline
}

type timeframeMap map[string]string
