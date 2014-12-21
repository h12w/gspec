// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"errors"
	"testing"

	exp "h12.me/gspec/expectation"
	ext "h12.me/gspec/extension"
	. "h12.me/gspec/reporter"
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
	expect := exp.Alias(exp.TFail(t.FailNow))
	r := &ReporterStub{}
	NewController(r).Start(Path{}, true, func(s S) {
		g := aliasGroup(s)
		g(func() {
			s.Fail(errors.New("err a"))
		})
	})
	expect("the root group", r.group).NotEqual(nil)
	expect(len(r.group.Children)).Equal(1)
	expect(r.group.Children[0].Error).Equal(errors.New("err a"))
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
	expect := exp.Alias(exp.TFail(t.FailNow))
	r := &ReporterStub{}
	ch := NewSS()
	NewController(r).Start(Path{}, true, func(s S) {
		g := aliasGroup(s)
		g(func() {
			defer func() {
				ch.Send("defer func")
			}()
			ch.Send("before FailNow")
			s.FailNow(errors.New("err a"))
			ch.Send("after FailNow")
		})
		g(func() {
			ch.Send("another test case")
		})
	})
	expect(r.group).NotEqual(nil)
	expect(len(r.group.Children)).Equal(2)
	expect(r.group.Children[0].Error).Equal(errors.New("err a"))
	expect(ch.Sorted()).Equal([]string{"another test case", "before FailNow", "defer func"})
}

/*
Scenario: test case panics
	Given a test case
	When it panics
	Then the error will be recorded and sent to the reporter
*/
func TestCasePanics(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	r := &ReporterStub{}
	NewController(r).Start(Path{}, true, func(s S) {
		g := aliasGroup(s)
		g(func() {
			panic("panic error")
		})
	})
	expect(r.group).NotEqual(nil)
	expect(len(r.group.Children)).Equal(1)
	expect(r.group.Children[0].Error).Equal(errors.New("panic error"))
}

/*
Scenario: Plain text progress indicator
	Given a nested test group with 5 leaves
	When the tests finished but 3 of test panics
	Then I should see 2 dots with 3 F: "F.F.F"
*/
func Test3Pass2Fail(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	var buf bytes.Buffer
	NewController(NewTextProgresser(&buf)).Start(Path{}, true, func(s S) {
		g := s.Alias("")
		g("a", func() {
			g("a-b", func() {
			})
			g("a-c", func() {
				g("a-c-d", func() {
					s.Fail(errors.New("err: a-c-d"))
				})
			})
			g("a-e", func() {
				g("a-e-f", func() {
					s.Fail(errors.New("err: a-e-f"))
				})
				g("a-e-g", func() {
					s.Fail(errors.New("123"))
				})
			})
		})
		g("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(out).HasPrefix("^")
	expect(out).HasSuffix("$\n")
	out = out[1 : len(out)-2]
	expect(sortBytes(out)).Equal("..FFF")
}

/*
Scenario: Pending test group
	Given nested test groups
	When one of them has no test closure
	Then a PendingError is set.
*/
func TestPending(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	r := &ReporterStub{}
	NewController(r).Start(Path{}, true, func(s S) {
		g := aliasGroup(s)
		g(func() {
		})
		g(nil)
		g(func() {
		})
	})
	expect(r.group).NotEqual(nil)
	expect(len(r.group.Children)).Equal(3)
	expect(r.group.Children[1].Error).Equal(&ext.PendingError{})
}
