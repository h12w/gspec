// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import "fmt"

// Actual provides checking methods for an actual value in an expectation.
type Actual struct {
	v    interface{}
	fail func(error)
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

// Alias registers a fail function and returns an ExpectFunc.
func Alias(fail func(error)) ExpectFunc {
	return func(actual interface{}) *Actual {
		return &Actual{actual, fail}
	}
}

// T is a subset of testing.T used in this package.
type T interface {
	Fail()
}

// AliasForT registers T as the fail handler and returns an ExpectFunc.
func AliasForT(t T) ExpectFunc {
	return Alias(func(err error) {
		fmt.Println(decorate(err.Error(), 4))
		t.Fail()
	})
}
