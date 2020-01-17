package ethereum

import (
	"context"
	"errors"
	"log"
	"math/big"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var Testnet = false
var InfuraAPIKey = os.Getenv("INFURA_API_KEY")

func endpoint() (url string) {
	if Testnet {
		url = "https://ropsten.infura.io"
	} else {
		url = "https://mainnet.infura.io"
	}

	if InfuraAPIKey != "" {
		url += "/v3/" + InfuraAPIKey
	}
	return
}

func accountFromMnemonic(seed string) (account accounts.Account, err error) {
	wallet, err := hdwallet.NewFromMnemonic(seed)
	if err != nil {
		return
	}

	return wallet.Derive(hdwallet.DefaultBaseDerivationPath, false)
}

func GenWallet() (seed, address string, err error) {
	seed, err = hdwallet.NewMnemonic(256)
	if err != nil {
		return
	}

	account, err := accountFromMnemonic(seed)
	if err != nil {
		return
	}

	address = account.Address.Hex()
	return
}

func BalanceWei(seed string) (amount *big.Int, err error) {
	account, err := accountFromMnemonic(seed)
	if err != nil {
		return
	}

	client, err := ethclient.Dial(endpoint())
	if err != nil {
		return
	}

	ctx := context.Background()
	amount, err = client.BalanceAt(ctx, account.Address, nil)
	return
}

func Balance(seed string) (amount float64, err error) {
	wei, err := BalanceWei(seed)
	if err != nil {
		return
	}

	weiFloat := new(big.Float).SetInt(wei)
	oneEth := big.NewFloat(1000000000000000000)
	amount, _ = new(big.Float).Quo(weiFloat, oneEth).Float64()
	return
}

func Validate(address string) (valid bool, err error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	valid = re.MatchString(address)
	return
}

func UnconfirmedBalanceWei(seed string) (unconfirmed *big.Int, err error) {
	account, err := accountFromMnemonic(seed)
	if err != nil {
		return
	}

	client, err := ethclient.Dial(endpoint())
	if err != nil {
		return
	}

	return client.PendingBalanceAt(context.Background(), account.Address)
}

func SendWei(seed, destination string, amount *big.Int) (tx string, err error) {
	wallet, err := hdwallet.NewFromMnemonic(seed)
	if err != nil {
		return
	}

	account, err := wallet.Derive(hdwallet.DefaultBaseDerivationPath, true)
	if err != nil {
		return
	}

	client, err := ethclient.Dial(endpoint())
	if err != nil {
		return
	}

	ctx := context.Background()

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return
	}

	nonce, err := client.PendingNonceAt(ctx, account.Address)
	if err != nil {
		return
	}

	toAddress := common.HexToAddress(destination)

	// The gas limit for a standard ETH transfer is 21000 units.
	gasLimit := uint64(21000) // in units

	ethTx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	signedTx, err := wallet.SignTx(account, ethTx, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return
	}

	tx = signedTx.Hash().Hex()
	return
}

func Send(seed, destination string, amount float64) (tx string, err error) {
	oneEth := big.NewFloat(1000000000000000000)
	amountWeiFloat := big.NewFloat(amount)
	amountWei, _ := new(big.Float).Mul(oneEth, amountWeiFloat).Int(nil)
	return SendWei(seed, destination, amountWei)
}

func SendAll(seed, destination string) (tx string, err error) {
	wei, err := BalanceWei(seed)
	if err != nil {
		return
	}

	if big.NewInt(0).Cmp(wei) == 0 {
		err = errors.New("zero balance")
		return
	}

	client, err := ethclient.Dial(endpoint())
	if err != nil {
		return
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}

	// The gas limit for a standard ETH transfer is 21000 units.
	gasLimit := new(big.Int).SetUint64(uint64(21000)) // in units

	fee := new(big.Int).Mul(gasPrice, gasLimit)

	return SendWei(seed, destination, new(big.Int).Sub(wei, fee))
}
