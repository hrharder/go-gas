package gas

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// ETHGasStationURL is the API URL for the ETH Gas Station API.
//
// More information available at https://ethgasstation.info
const ETHGasStationURL = "https://ethgasstation.info/json/ethgasAPI.json"

// GasPriority is a type alias for a string, with supported priorities included in this package.
type GasPriority string

// GasPriceSuggester is type alias  for a function that returns a reccomended gas price in base units for a given priority level.
type GasPriceSuggester func(GasPriority) (*big.Int, error)

const (
	// GasPriorityFast is the recommended gas price for a transaction to be mined in less than 2 minutes.
	GasPriorityFast = GasPriority("fast")

	// GasPriorityFastest is the recommended gas price for a transaction to be mined in less than 30 seconds.
	GasPriorityFastest = GasPriority("fastest")

	// GasPrioritySafeLow is the recommended cheapest gas price for a transaction to be mined in less than 30 minutes.
	GasPrioritySafeLow = GasPriority("safeLow")

	// GasPriorityAverage is the recommended average gas price for a transaction to be mined in less than 5 minutes.
	GasPriorityAverage = GasPriority("average")
)

// SuggestGasPrice returns a suggested gas price value in wei (base units) for timely transaction execution.
//
// The returned price depends on the priority specified, and supports all priorities supported by the ETH Gas Station API.
func SuggestGasPrice(priority GasPriority) (*big.Int, error) {
	prices, err := loadGasPrices()
	if err != nil {
		return nil, err
	}
	return parseSuggestedGasPrice(priority, prices)
}

// SuggestFastGasPrice is a helper method that calls SuggestGasPrice with GasPriorityFast
func SuggestFastGasPrice() (*big.Int, error) {
	return SuggestGasPrice(GasPriorityFast)
}

// NewGasPriceSuggester returns a function that can be used to either load a new gas price response, or use a cached
// response if it is within the age range defined by maxResultAge.
//
// The returned function loads from the cache or pulls a new response if the stored result is older than maxResultAge.
func NewGasPriceSuggester(maxResultAge time.Duration) (GasPriceSuggester, error) {
	prices, err := loadGasPrices()
	if err != nil {
		return nil, err
	}
	m := gasPriceManager{
		latestResponse: prices,
		fetchedAt:      time.Now(),
		maxResultAge:   maxResultAge,
	}
	return func(priority GasPriority) (*big.Int, error) {
		return m.suggestCachedGasPrice(priority)
	}, nil
}

type gasPriceManager struct {
	latestResponse *ethGasStationResponse
	fetchedAt      time.Time
	maxResultAge   time.Duration

	mux sync.Mutex
}

func (m *gasPriceManager) suggestCachedGasPrice(priority GasPriority) (*big.Int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// fetch new values if stored result is older than the maximum age
	if time.Since(m.fetchedAt) > m.maxResultAge {
		prices, err := loadGasPrices()
		if err != nil {
			return nil, err
		}
		m.latestResponse = prices
		m.fetchedAt = time.Now()
	}

	return parseSuggestedGasPrice(priority, m.latestResponse)
}

// conversion factor to go from (gwei * 10) to wei
// equal to: (raw / 10) => gwei => gwei * 1e9 => wei
// simplifies to: raw * 1e8 => wei
var conversionFactor = big.NewFloat(100000000)

type ethGasStationResponse struct {
	Fast    float64 `json:"fast"`
	Fastest float64 `json:"fastest"`
	SafeLow float64 `json:"safeLow"`
	Average float64 `json:"average"`
}

func loadGasPrices() (*ethGasStationResponse, error) {
	res, err := http.Get(ETHGasStationURL)
	if err != nil {
		return nil, err
	}

	var body ethGasStationResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func parseSuggestedGasPrice(priority GasPriority, prices *ethGasStationResponse) (*big.Int, error) {
	switch priority {
	case GasPriorityFast:
		return parseGasPriceToWei(prices.Fast)
	case GasPriorityFastest:
		return parseGasPriceToWei(prices.Fastest)
	case GasPrioritySafeLow:
		return parseGasPriceToWei(prices.SafeLow)
	case GasPriorityAverage:
		return parseGasPriceToWei(prices.Average)
	default:
		return nil, errors.New("eth: unknown/unsupported gas priority")
	}
}

// convert eth gas station units to wei
// (raw result / 10) * 1e9 = base units (wei)
func parseGasPriceToWei(raw float64) (*big.Int, error) {
	gwei := new(big.Float).Mul(big.NewFloat(raw), conversionFactor)
	if !gwei.IsInt() {
		return nil, errors.New("eth: unable to represent gas price as integer")
	}

	// we can skip the accuracy check because we know from above that gwei is an integer
	wei, _ := gwei.Int(new(big.Int))
	return wei, nil
}
