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
	expect(nil).To(func(actual, expected interface{}) error {
		return errors.Expect("x")
	}, nil)
	if e, ok := m.Error(); !ok || e.Text != "x" {
		t.Errorf("Expect error message %v to be x", e.Text)
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
