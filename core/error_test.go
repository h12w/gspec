// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"errors"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
	. "github.com/hailiang/gspec/reporter"
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
	r := &ReporterStub{}
	NewScheduler(&TStub{}, r).Start(false, func(s S) {
		do := aliasGroup(s)
		do(func() {
			s.Fail(errors.New("err a"))
		})
	})
	expect(len(r.groups)).Equal(1)
	if len(r.groups) == 1 {
		expect(r.groups[0].Error).Equal(errors.New("err a"))
	}
}

/*
Scenario: FailNow
	Given the FailNow method of S
	When it is called with an error object
	Then the error Will be recorded and sent to the reporter
	And the defer functions get called
	And the rest of the test cases continue to run
*/
func TestFailNow(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	r := &ReporterStub{}
	ch := NewSChan()
	NewScheduler(&TStub{}, r).Start(false, func(s S) {
		do := aliasGroup(s)
		do(func() {
			defer func() {
				ch.Send("defer func")
			}()
			ch.Send("before FailNow")
			s.FailNow(errors.New("err a"))
			ch.Send("after FailNow")
		})
		do(func() {
			ch.Send("another test case")
		})
	})
	expect(len(r.groups)).Equal(2)
	if len(r.groups) > 0 {
		expect(r.groups[0].Error).Equal(errors.New("err a"))
	}
	expect(ch.Sorted()).Equal([]string{"another test case", "before FailNow", "defer func"})
}

/*
Scenario: test case panics
	Given a test case
	When it panics
	Then the error will be recorded and sent to the reporter
*/
func TestCasePanics(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	r := &ReporterStub{}
	NewScheduler(&TStub{}, r).Start(false, func(s S) {
		do := aliasGroup(s)
		do(func() {
			panic("panic error")
		})
	})
	expect(len(r.groups)).Equal(1)
	if len(r.groups) == 1 {
		expect(r.groups[0].Error).Equal(errors.New("panic error"))
	}
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
	NewScheduler(&TStub{}, NewTextProgresser(&buf)).Start(false, func(s S) {
		do := s.Alias("")
		do("a", func() {
			do("a-b", func() {
			})
			do("a-c", func() {
				do("a-c-d", func() {
					s.Fail(errors.New("err: a-c-d"))
				})
			})
			do("a-e", func() {
				do("a-e-f", func() {
					s.Fail(errors.New("err: a-e-f"))
				})
				do("a-e-g", func() {
					s.Fail(errors.New("123"))
				})
			})
		})
		do("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(out).HasPrefix("^")
	expect(out).HasSuffix("$\n")
	out = out[1 : len(out)-2]
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
	mt := &TStub{}
	r := &ReporterStub{}
	NewScheduler(mt, r).Start(false, func(s S) {
		do := aliasGroup(s)
		do(func() {
			s.Fail(errors.New(""))
		})
	})
	expect(mt.s).Equal("Fail.")
}
