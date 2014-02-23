// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
	"reflect"
	"strings"

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

// HasPrefix checks if the actual value has a prefix of expected value.
func HasPrefix(actual, expected interface{}) error {
	a, ok := actual.(string)
	if !ok {
		return errors.Expect("actual value is a string.")
	}
	e, ok := expected.(string)
	if !ok {
		return errors.Expect("expected value is a string.")
	}
	if strings.HasPrefix(a, e) {
		return nil
	}
	return errors.Compare(actual, expected, "to has the prefix of")
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

// HasPrefix is the fluent method for checker HasPrefix.
func (a *Actual) HasPrefix(expected interface{}) {
	a.To(HasPrefix, expected)
}

// Panic is the fluent method for checker Panic.
func (a *Actual) Panic() {
	a.To(Panic, nil)
}
