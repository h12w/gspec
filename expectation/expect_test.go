// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"

	ge "github.com/hailiang/gspec/error"
)

func TestExpectTo(t *testing.T) {
	m, expect := mockExpect()
	expect(nil).To(func(actual, expected interface{}, name string, skip int) error {
		return ge.Expect("x", skip+1)
	}, nil)
	e, ok := m.Error()
	if !ok {
		t.Errorf("Expect error to be type ExpectError, got %v.", e)
	}
	if e.Text != "x" {
		t.Errorf("Expect error message %v to be x.", e.Text)
	}
	if e.Pos.BasePath() != "expect_test.go" {
		t.Errorf("Expect error position in file expect_test.go, but got %s.",
			e.Pos.BasePath())
	}
	if e.Pos.Line != 17 {
		t.Errorf("Expect error position at line 17, but got %d.", e.Pos.Line)
	}
}

type expectTestCase struct {
	msg string
	f   func(ExpectFunc)
	r   bool // true if expectation is mean to succeed.
}

func testExpectations(t *testing.T, expectTestCases []expectTestCase) {
	for _, c := range expectTestCases {
		m, expect := mockExpect()
		c.f(expect)
		if c.r != (m.err == nil) {
			t.Error(c.msg)
		}
	}
}

type expectMock struct {
	err error
}

func mockExpect() (*expectMock, ExpectFunc) {
	m := &expectMock{}
	expect := Alias(func(err error) {
		m.err = err
	})
	return m, expect
}

func (m *expectMock) Error() (*ge.ExpectError, bool) {
	e, ok := m.err.(*ge.ExpectError)
	return e, ok
}
