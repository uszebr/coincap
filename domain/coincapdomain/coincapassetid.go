package coincapdomain

// represent obj from {{host}}api.coincap.io/v2/assets/bitcoin endpoint(can be various ids in the end of request)

// Struct to represent the JSON structure('data' with one entity)
type CoincapAssetIdResponse struct {
	Data      CoincapAssets `json:"data"`
	Timestamp int64         `json:"timestamp"`
}

// Example
// {
// 	"data": {
// 	  "id": "bitcoin",
// 	  "rank": "1",
// 	  "symbol": "BTC",
// 	  "name": "Bitcoin",
// 	  "supply": "17193925.0000000000000000",
// 	  "maxSupply": "21000000.0000000000000000",
// 	  "marketCapUsd": "119179791817.6740161068269075",
// 	  "volumeUsd24Hr": "2928356777.6066665425687196",
// 	  "priceUsd": "6931.5058555666618359",
// 	  "changePercent24Hr": "-0.8101417214350335",
// 	  "vwap24Hr": "7175.0663247679233209"
// 	},
// 	"timestamp": 1533581098863
//   }
