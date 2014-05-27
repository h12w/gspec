// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package expectation provides expectation (assertion) helpers.
*/
package expectation

// FailFunc is the function type that is used to notify expectation failures.
type FailFunc func(error)

// Actual provides checking methods for an actual value in an expectation.
type Actual struct {
	v    interface{}
	fail FailFunc
}

// To is a general method for checking an expectation.
func (a *Actual) To(check Checker, expected interface{}) {
	a.to(check, expected, 1)
}

func (a *Actual) to(check Checker, expected interface{}, skip int) {
	if err := check(a.v, expected, skip+1); err != nil {
		a.fail(err)
	}
}

// ExpectFunc is the type of function that returns an Actual object given an
// actual value.
type ExpectFunc func(actual interface{}) *Actual

// Alias registers a fail function and returns an ExpectFunc.
func Alias(fail FailFunc) ExpectFunc {
	return func(actual interface{}) *Actual {
		return &Actual{actual, func(e error) {
			fail(e)
		}}
	}
}
