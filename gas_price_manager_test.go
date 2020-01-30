package gas

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestGasPriceManager(t *testing.T) {
	mgr := NewManager(100 * time.Millisecond)

	// 1. should use a cached result up til duration has passed
	// - we can ensure a cached result is used by manually setting a cached result as -1
	mgr.lastPrices[GasPriorityFast] = priceResult{
		fetchedAt: time.Now(),
		price:     big.NewInt(-1),
	}
	price1, err := mgr.SuggestGasPrice(GasPriorityFast)
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(-1), price1)

	// 2. should fetch a new result after duration has passed
	time.Sleep(101 * time.Millisecond)
	price2, err := mgr.SuggestGasPrice(GasPriorityFast)
	require.NoError(t, err)
	assert.NotEqual(t, big.NewInt(-1), price2)
	assert.Equal(t, price2.Cmp(big.NewInt(0)), 1)
}
