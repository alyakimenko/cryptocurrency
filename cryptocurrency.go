// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cryptocurrency

import (
	"errors"
	"math/big"

	"code.dumpstack.io/lib/cryptocurrency/bitcoin"
	"code.dumpstack.io/lib/cryptocurrency/cardano"
	"code.dumpstack.io/lib/cryptocurrency/ethereum"
)

type Cryptocurrency int

const (
	Bitcoin Cryptocurrency = iota
	Ethereum
	Cardano
	// TODO:
	// Monero
)

// Cryptocurrencies list
var Cryptocurrencies = []Cryptocurrency{
	Bitcoin,
	Ethereum,
	Cardano,
}

func (t Cryptocurrency) MarshalText() (b []byte, err error) {
	b = []byte(t.Symbol())
	return
}

func (t *Cryptocurrency) UnmarshalText(data []byte) (err error) {
	*t, err = FromSymbol(string(data))
	return
}

// Symbol of cryptocurrency (btc, eth, etc.)
func (t Cryptocurrency) Symbol() string {
	switch t {
	case Bitcoin:
		return "btc"
	case Ethereum:
		return "eth"
	case Cardano:
		return "ada"
	}
	panic(nil)
}

func FromSymbol(symbol string) (cc Cryptocurrency, err error) {
	switch symbol {
	case "btc":
		cc = Bitcoin
	case "eth":
		cc = Ethereum
	case "ada":
		cc = Cardano
	default:
		err = errors.New("unknown cryptocurrency")
	}
	return
}

// Testnet enable or disable
func (t Cryptocurrency) Testnet(state bool) (err error) {
	switch t {
	case Bitcoin:
		bitcoin.Testnet = true
		return
	case Ethereum:
		ethereum.Testnet = true
		return
	}

	err = errors.New("Not supported yet")
	return
}

// GenWallet for specified cryptocurrency
func (t Cryptocurrency) GenWallet() (seed, address string, err error) {
	switch t {
	case Bitcoin:
		seed, address, err = bitcoin.GenWallet()
		return
	case Ethereum:
		seed, address, err = ethereum.GenWallet()
		return
	case Cardano:
		seed, address, err = cardano.GenWallet()
		return
	}

	err = errors.New("Not supported yet")
	return
}

// GetAddress for the wallet (can be the same or different every time)
func (t Cryptocurrency) GetAddress(seed string) (address string, err error) {
	switch t {
	case Bitcoin:
		address, err = bitcoin.GetAddress(seed)
		return
	case Ethereum:
		address, err = ethereum.GetAddress(seed)
		return
	case Cardano:
		address, err = cardano.GetAddress(seed)
		return
	}

	err = errors.New("Not supported yet")
	return
}

// Balance of the wallet (not address!)
func (t Cryptocurrency) Balance(seed string) (amount float64, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Balance(seed)
	case Ethereum:
		return ethereum.Balance(seed)
	case Cardano:
		return cardano.Balance(seed)
	}

	err = errors.New("Not supported yet")
	return
}

// BalanceUnits returns the balance of the wallet (not address!) in Satoshi/Wei/etc.
func (t Cryptocurrency) BalanceUnits(seed string) (units *big.Int, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.BalanceSatoshi(seed)
	case Ethereum:
		return ethereum.BalanceWei(seed)
	case Cardano:
		return cardano.BalanceLovelace(seed)
	}

	err = errors.New("Not supported yet")
	return
}

// SendUnits amount of Satoshi/Wei/etc. to the address dest
func (t Cryptocurrency) SendUnits(seed, dest string, units *big.Int) (tx string, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.SendSatoshi(seed, dest, units)
	case Ethereum:
		return ethereum.SendWei(seed, dest, units)
	}

	err = errors.New("Not supported yet")
	return
}

// Send the amount of cryptocurrency to destination address
func (t Cryptocurrency) Send(seed, dest string, amount float64) (tx string, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Send(seed, dest, amount)
	case Ethereum:
		return ethereum.Send(seed, dest, amount)
	}

	err = errors.New("Not supported yet")
	return
}

// SendAll cryptocurrency to destination address
func (t Cryptocurrency) SendAll(seed, dest string) (tx string, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.SendAll(seed, dest)
	case Ethereum:
		return ethereum.SendAll(seed, dest)
	}

	err = errors.New("Not supported yet")
	return
}

// Validate cryptocurrency address
func (t Cryptocurrency) Validate(address string) (valid bool, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Validate(address)
	case Ethereum:
		return ethereum.Validate(address)
	}

	err = errors.New("Not supported yet")
	return
}
