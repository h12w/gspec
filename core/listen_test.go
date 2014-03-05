// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"os"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
	. "github.com/hailiang/gspec/extension"
	. "github.com/hailiang/gspec/reporter"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

/*
Scenario: reconstruct nested test group to a tree
	Given a treeListener
	When it's groupStart method is called like real tests running
	Then it is able to reconstruct the tree structure
*/
func TestTreeListener(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	co := newListener(NewTextReporter(os.Stdout))
	a := &TestGroup{
		Description: "a",
	}
	b := &TestGroup{
		Description: "b",
	}
	c := &TestGroup{
		Description: "c",
	}
	cp := path{1, 2, 3}
	d := &TestGroup{
		Description: "d",
	}
	z := &TestGroup{
		Description: "z",
	}
	co.groupStart(a, path{1})
	co.groupStart(b, path{1, 2})
	co.groupStart(c, cp)
	c.Error = errors.New("c err")
	co.groupStart(a, path{1})
	co.groupStart(b, path{1, 2})
	co.groupStart(d, path{1, 2, 4})
	co.groupStart(z, path{5})

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
	expect(co.groups).Equal(exp) //, "TreeListener fail to reconstruct correct tree"
}
