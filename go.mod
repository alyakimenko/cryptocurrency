module code.dumpstack.io/lib/cryptocurrency

replace code.dumpstack.io/lib/cryptocurrency/bitcoin => ./bitcoin

replace code.dumpstack.io/lib/cryptocurrency/cardano => ./cardano

replace code.dumpstack.io/lib/cryptocurrency/ethereum => ./ethereum

go 1.12

require (
	github.com/ethereum/go-ethereum v1.9.9
	github.com/jollheef/go-js-cardano-wasm-bin v0.0.0-20200118195930-07db337aa4d2
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20191219011559-4d9d106de2bd
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/wasmerio/go-ext-wasm v0.0.0-20200116105756-c27bd6f62e71
)
