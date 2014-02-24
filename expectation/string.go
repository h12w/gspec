// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"strings"

	"github.com/hailiang/gspec/errors"
)

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

// HasSuffix checks if the actual value has a suffix of expected value.
func HasSuffix(actual, expected interface{}) error {
	a, ok := actual.(string)
	if !ok {
		return errors.Expect("actual value is a string.")
	}
	e, ok := expected.(string)
	if !ok {
		return errors.Expect("expected value is a string.")
	}
	if strings.HasSuffix(a, e) {
		return nil
	}
	return errors.Compare(actual, expected, "to has the suffix of")
}

// HasPrefix is the fluent method for checker HasPrefix.
func (a *Actual) HasPrefix(expected interface{}) {
	a.To(HasPrefix, expected)
}

// HasSuffix is the fluent method for checker HasSuffix.
func (a *Actual) HasSuffix(expected interface{}) {
	a.To(HasSuffix, expected)
}
