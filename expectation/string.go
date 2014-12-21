// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"strings"

	ge "h12.me/gspec/errors"
)

// HasPrefix checks if the actual value has a prefix of expected value.
func HasPrefix(actual, expected interface{}, name string, skip int) error {
	a, e, err := checkStringType(actual, expected, skip+1)
	if err != nil {
		return err
	}
	if strings.HasPrefix(a, e) {
		return nil
	}
	return ge.Compare(actual, expected, "to has the prefix of", name, skip+1)
}

// HasSuffix checks if the actual value has a suffix of expected value.
func HasSuffix(actual, expected interface{}, name string, skip int) error {
	a, e, err := checkStringType(actual, expected, skip+1)
	if err != nil {
		return err
	}
	if strings.HasSuffix(a, e) {
		return nil
	}
	return ge.Compare(actual, expected, "to has the suffix of", name, skip+1)
}

// Contains checks if the actual value contains expected value.
func Contains(actual, expected interface{}, name string, skip int) error {
	a, e, err := checkStringType(actual, expected, skip+1)
	if err != nil {
		return err
	}
	if strings.Contains(a, e) {
		return nil
	}
	return ge.Compare(actual, expected, "to contain", name, skip+1)
}

func checkStringType(actual, expected interface{}, skip int) (string, string, error) {
	a, ok := actual.(string)
	if !ok {
		return "", "", ge.Expect("actual value is a string.", skip+1)
	}
	e, ok := expected.(string)
	if !ok {
		return "", "", ge.Expect("expected value is a string.", skip+1)
	}
	return a, e, nil
}

// HasPrefix is the fluent method for checker HasPrefix.
func (a *Actual) HasPrefix(expected interface{}) {
	a.to(HasPrefix, expected, 1)
}

// HasSuffix is the fluent method for checker HasSuffix.
func (a *Actual) HasSuffix(expected interface{}) {
	a.to(HasSuffix, expected, 1)
}

// Contains is the fluent method for checker Contains.
func (a *Actual) Contains(expected interface{}) {
	a.to(Contains, expected, 1)
}
