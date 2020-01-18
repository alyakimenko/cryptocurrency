// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"github.com/tyler-smith/go-bip39"
)

func genSeed() (seed string, err error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return
	}

	return bip39.NewMnemonic(entropy)
}
