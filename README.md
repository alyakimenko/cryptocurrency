# Cryptocurrency API

Stateless cryptocurrency API.

Usage:

    package main

    import (
        "log"

        "code.dumpstack.io/lib/cryptocurrency"
        // "code.dumpstack.io/lib/cryptocurrency/bitcoin"
    )

    func main() {
        // bitcoin.Testnet = true

        seed, address, err := cryptocurrency.Bitcoin.GenWallet()
        if err != nil {
            log.Fatal(err)
        }
        log.Println(seed, address)

        balance, err := cryptocurrency.Bitcoin.Balance(seed)
        log.Println(balance)
        if err != nil {
            log.Fatal(err)
        }

        dest := "bc1q23fyuq7kmngrgqgp6yq9hk8a5q460f39m8nv87"
        amount := float64(0.1)
        tx, err := cryptocurrency.Bitcoin.Send(seed, dest, amount)
        if err != nil {
            // here it'll exit because there's no money inside new wallet
            log.Fatal(err)
        }
        log.Println(tx)
    }
