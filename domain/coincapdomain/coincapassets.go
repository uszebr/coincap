package coincapdomain

// represent obj from {{host}}/v2/assets endpoint
type CoincapAssets struct {
	ID                string   `json:"id"`
	Rank              string   `json:"rank"`
	Symbol            string   `json:"symbol"`
	Name              string   `json:"name"`
	Supply            float64  `json:"supply,string"`            // Parse as float64
	MaxSupply         *float64 `json:"maxSupply,string"`         // Nullable and float64
	MarketCapUsd      float64  `json:"marketCapUsd,string"`      // Parse as float64
	VolumeUsd24Hr     float64  `json:"volumeUsd24Hr,string"`     // Parse as float64
	PriceUsd          float64  `json:"priceUsd,string"`          // Parse as float64
	ChangePercent24Hr float64  `json:"changePercent24Hr,string"` // Parse as float64
	Vwap24Hr          float64  `json:"vwap24Hr,string"`          // Parse as float64
	Explorer          string   `json:"explorer"`
}

// Struct to represent the outer JSON structure(map sith 'data' key)
type CoincapAssetsResponse struct {
	Data      []CoincapAssets `json:"data"`
	Timestamp int64           `json:"timestamp"`
}

// Example:
// {
// 	"data": [
// 	  {
// 		"id": "bitcoin",
// 		"rank": "1",
// 		"symbol": "BTC",
// 		"name": "Bitcoin",
// 		"supply": "17193925.0000000000000000",
// 		"maxSupply": "21000000.0000000000000000",
// 		"marketCapUsd": "119150835874.4699281625807300",
// 		"volumeUsd24Hr": "2927959461.1750323310959460",
// 		"priceUsd": "6929.8217756835584756",
// 		"changePercent24Hr": "-0.8101417214350335",
// 		"vwap24Hr": "7175.0663247679233209"
// 	  },
// 	  {
// 		"id": "ethereum",
// 		"rank": "2",
// 		"symbol": "ETH",
// 		"name": "Ethereum",
// 		"supply": "101160540.0000000000000000",
// 		"maxSupply": null,
// 		"marketCapUsd": "40967739219.6612727047843840",
// 		"volumeUsd24Hr": "1026669440.6451482672850841",
// 		"priceUsd": "404.9774667045200896",
// 		"changePercent24Hr": "-0.0999626159535347",
// 		"vwap24Hr": "415.3288028454417241"
// 	  },
