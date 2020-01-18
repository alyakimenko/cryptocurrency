// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import "testing"

func TestStatus(t *testing.T) {
	ok, err := status()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("server is dead")
	}
}
