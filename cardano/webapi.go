// Copyright 2020 Mikhail Klementev. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package cardano

import (
	"encoding/json"
	"net/http"
)

func endpoint() (url string) {
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
