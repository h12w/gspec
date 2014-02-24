// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"errors"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
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
	expect := exp.Alias(exp.TFail(t))
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
	g.For(func(path TestGroups) {
		s := ""
		for _, g := range path {
			s += g.Description
			if g.Error != nil {
				s += ":" + g.Error.Error()
			}
		}
		cases = append(cases, s)
	})
	expect(cases).Equal([]string{
		"abc",
		"abd:err",
		"ae",
	})
}
