// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec // import "h12.me/gspec"

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"h12.me/gspec/core"
	"h12.me/gspec/errors"
	"h12.me/gspec/expectation"
	ext "h12.me/gspec/extension"
	"h12.me/gspec/reporter"
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

// TestFunc is a trivial wrapper to core.TestFunc.
type TestFunc func(S)

// S is a trivial wrapper to core.S.
type S struct {
	core.S
}

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
func Add(fs ...TestFunc) int {
	for _, f := range fs {
		testFunctions = append(testFunctions, func(s core.S) { f(S{s}) })
	}
	return 0
}

// T is an interface that allows a testing.T to be passed without depending on
// the testing package.
type T interface {
	Fail()
}

// Test method runs all tests in the global test suite.
func Test(t T) {
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

// Alias2 accepts 2 names and returns 2 alias DescFuncs.
func Alias2(n1, n2 string, s S) (_, _ core.DescFunc) {
	return s.Alias(n1), s.Alias(n2)
}

// Alias3 accepts 3 names and returns 3 alias DescFuncs.
func Alias3(n1, n2, n3 string, s S) (_, _, _ core.DescFunc) {
	return s.Alias(n1), s.Alias(n2), s.Alias(n3)
}

// Alias4 accepts 4 names and returns 4 alias DescFuncs.
func Alias4(n1, n2, n3, n4 string, s S) (_, _, _, _ core.DescFunc) {
	return s.Alias(n1), s.Alias(n2), s.Alias(n3), s.Alias(n4)
}

// Alias5 accepts 5 names and returns 5 alias DescFuncs.
func Alias5(n1, n2, n3, n4, n5 string, s S) (_, _, _, _, _ core.DescFunc) {
	return s.Alias(n1), s.Alias(n2), s.Alias(n3), s.Alias(n4), s.Alias(n5)
}

// Expect is a trivial wrapper of expectation.Alias for GSpec or Go tests.
func Expect(fail interface{}, skip ...int) expectation.ExpectFunc {
	if f, ok := fail.(func(error)); ok {
		return expectation.Alias(f, skip...)
	} else if f, ok := fail.(func()); ok {
		return expectation.Alias(expectation.TFail(f), skip...)
	}
	panic("argument fail should be either an expectation.FailFunc or a func()")
}

// SetSprint is a trivial wrapper to set error.Sprint.
func SetSprint(sprint func(interface{}) string) {
	errors.Sprint = sprint
}
