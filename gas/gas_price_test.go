package gas

import (
	"encoding/json"
	"math/big"
	"testing"

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
	oneGweiInGasStationUnits := json.Number("10")
	oneGweiInBaseUnits := big.NewInt(int64(1e9))

	parsed, err := parseGasPriceToWei(oneGweiInGasStationUnits)
	require.NoError(t, err)

	assert.Equal(t, 0, oneGweiInBaseUnits.Cmp(parsed))
}

func TestLoadGasPrices(t *testing.T) {
	rawPrices, err := loadGasPrices()
	require.NoError(t, err)

	fast, err := rawPrices.Fast.Float64()
	require.NoError(t, err)
	fastest, err := rawPrices.Fastest.Float64()
	require.NoError(t, err)
	safeLow, err := rawPrices.SafeLow.Float64()
	require.NoError(t, err)
	average, err := rawPrices.Average.Float64()
	require.NoError(t, err)

	assert.GreaterOrEqual(t, fastest, fast)
	assert.GreaterOrEqual(t, fast, average)
	assert.GreaterOrEqual(t, average, safeLow)
	assert.GreaterOrEqual(t, safeLow, 0.0)
}
