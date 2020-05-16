// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"strings"
	"testing"
)

func TestGenerateSeed(t *testing.T) {
	seed, err := genSeed()
	if err != nil {
		t.Fatal(err)
	}

	if len(seed) == 0 {
		t.Fatal("seed is empty")
	}

	if len(strings.Split(seed, " ")) != 15 {
		t.Fatal("seed is not 15 words")
	}
}

func TestValidate(t *testing.T) {
	valid, err := Validate("DdzFFzCqrhsiZuCYMRy6R7SBkgQj3J4NxnSbVvrcg8sHneqdhe7cjUjH9AZ5C8mZxWXt6JCciNrnXrWtHzWywbU5RqVW511gJrE8Uk2d")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("!valid for valid")
	}

	valid, err = Validate("Ae2tdPwUPEZ68cfEjZjKKRabiqbazMtP69uGaM2pMZRg87fvn4FGvR95BEV")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal("!valid for valid")
	}
}
