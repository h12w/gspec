// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
	"reflect"

	"github.com/hailiang/gspec/errors"
)

// Checker is the type of function that checks between actual and expected value
// then returns an Error if the expectation fails.
type Checker func(actual, expected interface{}) error

// Equal checks for the equality of contents and is tolerant of type differences.
func Equal(actual, expected interface{}) error {
	if reflect.DeepEqual(actual, expected) {
		return nil
	}
	if fmt.Sprint(actual) == fmt.Sprint(expected) {
		return nil
	}
	return errors.Compare(actual, expected, "to equal")
}

// NotEqual is the reverse of Equal.
func NotEqual(actual, expected interface{}) error {
	if Equal(actual, expected) != nil {
		return nil
	}
	return errors.Compare(actual, expected, "not to equal")
}

// Panic checks if a function panics.
func Panic(actual, expected interface{}) (ret error) {
	f, ok := actual.(func())
	if !ok {
		ret = errors.Expect("the argument of Panic has to be a function of type func().")
	}
	defer func() {
		if err := recover(); err == nil {
			ret = errors.Expect("panicking")
		}
	}()
	f()
	return nil
}

// Equal is the fluent method for checker Equal.
func (a *Actual) Equal(expected interface{}) {
	a.To(Equal, expected)
}

// NotEqual is the fluent method for checker NotEqual.
func (a *Actual) NotEqual(expected interface{}) {
	a.To(NotEqual, expected)
}

// Panic is the fluent method for checker Panic.
func (a *Actual) Panic() {
	a.To(Panic, nil)
}
