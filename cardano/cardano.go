// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"regexp"

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

func SendAll(seed, destination string) (tx string, err error) {
	balance, err := BalanceLovelace(seed)
	if err != nil {
		return
	}

	// TODO calculate fees (right now just cutoff 0.5 ADA)
	amount := new(big.Int).Sub(balance, big.NewInt(500000))

	return SendLovelace(seed, destination, amount)
}

func Send(seed, destination string, amountf float64) (tx string, err error) {
	amountbigf := new(big.Float).SetFloat64(amountf)
	ada := big.NewFloat(1000000)
	amount, _ := new(big.Float).Mul(amountbigf, ada).Int(nil)
	return SendLovelace(seed, destination, amount)
}

func SendLovelace(seed, destination string, amount *big.Int) (tx string, err error) {
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

	var addresses []string
	err = json.Unmarshal([]byte(sAddresses), &addresses)
	if err != nil {
		return
	}

	utxos, err := utxoForAddresses(addresses)
	if err != nil {
		return
	}
	if len(utxos) == 0 {
		err = errors.New("no utxos")
		return
	}

	type input struct {
		TxIN struct {
			Index int    `json:"index"`
			ID    string `json:"id"`
		} `json:"ptr"`
		Value struct {
			Address string `json:"address"`
			Value   string `json:"value"`
		} `json:"value"`
		Addressing struct {
			Account int `json:"account"`
			Change  int `json:"change"`
			Index   int `json:"index"`
		} `json:"addressing"`
	}

	var inputs []input
	for _, utxo := range utxos {
		var in input

		in.TxIN.Index = utxo.TxIndex
		in.TxIN.ID = utxo.TxHash

		in.Value.Address = utxo.Receiver
		in.Value.Value = utxo.Amount

		in.Addressing.Index = utxo.TxIndex

		inputs = append(inputs, in)
	}

	bytes, err := json.Marshal(inputs)
	if err != nil {
		return
	}

	type output struct {
		Address string `json:"address"`
		Value   string `json:"value"`
	}
	var outputs []output
	outputs = append(outputs, output{
		Address: destination,
		Value:   amount.String(),
	})

	outputBytes, err := json.Marshal(outputs)
	if err != nil {
		return
	}

	// TODO here may be the issue with change address
	// we should choose empty address each time, probably
	raw, err := xWalletSpend(wallet, string(bytes), string(outputBytes), addresses[1])
	if err != nil {
		return
	}

	var result struct {
		Tx []byte `json:"cbor_encoded_tx"`
	}
	err = json.Unmarshal([]byte(raw), &result)
	if err != nil {
		return
	}

	signedTx := base64.StdEncoding.EncodeToString(result.Tx)
	err = sendSignedTx(signedTx)
	if err != nil {
		return
	}

	// TODO there should be a way to do it normally without
	// parsing transactions from the server
	return getTx(addresses, utxos)
}

func Validate(address string) (valid bool, err error) {
	// TODO
	re := regexp.MustCompile("^[0-9a-zA-Z]{32,256}$")
	valid = re.MatchString(address)
	return
}
