// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"

	"h12.me/gspec/errors"
)

func TestExpectTo(t *testing.T) {
	m, expect := mockExpect()
	expect(nil).To(func(actual, expected interface{}, name string, skip int) error {
		return errors.Expect("x", skip+1)
	}, nil)
	if m.err == nil {
		t.Error("Expect error but got nil.")
	}
	if m.err.Error() != "expect_test.go:17: expect x." {
		t.Errorf("Got error message %v.", m.err.Error())
	}
}

func TestExpectName(t *testing.T) {
	m, expect := mockExpect()
	expect("integer", 0).Equal(1)
	if m.err == nil {
		t.Error("Expect error but got nil.")
	}
	if m.err.Error() != "expect_test.go:28: expect integer 0 to equal 1." {
		t.Errorf("Got error message %v.", m.err.Error())
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
