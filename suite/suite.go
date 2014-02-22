// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"os"
	"github.com/hailiang/gspec"
)

var (
	// Reporter is the test reporter used during the test.
	Reporter = gspec.NewTextReporter(os.Stdout)
	testFunctions []gspec.TestFunc
)

// T is subset of testing.T
type T interface {
	Fail()
}

// Add GSpec test functions to the global test suite.
func Add(fs ...gspec.TestFunc) {
	testFunctions = append(testFunctions, fs...)
}

// Run all tests in the global test suite.
func Run(t T, sequential bool) {
	s := gspec.NewScheduler(Reporter)
	s.Start(sequential, testFunctions...)
}


