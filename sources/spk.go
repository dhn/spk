// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/dhn/spk/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/projectdiscovery/gologger"
)

// SPK JSON results
type spk []struct {
	Inetnum string `json:"inetnum"`
	Netname string `json:"netname"`
	Descr   string `json:"descr"`
	Country string `json:"country"`
	Source  string `json:"source"`
}

// GetSPKData function returns all netranges based on the given organization name
func GetSPKData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getSPKData(fmt.Sprintf("https://spk.osint.quest/descr?v=%s",
			url.QueryEscape(organization)), results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the JSON response
func getSPKData(sourceURL string, results chan utils.Result) {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{})

	var response spk
	err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
	if err != nil {
		gologger.Fatal().Msgf(err.Error())
	}

	for _, result := range response {
		if strings.Contains(result.Inetnum, ":") {
			results <- utils.Result{Value: result.Inetnum, Source: "spk"}
		} else {
			for _, cidr := range utils.RipeToCIDR(result.Inetnum) {
				results <- utils.Result{Value: cidr.String(), Source: "spk"}
			}
		}
	}
}
