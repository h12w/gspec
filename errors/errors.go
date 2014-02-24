// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
)

const skip = 5

// ExpectError is the base type of an expectation error.
type ExpectError struct {
	Pos  *Pos
	Text string
}

// Expect returns a new ExpectError object.
func Expect(text string) error {
	return &ExpectError{GetPos(skip), text}
}

func (e *ExpectError) str(msg string) string {
	return e.Pos.Decorate(fmt.Sprintf("expect %s.", msg))
}

// Error implements error interface.
func (e *ExpectError) Error() string {
	return e.str(e.Text)
}

// CompareError is the error of comparing two values.
type CompareError struct {
	ExpectError
	Actual, Expected interface{}
}

// Compare returns a new CompareError object.
func Compare(actual, expected interface{}, verb string) error {
	return &CompareError{ExpectError{GetPos(skip), verb}, actual, expected}
}

func (e *CompareError) verb() string {
	return e.Text
}

// Error implements error interface.
func (e *CompareError) Error() string {
	return e.str(fmt.Sprintf("%v %s %v", e.Actual, e.verb(), e.Expected))
}
