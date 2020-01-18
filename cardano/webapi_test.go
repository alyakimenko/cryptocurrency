// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"testing"
)

func TestStatus(t *testing.T) {
	ok, err := status()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("server is dead")
	}
}

func Test_utxoSumForAddresses(t *testing.T) {
	addresses := []string{
		"Ae2tdPwUPEYz56j2qffuLfQEN9y4YuMbefJRqyfsHuuWX2TzRENMTX23pGr",
		"Ae2tdPwUPEYwJGMJfPkC3bRn6c5oSGTwwuEWuR3ofDG5T6c6vH7XL5DMmGH",
		"Ae2tdPwUPEZCdLWob8bJRXumWHNZ7rXWPjHGTq5qT9YeuaeXTicoNxFhmZt",
		"Ae2tdPwUPEZ2SU9qv63oVGTLpPL1iPw6Xmm1xRT98KWvQpgYZBbj8gWfmJt",
	}

	sum, err := utxoSumForAddresses(addresses)
	if err != nil {
		t.Fatal(err)
	}

	if sum != 32133700 {
		t.Fatal("Wrong UTXO Sum")
	}

	addresses = []string{
		"Ae2tdPwUPEZ6GQPHgXuM31aEFKm4ZHGULQfT2iSPDcqkMVbyDLfMhaoAaSd",
	}

	sum, err = utxoSumForAddresses(addresses)
	if err != nil {
		t.Fatal(err)
	}

	if sum != 0 {
		t.Fatal("Wrong UTXO Sum")
	}
}
