package coincapdomain

// represent obj from api.coincap.io/v2/assets/bitcoin/markets
type CoincapMarkets struct {
	ExchangeId    string  `json:"exchangeId"`
	BaseId        string  `json:"baseId"`
	QuoteId       string  `json:"quoteId"`
	BaseSymbol    string  `json:"baseSymbol"`
	QuoteSymbol   string  `json:"quoteSymbol"`
	VolumeUsd24Hr float64 `json:"volumeUsd24Hr,string"` // Parse as float64
	PriceUsd      float64 `json:"priceUsd,string"`      // Parse as float64
	VolumePercent float64 `json:"volumePercent,string"` // Parse as float64
}

// Struct to represent the outer JSON structure(map sith 'data' key)
type CoincapMarketsResponse struct {
	Data      []CoincapMarkets `json:"data"`
	Timestamp int64            `json:"timestamp"`
}

// example:
// {
// 	"data": [
// 	  {
// 		"exchangeId": "Binance",
// 		"baseId": "bitcoin",
// 		"quoteId": "tether",
// 		"baseSymbol": "BTC",
// 		"quoteSymbol": "USDT",
// 		"volumeUsd24Hr": "277775213.1923032624064566",
// 		"priceUsd": "6263.8645034633024446",
// 		"volumePercent": "7.4239157877678087"
// 	  },
// 	  {
// 		"exchangeId": "Bitfinex",
// 		"baseId": "bitcoin",
// 		"quoteId": "united-states-dollar",
// 		"baseSymbol": "BTC",
// 		"quoteSymbol": "USD",
// 		"volumeUsd24Hr": "246147462.2180977690000000",
// 		"priceUsd": "6280.1000000000000000",
// 		"volumePercent": "6.5786216483427765"
// 	  },
