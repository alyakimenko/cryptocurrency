// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"encoding/json"
	"math/big"

	"github.com/tyler-smith/go-bip39"
)

func genSeed() (seed string, err error) {
	// TODO generate Yoroi-compatible mnemonic (15 words)
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return
	}

	return bip39.NewMnemonic(entropy)
}

func getAddresses(seed string) (addresses []string, err error) {
	wallet, err := xWalletCreateDaedalusMnemonic(seed)
	if err != nil {
		return
	}

	account, err := xWalletAccount(wallet, 0)
	if err != nil {
		return
	}

	sAddresses, err := xWalletAddresses(account)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(sAddresses), &addresses)
	return
}

// GetAddress for the wallet
func GetAddress(seed string) (address string, err error) {
	addresses, err := getAddresses(seed)
	if err != nil {
		return
	}

	address = addresses[0]
	return
}

// GenWallet for Cardano
func GenWallet() (seed, address string, err error) {
	seed, err = genSeed()
	if err != nil {
		return
	}

	address, err = GetAddress(seed)
	return
}

// BalanceLovelace returns the balance of the wallet (not address!) in lovelace
func BalanceLovelace(seed string) (amount *big.Int, err error) {
	addresses, err := getAddresses(seed)
	if err != nil {
		return
	}

	sum, err := utxoSumForAddresses(addresses)
	if err != nil {
		return
	}

	amount = new(big.Int).SetInt64(sum)
	return
}

// Balance of the wallet (not address!)
func Balance(seed string) (amount float64, err error) {
	lovelace, err := BalanceLovelace(seed)
	if err != nil {
		return
	}

	lovelaceFloat := new(big.Float).SetInt(lovelace)
	oneADA := big.NewFloat(1000000)
	amount, _ = new(big.Float).Quo(lovelaceFloat, oneADA).Float64()
	return
}
