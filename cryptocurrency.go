package cryptocurrency

import (
	"errors"
	"math/big"

	"code.dumpstack.io/lib/cryptocurrency/bitcoin"
	"code.dumpstack.io/lib/cryptocurrency/ethereum"
)

type Cryptocurrency int

const (
	Bitcoin Cryptocurrency = iota
	Ethereum
	Monero
	Cardano
)

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

func (t Cryptocurrency) GenWallet() (seed, address string, err error) {
	switch t {
	case Bitcoin:
		seed, address, err = bitcoin.GenWallet()
		return
	case Ethereum:
		seed, address, err = ethereum.GenWallet()
		return
	}

	err = errors.New("Not supported yet")
	return
}

func (t Cryptocurrency) Balance(seed string) (amount float64, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Balance(seed)
	case Ethereum:
		return ethereum.Balance(seed)
	}

	err = errors.New("Not supported yet")
	return
}

// BalanceUnits returns the balance of the wallet (not address!) in Satoshi/Wei/etc.
func (t Cryptocurrency) BalanceUnits(seed string) (units *big.Int, err error) {
	switch t {
	case Ethereum:
		return ethereum.BalanceWei(seed)
	}

	err = errors.New("Not supported yet")
	return
}

// SendUnits send units amount of Satoshi/Wei/etc. to the address dest
func (t Cryptocurrency) SendUnits(seed, dest string, units *big.Int) (tx string, err error) {
	switch t {
	case Ethereum:
		return ethereum.SendWei(seed, dest, units)
	}

	err = errors.New("Not supported yet")
	return
}

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
