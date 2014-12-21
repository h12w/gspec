// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"testing"

	exp "h12.me/gspec/expectation"
	. "h12.me/gspec/extension"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

/*
Scenario: reconstruct nested test group to a tree
	Given a treeCollector
	When it's groupStart method is called like real tests running
	Then it is able to reconstruct the tree structure
*/
func TestTreeCollector(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	co := newCollector(&ReporterStub{})
	r := &TestGroup{}
	a := &TestGroup{
		Description: "a",
	}
	b := &TestGroup{
		Description: "b",
	}
	c := &TestGroup{
		Description: "c",
	}
	d := &TestGroup{
		Description: "d",
	}
	z := &TestGroup{
		Description: "z",
	}
	co.groupStart(r, Path{1})
	co.groupStart(a, Path{1, 2})
	co.groupStart(b, Path{1, 2, 3})
	co.groupStart(c, Path{1, 2, 3, 4})
	c.Error = errors.New("c err")
	co.groupStart(r, Path{1})
	co.groupStart(a, Path{1, 2})
	co.groupStart(b, Path{1, 2, 3})
	co.groupStart(d, Path{1, 2, 3, 5})
	co.groupStart(r, Path{1})
	co.groupStart(z, Path{1, 6})

	exp := TestGroups{
		&TestGroup{
			Description: "a",
			Children: TestGroups{
				&TestGroup{
					Description: "b",
					Children: TestGroups{
						&TestGroup{
							Description: "c",
							Error:       c.Error,
						},
						&TestGroup{
							Description: "d",
						},
					},
				},
			},
		},
		&TestGroup{
			Description: "z",
		},
	}
	expect("the root group", co.group).NotEqual(nil)
	expect(co.group.Children).Equal(exp) //, "TreeCollector fail to reconstruct correct tree"
}

func TestSortingTestGroups(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	g0 := &TestGroup{ID: "0"}
	g1 := &TestGroup{ID: "1"}
	g2 := &TestGroup{ID: "2"}
	group := &TestGroup{Children: TestGroups{g2, g0, g1}}
	sortTestGroup(group)
	expect(group.Children).Equal(TestGroups{g0, g1, g2})
}
