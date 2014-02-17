// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
	"reflect"
)

var (
	// Sprint is used by checker to format actual and expected variables to
	// strings. It has a default implementation and can be replaced with an
	// external function.
	Sprint = func(v interface{}) string {
		return fmt.Sprintf("%#v", v)
	}
)

// Checker is the type of function that checks between actual and expected value
// then returns an Error if the expectation fails.
type Checker func(actual, expected interface{}) *Error

// Equal checks for the equality of contents and is tolerant of type differences.
func Equal(actual, expected interface{}) *Error {
	if reflect.DeepEqual(actual, expected) {
		return nil
	}
	if fmt.Sprint(actual) == fmt.Sprint(expected) {
		return nil
	}
	return &Error{fmt.Sprintf("\nExpect\n    %s\nequals\n    %s\n",
		Sprint(actual), Sprint(expected))}
}

// Panic checks if a function panics.
func Panic(actual, expected interface{}) (ret *Error) {
	f, ok := actual.(func())
	if !ok {
		ret = &Error{"the argument of Panic has to be a function of type func()."}
	}
	defer func() {
		if err := recover(); err == nil {
			ret = &Error{"Expect panicking but not occurred."}
		}
	}()
	f()
	return nil
}

// Equal is the fluent method for checker Equal.
func (a *Actual) Equal(expected interface{}) {
	a.To(Equal, expected)
}

// Panic is the fluent method for checker Panic.
func (a *Actual) Panic() {
	a.To(Panic, nil)
}
