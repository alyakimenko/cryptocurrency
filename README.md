# Cryptocurrency API

Stateless cryptocurrency API.

Requirements:
- Bitcoin: electrum wallet in $PATH
- Ethereum: `export INFURA_API_KEY=...` from infura.io (it's free). It'll work even without, but there are some limits on queries.

Start:

    go get -u code.dumpstack.io/lib/cryptocurrency
    go test -v code.dumpstack.io/lib/cryptocurrency/...

Usage:

    package main

    import (
        "log"

        "code.dumpstack.io/lib/cryptocurrency"
    )

    func main() {
        c := cryptocurrency.Bitcoin
        dest := "mk84dHbQoUHWaWGuYspx6GXWgcjB9CuQqw"
        // c := cryptocurrency.Ethereum
        // dest := "0xD98660C76443A8A043a19499048EeC4FB06f2581"

        err := c.Testnet(true)
        if err != nil {
                log.Fatal(err)
        }

        seed, address, err := c.GenWallet()
        if err != nil {
                log.Fatal(err)
        }
        log.Println(seed, address)

        balance, err := c.Balance(seed)
        log.Println(balance)
        if err != nil {
                log.Fatal(err)
        }

        valid, err := c.Validate(dest)
        if err != nil {
                log.Fatal(err)
        }
        if valid {
                log.Println("address", dest, "is valid")
        } else {
                log.Fatal("address", dest, "is invalid")
        }

        amount := float64(0.1)
        tx, err := c.Send(seed, dest, amount)
        // tx, err := c.SendUnits(seed, dest, wei) // precise version
        if err != nil {
                log.Println("here it'll exit because there's no money inside new wallet")
                log.Fatal(err)
        }
        log.Println(tx)
    }
