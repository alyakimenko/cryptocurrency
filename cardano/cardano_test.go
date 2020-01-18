// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import "testing"

func TestGenerateSeed(t *testing.T) {
	seed, err := genSeed()
	if err != nil {
		t.Fatal(err)
	}

	if len(seed) == 0 {
		t.Fatal("seed is empty")
	}
}
