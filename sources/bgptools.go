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
	"github.com/projectdiscovery/gologger"
)

// GetBGPToolsData function returns all netranges based on the given organization name
func GetBGPToolsData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getBGPToolsData(fmt.Sprintf("https://bgp.tools/search?q=%s", url.QueryEscape(organization)), organization, results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the HTML response
func getBGPToolsData(sourceURL string, organization string, results chan utils.Result) {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{})

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		gologger.Error().Msgf("%s", err)
	}

	// Parse html to get the needed data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// Example HTML tr entry:
		// <tr>
		//     <td class="nowrap"><img class="flag-img" src="/assets/flags/US.png"></td>
		//     <td><a href="?q=204.136.64.0%2F23">204.136.64.0/23</a></td>
		//     <td class="nowrap">Pepsi Cola Company / North America (PCNA)</td>
		// </tr>
		cidr := s.Find("td").Find("a").Text()
		name, _ := s.Find("td").Attr("class")

		if strings.Contains(name, "nowrap") {
			// Print if the description contains the organization and if the CIDR don't start with AS
			if strings.Contains(strings.ToLower(s.Text()), strings.ToLower(organization)) && !strings.HasPrefix(cidr, "AS") {
				if utils.CheckNumber(cidr) {
					results <- utils.Result{Value: cidr, Source: "bgptools"}
				}
			}
		}
	})
}
