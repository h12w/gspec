// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package error provides all types of test error.
*/
package error

import (
	"fmt"
)

// ExpectError is the base type of an expectation error.
type ExpectError struct {
	Pos  *Pos
	Text string
}

// Expect returns a new ExpectError object.
func Expect(text string, skip int) error {
	return &ExpectError{GetPos(skip + 1), text}
}

func (e *ExpectError) str(msg string) string {
	format := "expect %s."
	if startWithBreak(msg) {
		format = "expect%s."
	}
	return e.Pos.Decorate(fmt.Sprintf(format, msg), "")
}

// Error of ExpectError print the Text field.
func (e *ExpectError) Error() string {
	return e.str(e.Text)
}

// CompareError is the error of comparing two values.
type CompareError struct {
	ExpectError
	Actual, Expected interface{}
}

// Compare returns a new CompareError object.
func Compare(actual, expected interface{}, verb string, skip int) error {
	return &CompareError{ExpectError{GetPos(skip + 1), verb}, actual, expected}
}

func (e *CompareError) verb() string {
	return e.Text
}

// Error of CompareError formats an error message with the actual, expected
// value and the verb. When the actual value ends with break, it will add indent
// accordingly.
func (e *CompareError) Error() string {
	actual := Sprint(e.Actual)
	expect := Sprint(e.Expected)
	format := "%s %s %s"
	if endWithBreak(actual) {
		format = "%s%s%s"
		actual = Indent(actual, IndentString)
		expect = Indent(expect, IndentString)
	}
	return e.str(fmt.Sprintf(format, actual, e.verb(), expect))
}

func endWithBreak(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\n'
}

func startWithBreak(s string) bool {
	return len(s) > 0 && s[0] == '\n'
}
