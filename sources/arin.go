// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sources

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dhn/spk/utils"
	jsoniter "github.com/json-iterator/go"
)

// ARIN JSON results
type arin struct {
	Net struct {
		NetBlocks struct {
			NetBlock struct {
				CidrLength struct {
					Value string `json:"$"`
				} `json:"cidrLength"`
				StartAddress struct {
					Value string `json:"$"`
				} `json:"startAddress"`
			} `json:"netBlock"`
		} `json:"netBlocks"`
	} `json:"net"`
}

// ARIN JSON results if the netblock element is an array
type arinArray struct {
	Net struct {
		NetBlocks struct {
			NetBlock []struct {
				CidrLength struct {
					Value string `json:"$"`
				} `json:"cidrLength"`
				StartAddress struct {
					Value string `json:"$"`
				} `json:"startAddress"`
			} `json:"netBlock"`
		} `json:"netBlocks"`
	} `json:"net"`
}

// GetARINData function returns all netranges based on the given organization name
func GetARINData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getARINData("https://whois.arin.net/ui/query.do", organization, results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the HTML response
func getARINData(sourceURL string, organization string, results chan utils.Result) {
	data := []byte(fmt.Sprintf("advanced=true&q=%s&r=NETWORK&NETWORK=handle&NETWORK=name", url.QueryEscape(organization)))
	resp := utils.PostHTTPRequest(sourceURL, data)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return
	}

	// Parse html to get the needed data
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		net := s.Find("a").Text()
		if strings.Contains(net, "NET-") {
			cidrs := getARINCIDR(fmt.Sprintf("http://whois.arin.net/rest/net/%s", url.QueryEscape(net)))
			for _, cidr := range cidrs {
				results <- utils.Result{Value: cidr, Source: "arin"}
			}
		}
	})
}

func getARINCIDR(sourceURL string) []string {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{"Accept": "application/json"})

	var response arin
	var startAddr, cidrLen string
	var cidr []string

	err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
	if err != nil {
		// Let's try to decode the JSON object again because sometimes the JSON
		// netblock element could be an array rather than a list of one element
		var response arinArray

		err := jsoniter.NewDecoder(bytes.NewReader(resp.Body())).Decode(&response)
		if err != nil {
			return cidr
		}

		for _, netblocks := range response.Net.NetBlocks.NetBlock {
			startAddr = netblocks.StartAddress.Value
			cidrLen = netblocks.CidrLength.Value
			cidr = append(cidr, fmt.Sprintf("%s/%s", startAddr, cidrLen))
		}

		return cidr
	}

	startAddr = response.Net.NetBlocks.NetBlock.StartAddress.Value
	cidrLen = response.Net.NetBlocks.NetBlock.CidrLength.Value
	cidr = append(cidr, fmt.Sprintf("%s/%s", startAddr, cidrLen))

	return cidr
}
