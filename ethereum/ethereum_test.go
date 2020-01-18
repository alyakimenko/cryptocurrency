// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package ethereum

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func init() {
	Testnet = true
}

func TestGenWallet(t *testing.T) {
	seed, address, err := GenWallet()
	if err != nil {
		t.Fatal(err)
	}

	balance, err := Balance(seed)
	if err != nil {
		return
	}

	if balance != 0 {
		t.Fatal("BINGO (balance of new wallet != 0)", seed, balance)
	}

	valid, err := Validate(address)
	if err != nil || !valid {
		t.Fatal("Generated address is invalid", address, err)
	}
}

func TestValidate(t *testing.T) {
	valid, err := Validate("0x2aC8A61254884D847F460d318c598ab05D0ab383")
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("should be valid")
	}

	valid, err = Validate("0xINVALID254884D847F460d318c598ab05D0ab383")
	if err != nil {
		t.Fatal(err)
	}

	if valid {
		t.Fatal("should be invalid")
	}
}

func TestBalance(t *testing.T) {
	// wallet with zero balance
	seed := "dish fence ask vapor plate cart rival bomb " +
		"when snake domain asset february nurse vocal divorce " +
		"draft recycle chronic symptom exist average marine know"
	balance, err := Balance(seed)
	if err != nil {
		t.Fatal(err)
	}
	if balance != 0 {
		t.Fatal("BINGO (balance of test wallet != 0)", seed, balance)
	}

	_, err = Balance("some garbage instead of seed")
	if err == nil {
		t.Fatal("Balance does not returns error on invalid seed")
	}

	// wallet with some test eth inside
	seed = "harvest spider shine error require stereo syrup ritual " +
		"hurt inquiry faculty siege legend reveal happy bonus " +
		"grant tomorrow battle unfold security razor plastic must"
	balance, err = Balance(seed)
	if err != nil {
		t.Fatal(err)
	}
	if balance == 0 {
		t.Fatal("Balance of test wallet should not be zero", seed, balance)
	}
}

func TestSend(t *testing.T) {
	// wallet with some test eth inside
	// do not forget to put some eth to
	//     0x2aC8A61254884D847F460d318c598ab05D0ab383
	// from time to time
	seed := "harvest spider shine error require stereo syrup ritual " +
		"hurt inquiry faculty siege legend reveal happy bonus " +
		"grant tomorrow battle unfold security razor plastic must"
	balance, err := Balance(seed)
	if err != nil {
		t.Fatal(err)
	}

	if balance == 0 {
		t.Fatal("Please, open https://faucet.ropsten.be/ " +
			"and put address `0x2aC8A61254884D847F460d318c598ab05D0ab383`")
	}

	newseed, address, err := GenWallet()
	if err != nil {
		t.Fatal(err)
	}

	tx, err := Send(seed, address, 0.0000001)
	if err != nil {
		t.Fatal(err)
	}

	_, err = hexutil.Decode(tx)
	if err != nil {
		t.Fatal(err)
	}

	received := false
	for start := time.Now(); time.Since(start) < time.Minute; {
		unconfirmed, err := UnconfirmedBalanceWei(newseed)
		if err != nil {
			t.Fatal(err)
		}

		if unconfirmed.Int64() != 0 {
			received = true
			break
		}
	}

	if !received {
		t.Fatal("Does not received eth")
	}
}

func TestGetAddress(t *testing.T) {
	seed, address, err := GenWallet()
	if err != nil {
		t.Fatal(err)
	}

	newaddress, err := GetAddress(seed)
	if err != nil {
		t.Fatal(err)
	}

	if address != newaddress {
		t.Fatal("the address is not the same")
	}
}
