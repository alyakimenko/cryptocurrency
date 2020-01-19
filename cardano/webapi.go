// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func endpoint() (url string) {
	// TODO support testnet
	return "https://iohk-mainnet.yoroiwallet.com/api/"
}

func status() (ok bool, err error) {
	resp, err := http.Get(endpoint() + "status")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct{ IsServerOK bool }
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	ok = result.IsServerOK
	return
}

func utxoSumForAddresses(addresses []string) (sum int64, err error) {
	payload := struct {
		Addresses []string `json:"addresses"`
	}{Addresses: addresses}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", endpoint()+"txs/utxoSumForAddresses", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct{ Sum string }
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	if result.Sum == "" {
		// not error, just zero
		return
	}

	// >>> ada = 45_000_000_000
	// >>> lovelace = 1_000_000
	// >>> math.log2(ada * lovelace)
	// 55.32077451964011
	sum, err = strconv.ParseInt(result.Sum, 10, 64)
	return
}

type utxo struct {
	Amount   string `json:"amount"`
	BlockNum int    `json:"block_num"`
	Receiver string `json:"receiver"`
	TxHash   string `json:"tx_hash"`
	TxIndex  int    `json:"tx_index"`
	UtxoID   string `json:"utxo_id"`
}

func utxoForAddresses(addresses []string) (utxos []utxo, err error) {
	payload := struct {
		Addresses []string `json:"addresses"`
	}{Addresses: addresses}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", endpoint()+"txs/utxoForAddresses", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&utxos)
	return
}

func sendSignedTx(tx string) (err error) {
	data := struct {
		SignedTx string `json:"signedTx"`
	}{SignedTx: tx}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", endpoint()+"txs/signed", body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if string(bytes) == "[]" {
		// success
		return
	}
	//
	// failure; read error message
	//
	var result struct{ Code, Message string }
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}

	err = errors.New(fmt.Sprintf("[%s] %s", result.Code, result.Message))
	return
}

func inArray(u utxo, array []utxo) bool {
	for _, utxoFromArray := range array {
		if u.TxHash == utxoFromArray.TxHash {
			return true
		}
	}
	return false
}

func getTxTry(addrs []string, oldutxos []utxo) (found bool, tx string, err error) {
	// Looking for a new transaction
	newutxos, err := utxoForAddresses(addrs)
	if err != nil {
		return
	}

	// Looking for a new transaction
	for _, newutxo := range newutxos {
		if !inArray(newutxo, oldutxos) {
			found = true
			tx = newutxo.TxHash
		}
	}
	return
}

func getTx(addrs []string, oldutxos []utxo) (tx string, err error) {
	for start := time.Now(); time.Since(start) < time.Minute; {
		var found bool
		found, tx, err = getTxTry(addrs, oldutxos)
		if err != nil || found {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
