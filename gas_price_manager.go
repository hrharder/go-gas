package gas

import (
	"math/big"
	"sync"
	"time"
)

type priceResult struct {
	fetchedAt time.Time
	price     *big.Int
}

// Manager is a gas price fetcher with a cache of recent results, and user-defined max age for those results.
//
// A Manager instance's cache is per priority level. It also has an default priority to simplify usage.
type Manager struct {
	defaultPriority GasPriority
	lastPrices      map[GasPriority]priceResult
	maxAge          time.Duration

	mu sync.Mutex
}

// NewManager creates a new gas price manager, with the default priority set to "fast" (~1 min confirmation).
func NewManager(maxResultAge time.Duration) *Manager {
	return &Manager{
		defaultPriority: GasPriorityFast,
		lastPrices:      make(map[GasPriority]priceResult),
		maxAge:          maxResultAge,
	}
}

// NewManager creates a new gas price manager, with the default priority set to the user defined defaultPriority
func NewManagerWithDefault(maxResultAge time.Duration, defaultPriority GasPriority) *Manager {
	return &Manager{
		defaultPriority: defaultPriority,
		lastPrices:      make(map[GasPriority]priceResult),
		maxAge:          maxResultAge,
	}
}

// SuggestGasPrice either fetches a new gas price, or uses the stored result if the most recent cache is present and not too old.
func (mgr *Manager) SuggestGasPrice(priority GasPriority) (*big.Int, error) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	// if we have a result saved for this priority level, and it's not to old, use it
	if res, ok := mgr.lastPrices[priority]; ok {
		if time.Since(res.fetchedAt) <= mgr.maxAge {
			return res.price, nil
		}
	}

	// if not, fetch the price for this priority level, store it, and return it
	price, err := SuggestGasPrice(priority)
	if err != nil {
		return nil, err
	}

	mgr.lastPrices[priority] = priceResult{
		fetchedAt: time.Now(),
		price:     price,
	}

	return price, nil
}

// SuggestDefaultGasPrice behaves like SuggestGasPrice, instead using the default priority level defined upon Manager construction.
func (mgr *Manager) SuggestDefaultGasPrice() (*big.Int, error) {
	return mgr.SuggestGasPrice(mgr.defaultPriority)
}
