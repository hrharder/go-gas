package gas

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuggestGasPrice(t *testing.T) {
	// 1. normal case, check that the value is greater or equal to 1 gwei and no error
	// - it is hard to test against a raw conversion due to the API updating frequently
	gasPrice, err := SuggestGasPrice(GasPriorityFast)
	require.NoError(t, err)
	baseline := big.NewInt(int64(1000000000))
	assert.GreaterOrEqual(t, gasPrice.Cmp(baseline), 0)

	// 2. check that we handle invalid priority selections
	_, err = SuggestGasPrice(GasPriority("foo"))
	assert.Error(t, err)
}

func TestParseGasPriceToWei(t *testing.T) {
	oneGweiInGasStationUnits := 10.0
	oneGweiInBaseUnits := big.NewInt(int64(1e9))

	parsed, err := parseGasPriceToWei(oneGweiInGasStationUnits)
	require.NoError(t, err)

	assert.Equal(t, 0, oneGweiInBaseUnits.Cmp(parsed))
}

func TestLoadGasPrices(t *testing.T) {
	rawPrices, err := loadGasPrices()
	require.NoError(t, err)

	require.GreaterOrEqual(t, rawPrices.Fastest, rawPrices.Fast)
	require.GreaterOrEqual(t, rawPrices.Fast, rawPrices.Average)
	require.GreaterOrEqual(t, rawPrices.Average, rawPrices.SafeLow)
	require.GreaterOrEqual(t, rawPrices.SafeLow, 0.0)
}

func TestGasPriceManager(t *testing.T) {
	// create "phony" negative price result so we know the cache is being used
	prices := &ethGasStationResponse{
		Fast:    -1.0,
		Fastest: -1.0,
		SafeLow: -1.0,
		Average: -1.0,
	}

	mgr := gasPriceManager{
		latestResponse: prices,
		fetchedAt:      time.Now(),
		maxResultAge:   50 * time.Millisecond,
	}

	// 1. should use a cached result up til duration has passed
	// - we can ensure a cached result is used by manually setting a cached result as -1
	cachedResult, err := mgr.suggestCachedGasPrice(GasPriorityFast)
	require.NoError(t, err)
	assert.Equal(t, "-100000000", cachedResult.String(), "cached result should be negative since we manually set the result")

	// 2. should fetch a new result after duration has passed
	time.Sleep(51 * time.Millisecond)
	newResult, err := mgr.suggestCachedGasPrice(GasPriorityFast)
	require.NoError(t, err)
	assert.Equal(t, newResult.Cmp(big.NewInt(0)), 1, "new result should be greater than 0")
}
