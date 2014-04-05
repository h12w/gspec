// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
//	"runtime"
)

// FailFunc is the function type that is used to notify expectation failures.
type FailFunc func(error)

// Actual provides checking methods for an actual value in an expectation.
type Actual struct {
	v    interface{}
	fail FailFunc
}

// To is a general method for checking an expectation.
func (a *Actual) To(check Checker, expected interface{}) {
	if err := check(a.v, expected); err != nil {
		a.fail(err)
	}
}

// ExpectFunc is the type of function that returns an Actual object given an
// actual value.
type ExpectFunc func(actual interface{}) *Actual

// Alias registers a fail function and returns an ExpectFunc. The fail function
// needs only to record the error and the returned ExpectFunc will terminate the
// goroutine as soon as jn expectation failure occurs.
func Alias(fail FailFunc) ExpectFunc {
	return func(actual interface{}) *Actual {
		return &Actual{actual, func(e error) {
			fail(e)
//			runtime.Goexit()
		}}
	}
}

// T is a subset of testing.T used in this package.
type T interface {
	Fail()
	FailNow()
}

// TFail return the FailFunc for testing.T.Fail
func TFail(t T) FailFunc {
	return func(err error) {
		t.Fail()
		fmt.Println(err.Error())
	}
}

// TFailNow return the FailFunc for testing.T.FailNow
func TFailNow(t T) FailFunc {
	return func(err error) {
		t.FailNow()
		fmt.Println(err.Error())
	}
}
