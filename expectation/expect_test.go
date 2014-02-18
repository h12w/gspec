// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"github.com/hailiang/gspec/errors"
	"testing"
)

type testCase struct {
	actual, expected interface{}
	succeeded        bool
}

func TestEqual(t *testing.T) {
	cases := []testCase{
		{nil, nil, true},
		{nil, 1, false},
		{1, 1, true},
		{1, 2, false},
		{int32(1), int32(1), true},
		{1, int32(1), true},
		{1, "1", true},
		{[]int{1}, []int{1}, true},
	}

	for _, c := range cases {
		m, expect := mockExpect()
		expect(c.actual).Equal(c.expected)
		if (m.err != nil) == c.succeeded {
			t.Errorf("Equal test fails: %v", c)
		}
	}
}

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
