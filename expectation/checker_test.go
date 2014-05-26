// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"
)

type testCase struct {
	actual, expected interface{}
	succeeded        bool
}

var equalCases = []testCase{
	{nil, nil, true},
	{nil, 1, false},
	{1, 1, true},
	{1, 2, false},
	{int32(1), int32(1), true},
	{1, int32(1), true},
	{1, "1", true},
	{[]int{1}, []int{1}, true},
}

func TestEqual(t *testing.T) {
	for _, c := range equalCases {
		m, expect := mockExpect()
		expect(c.actual).Equal(c.expected)
		if (m.err != nil) == c.succeeded {
			t.Errorf("Equal test: %v", c)
		}
	}
}

func TestNotEqual(t *testing.T) {
	for _, c := range equalCases {
		m, expect := mockExpect()
		expect(c.actual).NotEqual(c.expected)
		if (m.err != nil) == !c.succeeded {
			t.Errorf("NotEqual test: %v", c)
		}
	}
}

var panicTestCases = []expectTestCase{
	{`expect(func() { panic("") }).Panic()`, func(expect ExpectFunc) { expect(func() { panic("") }).Panic() }, true},
	{`expect(func() {}).Panic()`, func(expect ExpectFunc) { expect(func() {}).Panic() }, false},
	{`wrong signature: expect(func(int) { panic("") }).Panic()`, func(expect ExpectFunc) { expect(func(int) { panic("") }).Panic() }, false},
}

func TestPanic(t *testing.T) {
	testExpectations(t, panicTestCases)
}

var isTypeTestCases = []expectTestCase{
	{`expect(int(0)).IsType(int(1))`, func(expect ExpectFunc) { expect(int(0)).IsType(int(1)) }, true},
	{`expect(int(0)).IsType(bool(true))`, func(expect ExpectFunc) { expect(int(0)).IsType(bool(true)) }, false},
}

func TestIsType(t *testing.T) {
	testExpectations(t, isTypeTestCases)
}
