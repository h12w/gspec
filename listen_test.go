// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"errors"
	"os"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: Internal Tests
	Test internal types/functions
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
	cp := path{{p: 1}, {p: 2}, {p: 3}}
	d := &TestGroup{
		Description: "d",
	}
	z := &TestGroup{
		Description: "z",
	}
	co.groupStart(a, path{{p: 1}})
	co.groupStart(b, path{{p: 1}, {p: 2}})
	co.groupStart(c, cp)
	c.Error = errors.New("c err")
	co.groupStart(a, path{{p: 1}})
	co.groupStart(b, path{{p: 1}, {p: 2}})
	co.groupStart(d, path{{p: 1}, {p: 2}, {p: 4}})
	co.groupStart(z, path{{p: 5}})

	exp := []*TestGroup{
		&TestGroup{
			Description: "a",
			Children: []*TestGroup{
				&TestGroup{
					Description: "b",
					Children: []*TestGroup{
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
