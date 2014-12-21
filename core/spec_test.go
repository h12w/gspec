// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	exp "h12.me/gspec/expectation"
	. "h12.me/gspec/extension"
)

/*
Story: A developer write tests as specification

As a developer
I want to write my tests
So that I can get a structured, readable specification
*/

/*
Scenario: Attach customized alias name and description for each test group method.
	Given the S object provided by GSpec
	When I define a customized name for the test group method
	And use it to attach a description to the test group
	Then the alias name and description are combined and passed to the reporter
*/
func TestDescribeTests(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	r := &ReporterStub{}
	NewController(r).Start(Path{}, false, func(s S) {
		describe, context, it := s.Alias("describe"), s.Alias("context"), s.Alias("it")
		describe("a", func() {
			context("b", func() {
				it("c", func() {
				})
			})
		})
	})
	clearGroupForTest(r.group)
	expect(r.group.Children).Equal(TestGroups{
		{
			Description: "describe a",
			Children: TestGroups{
				{
					Description: "context b",
					Children: TestGroups{
						{Description: "it c"},
					},
				},
			},
		},
	})
}
