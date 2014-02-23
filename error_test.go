// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"bytes"
	"errors"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: A developer deals with test errors

As a developer
I want to see a helpful error message when an error occurs
So that I can easily find the cause and fix it
*/

/*
Scenario: test case fails
	Given a Fail method of S
	When it is called with an error object
	Then the error will be recorded and sent to the reporter
*/
func TestCaseFails(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	r := &MockReporter{}
	NewScheduler(&MockT{}, r).Start(false, func(s S) {
		do := aliasDo(s)
		do(func() {
			s.Fail(errors.New("err a"))
		})
	})
	expect(r.groups).Equal([]*TestGroup{
		{
			Error: errors.New("err a"),
		},
	})
}

/*
Scenario: Plain text progress indicator
	Given a nested test group with 5 leaves
	When the tests finished but 3 of test panics
	Then I should see 2 dots with 3 F: "F.F.F"
*/
func Test3Pass2Fail(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	var buf bytes.Buffer
	NewScheduler(&MockT{}, NewTextProgresser(&buf)).Start(false, func(s S) {
		do := s.Alias("")
		do("a", func() {
			do("a-b", func() {
			})
			do("a-c", func() {
				do("a-c-d", func() {
					panic("err: a-c-d")
				})
			})
			do("a-e", func() {
				do("a-e-f", func() {
					panic(errors.New("err: a-e-f"))
				})
				do("a-e-g", func() {
					panic(123)
				})
			})
		})
		do("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(sortBytes(out)).Equal("..FFF")
}

/*
Scenario: notify testing.T
	Given a Fail method of S
	When it is called
	Then testing.T.Fail should be called
*/
func TestNotifyT(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	mt := &MockT{}
	r := &MockReporter{}
	NewScheduler(mt, r).Start(false, func(s S) {
		do := aliasDo(s)
		do(func() {
			s.Fail(errors.New(""))
		})
	})
	expect(mt.s).Equal("Fail.")
}
