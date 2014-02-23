// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"os"

	"github.com/hailiang/gspec"
)

var (
	// Reporters are the test reporters used during the test.
	Reporters = []gspec.Reporter{
		gspec.NewTextProgresser(os.Stdout),
		gspec.NewTextReporter(os.Stdout),
	}
	testFunctions []gspec.TestFunc
)

// Add GSpec test functions to the global test suite.
// Return value has no meaning, allowing it to be called in global scope.
func Add(fs ...gspec.TestFunc) int {
	testFunctions = append(testFunctions, fs...)
	return 0
}

// Run all tests in the global test suite.
func Run(t gspec.T, sequential bool) {
	s := gspec.NewScheduler(t, Reporters...)
	s.Start(sequential, testFunctions...)
}
