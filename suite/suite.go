// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hailiang/gspec/core"
	ext "github.com/hailiang/gspec/extension"
	"github.com/hailiang/gspec/reporter"
)

var (
	// Reporters are the test reporters used during the test.
	Reporters = []ext.Reporter{
		reporter.NewTextProgresser(os.Stdout),
		reporter.NewTextReporter(os.Stdout, Verbose()),
	}

	testFunctions []core.TestFunc
	globalConfig  config
)

type config struct {
	focus      core.Path
	concurrent bool
}

func init() {
	flag.Var(&globalConfig.focus, "focus", "tell GSpec to run only focused test case.")
	flag.BoolVar(&globalConfig.concurrent, "concurrent", false, "tell GSpec to run concurrently, false by default.")
}

// Add GSpec test functions to the global test suite.
// Return value has no meaning, allowing it to be called in global scope.
func Add(fs ...core.TestFunc) int {
	testFunctions = append(testFunctions, fs...)
	return 0
}

// T is an interface that allows a testing.T to be passed without depending on
// the testing package.
type T interface {
	Fail()
}

// Run all tests in the global test suite.
func Run(t T) {
	if globalConfig.concurrent {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	reporters := append(Reporters, reporter.NewFailReporter(t))
	s := core.NewController(reporters...)
	err := s.Start(globalConfig.focus, globalConfig.concurrent, testFunctions...)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

// Verbose returns value of "test.v" flag without depending on the testing
// package.
func Verbose() bool {
	if f := flag.Lookup("test.v"); f != nil {
		return f.Value.String() == "true"
	}
	return false
}
