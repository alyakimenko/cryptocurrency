// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cryptocurrency

import "testing"

func TestGenWalletAPI(t *testing.T) {
	for _, cc := range Cryptocurrencies {
		seed, address, err := cc.GenWallet()
		if err != nil {
			t.Fatal(err)
		}

		if seed == "" {
			t.Fatal("no seed")
		}

		if address == "" {
			t.Fatal("no address")
		}
	}
}
