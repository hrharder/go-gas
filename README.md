# Module: `go-gas`

A simple Golang wrapper for the [ETH Gas Station](https://ethgasstation.info) API.

It provides a simple set of methods for loading appropriate gas prices for submitting Ethereum transactions in base units.

Built only on the standard library, it has no external dependencies other than `stretchr/testify` which is used for testing.

# Usage

## Install

The `go-gas` module supports Go modules. Add to your project with the following command.

```
go get -u github.com/hrharder/go-gas
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/hrharder/go-gas/gas"
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
}
```
