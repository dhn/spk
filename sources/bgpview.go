// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/dhn/spk/utils"
	jsoniter "github.com/json-iterator/go"
)

// BGPView JSON results
type bgpview struct {
	Data struct {
		Ipv4Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"ipv4_prefixes"`
		Ipv6Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"ipv6_prefixes"`
	}
}

// GetBGPData function returns all netranges based on the given organization name
func GetBGPData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getBGPViewData(fmt.Sprintf("https://api.bgpview.io/search?query_term=%s",
			url.QueryEscape(organization)), results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the JSON response
func getBGPViewData(sourceURL string, results chan utils.Result) {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{"Content-Type": "application/json"})

	var response bgpview
	err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
	if err != nil {
		return
	}

	ipv4 := response.Data.Ipv4Prefixes
	ipv6 := response.Data.Ipv6Prefixes

	for _, data := range append(ipv4, ipv6...) {
		results <- utils.Result{Value: data.Prefix, Source: "bgpview"}
	}
}
