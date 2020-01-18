package bitcoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Testnet = false

func init() {
	if exec.Command("which", "electrum").Run() != nil {
		log.Println("`electrum` is not found in $PATH")
	}
}

func electrum(wallet string, args ...string) (output []byte, err error) {
	params := []string{}
	if Testnet {
		params = append(params, "--testnet")
	}
	params = append(params, "--wallet", wallet)
	params = append(params, args...)

	return exec.Command("electrum", params...).Output()
}

func startDaemon(wallet string) (err error) {
	output, err := electrum(wallet, "daemon", "start")
	if err != nil {
		return
	}

	output, err = electrum(wallet, "daemon", "load_wallet")
	if err != nil {
		return
	}

	err = errors.New("Can't synchronize electrum")
	for start := time.Now(); time.Since(start) < time.Minute; {
		output, _ = electrum(wallet, "is_synchronized")
		if string(output) == "true\n" {
			err = nil
			break
		}
		time.Sleep(time.Second)
	}
	return
}

func stopDaemon(wallet string) (err error) {
	_, err = electrum(wallet, "daemon", "stop")
	return
}

func GenWallet() (seed, address string, err error) {
	dir, err := ioutil.TempDir("/tmp/", "cryptocurrency_")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	wallet := filepath.Join(dir, "wallet")

	output, err := electrum(wallet, "create")
	if err != nil {
		return
	}

	var result struct{ Seed string }
	err = json.Unmarshal(output, &result)
	if err != nil {
		return
	}
	seed = strings.Trim(result.Seed, " \r\n")

	err = startDaemon(wallet)
	if err != nil {
		return
	}
	defer stopDaemon(wallet)

	output, err = electrum(wallet, "getunusedaddress")
	if err != nil {
		return
	}

	address = strings.Trim(string(output), " \r\n")
	return
}

func parseBalance(output []byte) (confirmed, unconfirmed float64, err error) {
	var result struct{ Confirmed, Unconfirmed string }
	err = json.Unmarshal(output, &result)
	if err != nil {
		return
	}

	confirmed, err = strconv.ParseFloat(result.Confirmed, 64)
	if err != nil {
		return
	}

	if result.Unconfirmed != "" {
		unconfirmed, err = strconv.ParseFloat(result.Unconfirmed, 64)
		if err != nil {
			return
		}
	}

	return
}

func Balance(seed string) (amount float64, err error) {
	amount, unconfirmed, err := RawBalance(seed)
	if unconfirmed < 0 {
		amount += unconfirmed // subtraction
	}
	return
}

func BalanceSatoshi(seed string) (amount *big.Int, err error) {
	amountf, err := Balance(seed)
	if err != nil {
		return
	}

	amountbigf := new(big.Float).SetFloat64(amountf)
	btc := big.NewFloat(100000000)
	amount, _ = new(big.Float).Mul(amountbigf, btc).Int(nil)
	return
}

func RawBalance(seed string) (confirmed, unconfirmed float64, err error) {
	dir, err := ioutil.TempDir("/tmp/", "cryptocurrency_")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	wallet := filepath.Join(dir, "wallet")

	output, err := electrum(wallet, "restore", seed)
	if err != nil {
		return
	}

	err = startDaemon(wallet)
	if err != nil {
		return
	}
	defer stopDaemon(wallet)

	output, err = electrum(wallet, "getbalance")
	if err != nil {
		return
	}

	confirmed, unconfirmed, err = parseBalance(output)
	return
}

func Validate(btc string) (valid bool, err error) {
	output, err := electrum("", "validateaddress", btc)
	if err != nil {
		return
	}

	switch string(output) {
	case "true\n":
		valid = true
		break
	case "false\n":
		valid = false
		break
	default:
		err = errors.New("electrum output is invalid")
	}
	return
}

func send(seed, destination string, amount string) (tx string, err error) {
	dir, err := ioutil.TempDir("/tmp/", "cryptocurrency_")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)

	wallet := filepath.Join(dir, "wallet")

	_, err = electrum(wallet, "restore", seed)
	if err != nil {
		return
	}

	err = startDaemon(wallet)
	if err != nil {
		return
	}
	defer stopDaemon(wallet)

	output, err := electrum(wallet, "payto", destination, amount)
	if err != nil {
		return
	}

	var result struct {
		Complete bool
		Final    bool
		Hex      string
	}
	err = json.Unmarshal(output, &result)
	if err != nil {
		return
	}

	if !result.Complete {
		err = errors.New("Transaction is not complete")
		return
	}

	output, err = electrum(wallet, "broadcast", result.Hex)
	if err != nil {
		return
	}
	tx = strings.Trim(string(output), " \r\n")
	return
}

func SendSatoshi(seed, destination string, units *big.Int) (tx string, err error) {
	btc := big.NewFloat(100000000)
	unitsf := new(big.Float).SetInt(units)
	amount := new(big.Float).Quo(unitsf, btc).Text('f', 8)
	return send(seed, destination, amount)
}

func Send(seed, destination string, amount float64) (tx string, err error) {
	return send(seed, destination, fmt.Sprintf("%.8f", amount))
}

func SendAll(seed, destination string) (tx string, err error) {
	return send(seed, destination, "!")
}
