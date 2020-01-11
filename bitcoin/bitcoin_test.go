package bitcoin

import (
	"testing"
	"time"
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
	var addr string
	if Testnet {
		addr = "mpob68igSVfmaRyvucXHjJEdpbWG5Gt8dR"
	} else {
		addr = "1An9UvkeF1b57u448To7wqZ34HLEkSqCQ1"
	}
	valid, err := Validate(addr)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("should be valid")
	}

	valid, err = Validate("WRONGAn9UvkeF1b57u448To7wqZ34HLEkSqCQ1")
	if err != nil {
		t.Fatal(err)
	}

	if valid {
		t.Fatal("should be invalid")
	}
}

func TestBalance(t *testing.T) {
	// wallet with zero balance
	seed := "differ come sugar drift clump athlete " +
		"sweet fiscal uncle dilemma cage garbage"
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

	// wallet with some test btc inside
	seed = "flag release number shift amazing bacon " +
		"trend maximum lawsuit start traffic feel"
	balance, err = Balance(seed)
	if err != nil {
		t.Fatal(err)
	}
	if balance == 0 {
		t.Fatal("Balance of test wallet should not be zero", seed, balance)
	}
}

func TestSend(t *testing.T) {
	// wallet with some test btc inside
	// do not forget to put some btc to
	//     n4DSCXMeKjRRjBQHHvKURsLtmrcyfijTnN
	// from time to time
	seed := "act sentence begin build tornado note " +
		"then jungle jar govern bird dinner"
	balance, err := Balance(seed)
	if err != nil {
		t.Fatal(err)
	}

	if balance == 0 {
		t.Fatal("Please, open https://coinfaucet.eu/en/btc-testnet/ " +
			"and put address `n4DSCXMeKjRRjBQHHvKURsLtmrcyfijTnN`")
	}

	newseed, address, err := GenWallet()
	if err != nil {
		t.Fatal(err)
	}

	_, err = Send(seed, address, 0.0000001)
	if err != nil {
		t.Fatal(err)
	}

	received := false
	for start := time.Now(); time.Since(start) < time.Minute; {
		_, unconfirmed, err := RawBalance(newseed)
		if err != nil {
			t.Fatal(err)
		}

		if unconfirmed != 0 {
			received = true
			break
		}
	}

	if !received {
		t.Fatal("Does not received btc")
	}
}
