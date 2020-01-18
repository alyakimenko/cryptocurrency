// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
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
