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

// Ripe JSON results
type ripe struct {
	Result struct {
		NumFound int `json:"numFound"`
		Docs     []struct {
			Doc struct {
				Strs []struct {
					Str struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"str"`
				} `json:"strs"`
			} `json:"doc"`
		} `json:"docs"`
	} `json:"result"`
}

// GetRipeData function returns all netranges found in the RIPE database
func GetRipeData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getRipeDBData(fmt.Sprintf("https://apps.db.ripe.net/db-web-ui/api/rest/fulltextsearch/select?facet=true&hl=true&q=(descr:(%s))%%%%20AND%%%%20(object-type:inetnum)&start=%%d&wt=json",
			url.QueryEscape(organization)), results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the JSON response
func getRipeDBData(sourceURL string, results chan utils.Result) {
	entries := getEntries(sourceURL)

	for index := 0; index <= entries; index += 10 {
		sourceURL := fmt.Sprintf(sourceURL, index)
		resp := utils.GetHTTPRequest(sourceURL, map[string]string{"Accept": "application/json"})

		var response ripe
		err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
		if err != nil {
			gologger.Fatal().Msgf(err.Error())
		}

		for _, result := range response.Result.Docs {
			for _, doc := range result.Doc.Strs {
				if strings.Compare(doc.Str.Name, "inetnum") == 0 {
					for _, cidr := range utils.RipeToCIDR(doc.Str.Value) {
						results <- utils.Result{Value: cidr.String(), Source: "ripe"}
					}
				}
			}
		}
	}
}

// Receive the list of entries from the JSON object
func getEntries(sourceURL string) int {
	sourceURL = fmt.Sprintf(sourceURL, 0)
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{"Accept": "application/json"})

	var response ripe
	err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
	if err != nil {
		gologger.Fatal().Msgf(err.Error())
	}

	return response.Result.NumFound
}
