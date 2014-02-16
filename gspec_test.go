// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"bytes"
	exp "github.com/hailiang/gspec/expectation"
	"testing"
)

/*
TODO:
* report failure location
*/

/*
Story: Dveloper verify tests

As a developer
I want to run my tests
So that I can verify the result of the tests
*/

/*
Scenario: Plain text progress indicator
	Given a nested test group with 5 leaves
	When the tests finished without any error
	Then I should see 5 dots: "....."
*/
func Test5Pass(t *testing.T) {
	expect := exp.AliasForT(t)
	var buf bytes.Buffer
	NewScheduler(NewTextReporter(&buf)).Start(false, func(s S) {
		do := s.Alias("")
		do("a", func() {
			do("a-b", func() {
			})
			do("a-c", func() {
				do("a-c-d", func() {
				})
			})
			do("a-e", func() {
				do("a-e-f", func() {
				})
				do("a-e-g", func() {
				})
			})
		})
		do("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(sortBytes(out)).Equal(".....")
}

/*
Scenario: Plain text progress indicator
	Given a nested test group with 5 leaves
	When the tests finished but 1 of test panics
	Then I should see 4 dots with 1 F: "..F.."
*/
func Test4Pass1Fail(t *testing.T) {
	expect := exp.AliasForT(t)
	var buf bytes.Buffer
	NewScheduler(NewTextReporter(&buf)).Start(false, func(s S) {
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
				})
				do("a-e-g", func() {
				})
			})
		})
		do("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(sortBytes(out)).Equal("....F")
}

/*
Scenario: Plain text progress indicator
	Given a nested test group with 5 leaves
	When the tests finished but 2 of test panics
	Then I should see 3 dots with 2 F: "..F.."
*/
func Test3Pass2Fail(t *testing.T) {
	expect := exp.AliasForT(t)
	var buf bytes.Buffer
	NewScheduler(NewTextReporter(&buf)).Start(false, func(s S) {
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
				})
				do("a-e-g", func() {
					panic("err: a-e-g")
				})
			})
		})
		do("h", func() {
		})
	})
	out, _ := buf.ReadString('\n')
	expect(sortBytes(out)).Equal("...FF")
}
