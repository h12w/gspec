// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"fmt"
	"os"

	"github.com/hailiang/gspec/core"
	ext "github.com/hailiang/gspec/extension"
	"github.com/hailiang/gspec/reporter"
)

var (
	// Reporters are the test reporters used during the test.
	Reporters = []ext.Reporter{
		reporter.NewTextProgresser(os.Stdout),
		reporter.NewTextReporter(os.Stdout),
	}
	testFunctions []core.TestFunc
)

// Add GSpec test functions to the global test suite.
// Return value has no meaning, allowing it to be called in global scope.
func Add(fs ...core.TestFunc) int {
	testFunctions = append(testFunctions, fs...)
	return 0
}

// Run all tests in the global test suite.
func Run(t core.T, sequential bool) {
	s := core.NewController(t, Reporters...)
	err := s.Start(sequential, testFunctions...)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
