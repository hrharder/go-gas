package gas

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
)

// ETHGasStationURL is the API URL for the ETH Gas Station API.
//
// More information available at https://ethgasstation.info
const ETHGasStationURL = "https://ethgasstation.info/json/ethgasAPI.json"

// GasPriority is a type alias for a string, with supported priorities included in this package.
type GasPriority string

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

// SuggestFastGasPrice is a helper method that calls SuggestGasPrice with GasPriorityFast
func SuggestFastGasPrice() (*big.Int, error) {
	return SuggestGasPrice(GasPriorityFast)
}

type ethGasStationResponse struct {
	Fast    json.Number `json:"fast"`
	Fastest json.Number `json:"fastest"`
	SafeLow json.Number `json:"safeLow"`
	Average json.Number `json:"average"`
}

func loadGasPrices() (*ethGasStationResponse, error) {
	res, err := http.Get(ETHGasStationURL)
	if err != nil {
		return nil, err
	}

	dcr := json.NewDecoder(res.Body)
	dcr.UseNumber()

	var body ethGasStationResponse
	if err := dcr.Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

// convert eth gas station units to wei
// (raw result / 10) * 1e9 = base units (wei)
func parseGasPriceToWei(raw json.Number) (*big.Int, error) {
	num, ok := new(big.Float).SetString(raw.String())
	if !ok {
		return nil, errors.New("eth: unable to parse float value")
	}

	gwei := new(big.Float).Mul(num, big.NewFloat(100000000))
	if !gwei.IsInt() {
		return nil, errors.New("eth: unable to represent gas price as integer")
	}

	var wei *big.Int
	wei, _ = gwei.Int(wei)
	return new(big.Int).Set(wei), nil
}
