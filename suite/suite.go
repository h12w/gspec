// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"flag"
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
	globalConfig  core.Config
)

func init() {
	flag.Var(&globalConfig.Focus, "focus", "test case id to select one test case to run")
}

// Add GSpec test functions to the global test suite.
// Return value has no meaning, allowing it to be called in global scope.
func Add(fs ...core.TestFunc) int {
	testFunctions = append(testFunctions, fs...)
	return 0
}

// T is an interface that allows a testing.T to be passed to GSpec.
type T interface {
	Fail()
	Parallel()
}

// Run all tests in the global test suite.
func Run(t T, concurrent bool) {
	if concurrent {
		t.Parallel()
	}
	fr := reporter.NewFailReporter(t)
	s := core.NewController(&globalConfig, append(Reporters, fr)...)
	err := s.Start(concurrent, testFunctions...)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
