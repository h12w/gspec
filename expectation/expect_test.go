// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"

	"github.com/hailiang/gspec/errors"
)

func TestExpectTo(t *testing.T) {
	m, expect := mockExpect()
	expect(nil).To(func(actual, expected interface{}, skip int) error {
		return errors.Expect("x", skip+1)
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

func (m *expectMock) Error() (*errors.ExpectError, bool) {
	e, ok := m.err.(*errors.ExpectError)
	return e, ok
}
