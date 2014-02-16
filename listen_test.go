package gspec

import (
	exp "github.com/hailiang/gspec/expectation"
	"os"
	"testing"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestTreeListener(t *testing.T) {
	expect := exp.AliasForT(t)
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
	cp := path{{1}, {2}, {3}}
	d := &TestGroup{
		Description: "d",
	}
	z := &TestGroup{
		Description: "z",
	}
	co.groupStart(a, path{{1}})
	co.groupStart(b, path{{1}, {2}})
	co.groupStart(c, cp)
	c.Error = &TestError{}
	co.groupStart(a, path{{1}})
	co.groupStart(b, path{{1}, {2}})
	co.groupStart(d, path{{1}, {2}, {4}})
	co.groupStart(z, path{{5}})

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
