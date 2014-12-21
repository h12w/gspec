// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
	"reflect"

	"h12.me/gspec/errors"
)

// Checker is the type of function that checks between actual and expected value
// then returns an Error if the expectation fails.
type Checker func(actual, expected interface{}, name string, skip int) error

// Equal checks for the equality of contents and is tolerant of type differences.
func Equal(actual, expected interface{}, name string, skip int) error {
	if reflect.DeepEqual(actual, expected) {
		return nil
	}
	if fmt.Sprint(actual) == fmt.Sprint(expected) {
		return nil
	}
	return errors.Compare(actual, expected, "to equal", name, skip+1)
}

// NotEqual is the reverse of Equal.
func NotEqual(actual, expected interface{}, name string, skip int) error {
	if Equal(actual, expected, name, skip+1) != nil {
		return nil
	}
	return errors.Compare(actual, expected, "not to equal", name, skip+1)
}

// Panic checks if a function panics.
func Panic(actual, expected interface{}, name string, skip int) (ret error) {
	f, ok := actual.(func())
	if !ok {
		ret = errors.Expect("the argument of Panic has to be a function of type func().", skip)
	}
	defer func() {
		if err := recover(); err == nil {
			ret = errors.Expect("panicking", skip+1)
		}
	}()
	f()
	return nil
}

// IsType checks if the actual value is of the same type as the expected value.
func IsType(actual, expected interface{}, name string, skip int) error {
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return errors.Compare(actual, expected, "to have type of", name, skip+1)
	}
	return nil
}

// Equal is the fluent method for checker Equal.
func (a *Actual) Equal(expected interface{}) {
	a.to(Equal, expected, 1)
}

// NotEqual is the fluent method for checker NotEqual.
func (a *Actual) NotEqual(expected interface{}) {
	a.to(NotEqual, expected, 1)
}

// Panic is the fluent method for checker Panic.
func (a *Actual) Panic() {
	a.to(Panic, nil, 1)
}

// IsType is the fluent method for checker IsType.
func (a *Actual) IsType(expected interface{}) {
	a.to(IsType, expected, 1)
}
