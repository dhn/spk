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

// GetBGPData function returns all netranges based on the given organization name
func GetBGPData(organization string) <-chan utils.Result {
	results := make(chan utils.Result)

	go func() {
		getBGPViewData(fmt.Sprintf("https://bgpview.io/search/%s#results-v4",
			url.QueryEscape(organization)), results)
		close(results)
	}()

	return results
}

// Send a HTTP request and parse the HTML response
func getBGPViewData(sourceURL string, results chan utils.Result) {
	resp := utils.GetHTTPRequest(sourceURL, map[string]string{})

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		gologger.Error().Msgf("%s", err)
	}

	// Parse html to get the needed data
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "https://bgpview.io/prefix/") {
			cidr := s.Text()
			results <- utils.Result{Value: cidr, Source: "bgpview"}
		}
	})
}
