// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	. "github.com/onsi/gomega"
	"runtime"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

/*
TODO:
* TextFormatter
* handle failure
* assert
*/

/*
Story: Dveloper write tests

As a developer
I want to write tests
So that I can get my test code run (sequentially or concurrently)
*/

/*
Scenario: run a test defined in a closure
	Given a test defined in a closure
	When it is executed
	Then it should be executed once and only once
*/
func TestRunClosureTest(t *testing.T) {
	ch := NewSChan()
	RunSeq(func(g *G) {
		do := g.Group
		do(func() {
			ch.Send("a")
		})
	})
	if exp := []string{"a"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
}

/*
Scenario: setup a common test context for two tests (before each)
	Given 3 tests (a, b, c) with a common setup code (s)
	When executed
	Then the ordered execution sequence is: sa sb sc
*/
func TestBeforeEach(t *testing.T) {
	ch := NewSChan()
	Run(func(g *G) {
		do := g.Group
		do(func() {
			s := "s"
			do(func() {
				s += "a"
				ch.Send(s)
			})
			do(func() {
				s += "b"
				ch.Send(s)
			})
			do(func() {
				s += "c"
				ch.Send(s)
			})
		})
	})
	if exp := []string{"sa", "sb", "sc"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

/*
Scenario: teardown a common test context for two tests (after each)
	Given two tests (a, b) with a common teardown code (t)
	When executed
	Then the ordered execution sequence is: at bt
*/
func TestAfterEach(t *testing.T) {
	ch := NewSChan()
	Run(func(g *G) {
		do := g.Group
		do(func() {
			s := ""
			defer func() {
				s += "t"
				ch.Send(s)
			}()
			do(func() {
				s += "a"
			})
			do(func() {
				s += "b"
			})
		})
	})
	if exp := []string{"at", "bt"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

/*
Scenario: nested testing group
	Given a nested testing group defined by closures like pseudo code below:

		a {
		    b {}
		    c {
		        d {}
			}
		}
		e {
		    f {}
		}

	When the test is run.
	Then the ordered execution sequence is: ab acd aef
*/
func TestNestedTestingContext(t *testing.T) {
	ch := NewSChan()
	Run(func(g *G) {
		do := g.Group
		do(func() {
			s := ""
			defer func() {
				ch.Send(s)
			}()
			s = "a"
			defer func() {
				s += "A"
			}()
			do(func() {
				s += "b"
			})
			do(func() {
				s += "c"
				defer func() {
					s += "C"
				}()
				do(func() {
					s += "d"
				})
			})
			do(func() {
				s += "e"
				do(func() {
					s += "f"
				})
			})
		})
	})
	if exp := []string{"abA", "acdCA", "aefA"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

/*
Scenario: concurrent running tests
	Given 5 identical time consuming test cases
	When they are completed
	Then the time to run all should be closer to 2 times of one test rather than 5 times
		(2 times because the first test must return first then the rest tests can be discovered and run simultaneously)
*/
func TestConcurrentRunning(t *testing.T) {
	delay := 10 * time.Millisecond
	tm := time.Now()
	Run(func(g *G) {
		do := g.Group
		do(func() {
			time.Sleep(delay)
			do(func() {
			})
			do(func() {
			})
			do(func() {
			})
			do(func() {
			})
			do(func() {
			})
		})
	})
	d := time.Now().Sub(tm)
	if d > 3*delay {
		t.Fatalf("Tests are not run concurrently, duration: %v", d)
	}
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
func TestDescriptions(t *testing.T) {
	RegisterTestingT(t)
	ds := []string{}
	NewCollector = func() collector {
		return CollectFunc(func(g *TestGroup, path []FuncId) {
			ds = append(ds, g.Description)
		})
	}
	Run(func(g *G) {
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

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestTreeCollector(t *testing.T) {
	RegisterTestingT(t)
	co := &treeCollector{}
	co.Start()
	a := &TestGroup{
		Id:          1,
		Description: "a",
	}
	b := &TestGroup{
		Id:          2,
		Description: "b",
	}
	c := &TestGroup{
		Id:          3,
		Description: "c",
	}
	d := &TestGroup{
		Id:          4,
		Description: "d",
	}
	z := &TestGroup{
		Id:          5,
		Description: "z",
	}
	co.GroupStart(a, []FuncId{})
	co.GroupStart(b, []FuncId{1})
	co.GroupStart(c, []FuncId{1, 2})
	err := &TestError{}
	co.GroupEnd(3, err)
	co.GroupStart(a, []FuncId{})
	co.GroupStart(b, []FuncId{1})
	co.GroupStart(d, []FuncId{1, 2})
	co.GroupStart(z, []FuncId{})
	co.End()

	exp := []*TestGroup{
		&TestGroup{
			Id:          1,
			Description: "a",
			Children: []*TestGroup{
				&TestGroup{
					Id:          2,
					Description: "b",
					Children: []*TestGroup{
						&TestGroup{
							Id:          3,
							Description: "c",
							Error:       err,
						},
						&TestGroup{
							Id:          4,
							Description: "d",
						},
					},
				},
			},
		},
		&TestGroup{
			Id:          5,
			Description: "z",
		},
	}
	Expect(co.Groups).To(Equal(exp), "TreeCollector fail to reconstruct correct tree")
}

func TestFuncUniqueId(t *testing.T) {
	f1 := func() {}
	f2 := func() {}
	if getFuncId(f1) != getFuncId(f1) {
		t.Fatalf("Does not return the same id for the same function.")
	}
	if getFuncId(f1) == getFuncId(f2) {
		t.Fatalf("Return the same id for different functions.")
	}
}

func TestPath(t *testing.T) {
	RegisterTestingT(t)
	p := path{}
	p.push(1)
	p.push(2)
	Expect(p.a).To(Equal([]FuncId{1, 2}), "path.push failed")
	i := p.pop()
	Expect(p.a).To(Equal([]FuncId{1}), "path.pop failed")
	Expect(i).To(Equal(FuncId(2)), "path.pop failed")
	i = p.pop()
	Expect(p.a).To(Equal([]FuncId{}), "path.pop failed")
	Expect(func() { p.pop() }).To(Panic(), "path.pop should panic when empty")
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}
