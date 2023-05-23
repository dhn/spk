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
)

// APNIC JSON results
type apnic []struct {
	Attributes []struct {
		Name   string   `json:"name"`
		Values []string `json:"values,omitempty"`
	} `json:"attributes,omitempty"`
}

// GetAPNICData function returns all netranges found in the APNIC database
func GetAPNICData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getAPNICDData(fmt.Sprintf("http://wq.apnic.net/query?searchtext=%s",
			url.QueryEscape(organization)), organization, results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the JSON response
func getAPNICDData(sourceURL string, organization string, results chan utils.Result) {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{})

	var response apnic
	err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
	if err != nil {
		return
	}

	// Example Response:
	// {"type":"object","attributes":[
	// 	{"name":"inetnum","values":["61.88.26.140 - 61.88.26.143"]},
	//  {"name":"descr","values":["OPTUS Customer Network"]},...
	for _, result := range response {
		for _, attributes := range result.Attributes {
			if attributes.Name == "descr" {
				for _, values := range attributes.Values {
					if strings.Contains(values, organization) {
						if result.Attributes[0].Name == "inetnum" {
							for _, cidr := range utils.RipeToCIDR(result.Attributes[0].Values[0]) {
								results <- utils.Result{Value: cidr.String(), Source: "apnic"}
							}
						}
					}
				}
			}
		}
	}
}
