// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/dhn/spk/runner"
	"github.com/dhn/spk/utils"
)

func main() {
	options := utils.ParseOptions()

	if options.SearchString != "" {
		runner.Whois(options.SearchString, *options)
	}
}
