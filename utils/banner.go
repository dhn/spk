// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import "github.com/projectdiscovery/gologger"

const banner = `
               _    
     ___ _ __ | | __
    / __| '_ \| |/ /
    \__ \ |_) |   < 
    |___/ .__/|_|\_\
         |_|         

`

// Version is the current version of dnsx
const Version = `0.0.3`

// showBanner is used to show the banner to the user
func ShowBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\tmade with <3\n\n")
}
