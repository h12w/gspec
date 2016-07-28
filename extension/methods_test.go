// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"errors"
	"testing"

	exp "h12.me/gspec/expectation"
)

/*
Story: A developer views the test result as a specification
*/

/*
Scenario: Traversal over all test cases in the tree of nested test group
	Given a root TestGroup
	When the For method is called
	Then each test case will get visited once
*/
func TestTestCaseTraversal(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	g := &TestGroup{
		Description: "a",
		Children: TestGroups{
			&TestGroup{
				Description: "b",
				Children: TestGroups{
					&TestGroup{
						Description: "c",
					},
					&TestGroup{
						Description: "d",
						Error:       errors.New("err"),
					},
				},
			},
			&TestGroup{
				Description: "e",
			},
		},
	}
	cases := []string{}
	g.For(func(path TestGroups) bool {
		s := ""
		for _, g := range path {
			s += g.Description
			if g.Error != nil {
				s += ":" + g.Error.Error()
			}
		}
		cases = append(cases, s)
		return true
	})
	expect(cases).Equal([]string{
		"abc",
		"abd:err",
		"ae",
	})

	cases = []string{}
	g.For(func(path TestGroups) bool {
		s := ""
		hasError := false
		for _, g := range path {
			s += g.Description
			if g.Error != nil {
				s += ":" + g.Error.Error()
				hasError = true
			}
		}
		cases = append(cases, s)
		return !hasError
	})
	expect(cases).Equal([]string{
		"abc",
		"abd:err",
	})

}

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestGroupStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	s := groupStack{}
	s.push(&TestGroup{Description: "a"})
	s.push(&TestGroup{Description: "b"})
	expect(s.a).Equal(
		TestGroups{
			{Description: "a"},
			{Description: "b"},
		})
	g := s.pop()
	expect(s.a).Equal(
		TestGroups{
			{Description: "a"},
		})
	expect(g).Equal(&TestGroup{Description: "b"})
	g = s.pop()
	expect(g).Equal(&TestGroup{Description: "a"})
	expect(s.a).Equal(TestGroups{})
	expect(func() { s.pop() }).Panic()
}
