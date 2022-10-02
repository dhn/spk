// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/projectdiscovery/gologger"
)

// JSON object whois
type jsonWhois struct {
	CIDR   string `json:"cidr"`
	Source string `json:"source"`
}

// Print results as JSON
func WriteJSON(results <-chan Result) {
	encoder := jsoniter.NewEncoder(os.Stdout)
	var whois jsonWhois

	for result := range results {
		whois.CIDR = result.Value
		whois.Source = result.Source
		err := encoder.Encode(&whois)
		if err != nil {
			gologger.Fatal().Msgf(err.Error())
		}
	}
}
