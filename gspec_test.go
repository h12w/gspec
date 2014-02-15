// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"bytes"
	exp "github.com/hailiang/gspec/expectation"
	"os"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

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

/*
Story: Dveloper describe tests

As a developer
I want to describe my tests
So that a structured, readable specification can be generated
*/

/*
Scenario: Use customized alias names for group function.
	Given the G object provided by GSpec
	When I want to define a better name (or a set of different names) for it
	Then I can make aliases out of it and GSpec will pass them within description

Scenario: Attach a string as the description to the testing group
	Given a group function
	When attach a string as the description
	Then GSpec will pass them to collector
*/
/*
func TestDescriptions(t *testing.T) {
	RegisterTestingT(t)
	ds := []string{}
	NewCollector = func() collector {
		return CollectFunc(func(g *TestGroup, path path) {
			ds = append(ds, g.Description)
		})
	}
	Run(func(s S) {
		describe, context, it := g.Alias3("describe", "context", "it")
		describe, it = g.Alias2("describe", "it")
		describe("a", func() {
			context("b", func() {
				it("c", func() {
				})
			})
		})
	})
	Expect(ds).To(Equal([]string{"describe a", "context b", "it c"}), "description not stored correctly")
}
*/

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

func TestFuncUniqueID(t *testing.T) {
	f1 := func() {}
	f2 := func() {}
	if getFuncID(f1) != getFuncID(f1) {
		t.Fatalf("Does not return the same id for the same function.")
	}
	if getFuncID(f1) == getFuncID(f2) {
		t.Fatalf("Return the same id for different functions.")
	}
}

func TestPath(t *testing.T) {
	expect := exp.AliasForT(t)
	p := idStack{}
	p.push(funcID{1})
	p.push(funcID{2})
	expect(p.path).Equal(path{{1}, {2}}) //, "path.push failed")
	i := p.pop()
	expect(p.path).Equal(path{{1}}) //, "path.pop failed")
	expect(i).Equal(funcID{2})      //, "path.pop failed")
	i = p.pop()
	expect(p.path).Equal(path{})       //, "path.pop failed")
	expect(func() { p.pop() }).Panic() //, "path.pop should panic when empty")
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}
