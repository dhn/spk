// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runner

import (
	"github.com/dhn/spk/sources"
	"github.com/dhn/spk/utils"
)

// Get the net-ranges based on the "descr" field in the whois response
func Whois(organization string, options utils.Options) {
	netranges := utils.MergeChannels(
		sources.GetSPKData(organization),
		sources.GetAPNICData(organization),
		sources.GetBGPData(organization),
		sources.GetBGPToolsData(organization),
		sources.GetRipeData(organization),
		sources.GetARINData(organization),
	)
	netranges = utils.RemoveDuplicates(netranges)

	utils.PrintResults(options.JSON, netranges)
}
