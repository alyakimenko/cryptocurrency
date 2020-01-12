package cryptocurrency

import (
	"errors"

	"code.dumpstack.io/lib/cryptocurrency/bitcoin"
)

type Cryptocurrency int

const (
	Bitcoin Cryptocurrency = iota
	Ethereum
	Monero
	Cardano
)

func (t Cryptocurrency) GenWallet() (seed, address string, err error) {
	switch t {
	case Bitcoin:
		seed, address, err = bitcoin.GenWallet()
		return
	}

	err = errors.New("Not supported yet")
	return
}

func (t Cryptocurrency) Balance(seed string) (amount float64, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Balance(seed)
	}

	err = errors.New("Not supported yet")
	return
}

func (t Cryptocurrency) Send(seed, dest string, amount float64) (tx string, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.Send(seed, dest, amount)
	}

	err = errors.New("Not supported yet")
	return
}

func (t Cryptocurrency) SendAll(seed, dest string) (tx string, err error) {
	switch t {
	case Bitcoin:
		return bitcoin.SendAll(seed, dest)
	}

	err = errors.New("Not supported yet")
	return
}
