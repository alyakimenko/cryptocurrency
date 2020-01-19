![GitHub Actions](https://github.com/jollheef/cryptocurrency/workflows/Build%20and%20Test/badge.svg)
[![GoDoc](https://godoc.org/code.dumpstack.io/lib/cryptocurrency?status.svg)](https://godoc.org/code.dumpstack.io/lib/cryptocurrency)
[![Go Report Card](https://goreportcard.com/badge/code.dumpstack.io/lib/cryptocurrency)](https://goreportcard.com/report/code.dumpstack.io/lib/cryptocurrency)
[![Donate](https://img.shields.io/badge/donate-paypal-blue.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=R8W2UQPZ5X5JE&source=url)
[![Donate](https://img.shields.io/badge/donate-bitcoin-blue.svg)](https://blockchair.com/bitcoin/address/bc1q23fyuq7kmngrgqgp6yq9hk8a5q460f39m8nv87)

# Cryptocurrency API

Stateless cryptocurrency API.

Requirements:
- Bitcoin: electrum wallet in $PATH
- Ethereum: `export INFURA_API_KEY=...` from infura.io (it's free). It'll work even without, but there are some limits on queries.

Start:

    go get -u code.dumpstack.io/lib/cryptocurrency
    go test -v code.dumpstack.io/lib/cryptocurrency/...

Usage:
```go
package main

import (
	"log"

	"code.dumpstack.io/lib/cryptocurrency"
)

func main() {
	c := cryptocurrency.Bitcoin
	dest := "mk84dHbQoUHWaWGuYspx6GXWgcjB9CuQqw"
	// c := cryptocurrency.Ethereum
	// dest := "0xD98660C76443A8A043a19499048EeC4FB06f2581"
	// c := cryptocurrency.Cardano
	// dest := "Ae2tdPwUPEZ68cfEjZjKKRabiqbazMtP69uGaM2pMZRg87fvn4FGvR95BEV"

	err := c.Testnet(true)
	if err != nil {
		log.Fatal(err)
	}

	seed, address, err := c.GenWallet()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(seed, address)

	balance, err := c.Balance(seed)
	log.Println(balance)
	if err != nil {
		log.Fatal(err)
	}

	valid, err := c.Validate(dest)
	if err != nil {
		log.Fatal(err)
	}
	if valid {
		log.Println("address", dest, "is valid")
	} else {
		log.Fatal("address", dest, "is invalid")
	}

	amount := float64(0.1)
	tx, err := c.Send(seed, dest, amount)
	// tx, err := c.SendUnits(seed, dest, wei) // precise version
	if err != nil {
		log.Println("here it'll exit because there's no money inside new wallet")
		log.Fatal(err)
	}
	log.Println(tx)
}
```
