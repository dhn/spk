// Copyright (c) 2022 dhn. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"sync"
	"regexp"

	"github.com/projectdiscovery/gologger"
)

// Print results based on the given parameters
func PrintResults(json bool, results <-chan Result) {
	if json {
		WriteJSON(results)
	} else {
		for result := range results {
			gologger.Silent().Msg(result.Value)
		}
	}
}

// Remove duplicates from a channel and return a channel from type Result
func RemoveDuplicates(input <-chan Result) <-chan Result {
	output := make(chan Result)

	go func() {
		set := make(map[Result]struct{})
		for index := range input {
			if _, ok := set[index]; !ok {
				set[index] = struct{}{}
				output <- index
			}
		}
		close(output)
	}()

	return output
}

// Merge multiple channels from type Result
func MergeChannels(channels ...<-chan Result) <-chan Result {
	out := make(chan Result)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		go func(channel <-chan Result) {
			for value := range channel {
				out <- value
			}
			wg.Done()
		}(channel)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func CheckNumber(cidr string) bool {
	reg := regexp.MustCompile(`^\d.*`)
	if reg.MatchString(cidr) {
		return true
	} else {
		return false
	}
}
