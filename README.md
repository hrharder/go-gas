# Module: `go-gas`

A simple Golang client for the [ETH Gas Station](https://ethgasstation.info) API.

It provides a simple set of methods for loading appropriate gas prices for submitting Ethereum transactions in base units.

Built only on the standard library, it has no external dependencies other than `stretchr/testify` which is used for testing.

# Usage

## Install

The `go-gas` module supports Go modules. Add to your project with the following command.

```
go get -u github.com/hrharder/go-gas
```

## Usage

Package `gas` provides two main ways to fetch a gas price from the ETH Gas Station API.

1. Fetch the current recommended price for a given priority level with a new API call each time
   - Use `gas.SuggestGasPrice` for a specific priority level
   - Use `gas.SuggestFastGasPrice` to fetch the fast priority level (no arguments)
1. Create a new `GasPriceSuggester` which maintains a cache of results for a user-defined duration
   - Use `gas.NewGasPriceSuggester` and specify a max result age
   - Use the returned function to fetch new gas prices, or use the cache based on how old the results are


### Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/hrharder/go-gas"
)

func main() {
    // get a gas price in base units with one of the exported priorities (fast, fastest, safeLow, average)
    fastestGasPrice, err := gas.SuggestGasPrice(gas.GasPriorityFastest)
    if err != nil {
        log.Fatal(err)
    }

    // convenience wrapper for getting the fast gas price
    fastGasPrice, err := gas.SuggestFastGasPrice()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(fastestGasPrice)
    fmt.Println(fastGasPrice)

    // alternatively, use the NewGasPriceSuggester which maintains a cache of results until they are older than max age
    suggestGasPrice, err := gas.NewGasPriceSuggester(5 * time.Minute)
    if err != nil {
        log.Fatal(err)
    }

    fastGasPriceFromCache, err := suggestGasPrice(gas.GasPriorityFast)
    if err != nil {
        return nil, err
    }

    // after 5 minutes, the cache will be invalidated and new results will be fetched
    time.Sleep(5 * time.Minute)
    fasGasPriceFromAPI, err := suggestGasPrice(gas.GasPriorityFast)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(fastGasPriceFromCache)
    fmt.Println(fasGasPriceFromAPI)
}
```
