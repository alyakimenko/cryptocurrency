// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"encoding/json"
	"testing"
)

func isJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func Test_xWalletCreateDaedalusMnemonic(t *testing.T) {
	seed, err := genSeed()
	if err != nil {
		t.Fatal(err)
	}

	wallet, err := xWalletCreateDaedalusMnemonic(seed)
	if err != nil {
		t.Fatal(err)
	}

	if !isJSON(wallet) {
		t.Fatal("output is not json")
	}
}

func Test_xWalletAccount(t *testing.T) {
	seed, _ := genSeed()
	wallet, _ := xWalletCreateDaedalusMnemonic(seed)

	acc1, err := xWalletAccount(wallet, 0)
	if err != nil {
		t.Fatal(err)
	}

	if !isJSON(acc1) {
		t.Fatal("output is not json")
	}

	acc2, err := xWalletAccount(wallet, 1)
	if err != nil {
		t.Fatal(err)
	}

	if acc1 == acc2 {
		t.Fatal("different accounts are the same")
	}
}

func Test_xWalletAddresses(t *testing.T) {
	seed, _ := genSeed()
	wallet, _ := xWalletCreateDaedalusMnemonic(seed)
	account, _ := xWalletAccount(wallet, 0)

	addresses, err := xWalletAddresses(account)
	if err != nil {
		t.Fatal(err)
	}

	if !isJSON(addresses) {
		t.Fatal("output is not json")
	}
}
