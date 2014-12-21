// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package expectation provides expectation (assertion) helpers.
*/
package expectation // import "h12.me/gspec/expectation"

import "fmt"

// FailFunc is the function type that is used to notify expectation failures.
type FailFunc func(error)

// Actual provides checking methods for an actual value in an expectation.
type Actual struct {
	v    interface{}
	name string
	skip int
	fail FailFunc
}

// To is a general method for checking an expectation.
func (a *Actual) To(check Checker, expected interface{}) {
	a.to(check, expected, 1)
}

func (a *Actual) to(check Checker, expected interface{}, skip int) {
	if err := check(a.v, expected, a.name, a.skip+skip+1); err != nil {
		a.fail(err)
	}
}

// ExpectFunc is the type of function that returns an Actual object given an
// actual value or a name and an actual value.
type ExpectFunc func(actual ...interface{}) *Actual

// Alias registers a fail function and returns an ExpectFunc.
// The optional skip parameter is used to skip extra function calls in the stack
// trace in case the ExpectFunc is further wrapped by another function.
func Alias(fail FailFunc, skip ...int) ExpectFunc {
	_skip := getSkip(skip)
	return func(s ...interface{}) *Actual {
		actual, name := getActual(s)
		return &Actual{actual, name, _skip, fail}
	}
}

func getActual(s []interface{}) (actual interface{}, name string) {
	if len(s) == 1 {
		return s[0], ""
	} else if len(s) == 2 {
		if name, ok := s[0].(string); ok {
			return s[1], name
		}
	}
	panic(fmt.Errorf("invalid argument actual %v", s))
}

func getSkip(s []int) int {
	if len(s) == 0 {
		return 0
	} else if len(s) == 1 {
		return s[0]
	}
	panic(fmt.Errorf("invalid argument skip %v", s))
}
