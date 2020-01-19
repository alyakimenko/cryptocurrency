// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"fmt"
)

func xWalletCreateDaedalusMnemonic(seed string) (wallet string, err error) {
	return jsonapi("xwallet_create_daedalus_mnemonic", "\""+seed+"\"")
}

func xWalletAccount(wallet string, n int) (account string, err error) {
	payload := fmt.Sprintf("{ \"wallet\": %s, \"account\": %d}", wallet, n)
	return jsonapi("xwallet_account", payload)
}

func xWalletAddresses(account string) (addresses string, err error) {
	// TODO add parameter how many addresses to get
	payload := "{ \"account\": " + account +
		", \"address_type\": \"External\", \"indices\": [0,1,2,3,4]}"
	return jsonapi("xwallet_addresses", payload)
}

func xWalletSpend(wallet, inputs, output, change string) (signedTx string, err error) {
	payload := "{ \"wallet\": " + wallet +
		", \"inputs\": " + inputs + ", " +
		"\"outputs\": " + output + ", " +
		"\"change_addr\": \"" + change + "\"}"
	return jsonapi("xwallet_spend", payload)
}
