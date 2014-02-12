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
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

/*
TODO:
* report failure location
*/

/*
Story: A dveloper runs tests

As a developer
I want to write tests
So that I can get my test code run (sequentially or concurrently)
*/

func aliasDo(s S) func(func()) {
	return func(f func()) { s.Alias("")("", f) }
}

/*
Scenario: run a test defined in a closure
	Given a test defined in a closure
	When it is executed
	Then it should be executed once and only once
*/
func TestRunClosureTest(t *testing.T) {
	ch := NewSChan()
	RunSeq(func(s S) {
		do := aliasDo(s)
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
	Run(func(s S) {
		do := aliasDo(s)
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
	Run(func(s S) {
		do := aliasDo(s)
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
	Run(func(s S) {
		do := aliasDo(s)
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
	Given N (N > 2) identical time consuming test cases
	When they are completed
	Then the time to run all should be 2 to 3 times of one test rather than N times
		(the first test must return first then the rest tests can be discovered and run simultaneously)
*/
func TestConcurrentRunning(t *testing.T) {
	delay := 10 * time.Millisecond
	tm := time.Now()
	Run(func(s S) {
		do := aliasDo(s)
		do(func() {
			time.Sleep(delay)
			do(func() {
			})
			do(func() {
			})
			do(func() {
			})
		})
		do(func() {
			time.Sleep(delay)
		})
		do(func() {
			time.Sleep(delay)
		})
	})
	d := time.Now().Sub(tm)
	if d > time.Duration(2.3*float64(delay)) {
		t.Fatalf("Tests are not run concurrently, duration: %v", d)
	}
}

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
		return CollectFunc(func(g *TestGroup, path []FuncID) {
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
	co := newTreeListener(NewTextReporter(os.Stdout))
	a := &TestGroup{
		ID:          FuncID{1},
		Description: "a",
	}
	b := &TestGroup{
		ID:          FuncID{2},
		Description: "b",
	}
	c := &TestGroup{
		ID:          FuncID{3},
		Description: "c",
	}
	cp := []FuncID{{1}, {2}}
	d := &TestGroup{
		ID:          FuncID{4},
		Description: "d",
	}
	z := &TestGroup{
		ID:          FuncID{5},
		Description: "z",
	}
	co.groupStart(a, []FuncID{})
	co.groupStart(b, []FuncID{{1}})
	co.groupStart(c, cp)
	c.Error = &TestError{}
	co.groupStart(a, []FuncID{})
	co.groupStart(b, []FuncID{{1}})
	co.groupStart(d, []FuncID{{1}, {2}})
	co.groupStart(z, []FuncID{})

	exp := []*TestGroup{
		&TestGroup{
			ID:          FuncID{1},
			Description: "a",
			Children: []*TestGroup{
				&TestGroup{
					ID:          FuncID{2},
					Description: "b",
					Children: []*TestGroup{
						&TestGroup{
							ID:          FuncID{3},
							Description: "c",
							Error:       c.Error,
						},
						&TestGroup{
							ID:          FuncID{4},
							Description: "d",
						},
					},
				},
			},
		},
		&TestGroup{
			ID:          FuncID{5},
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
	p := path{}
	p.push(FuncID{1})
	p.push(FuncID{2})
	expect(p.a).Equal([]FuncID{{1}, {2}}) //, "path.push failed")
	i := p.pop()
	expect(p.a).Equal([]FuncID{{1}}) //, "path.pop failed")
	expect(i).Equal(FuncID{2})       //, "path.pop failed")
	i = p.pop()
	expect(p.a).Equal([]FuncID{})      //, "path.pop failed")
	expect(func() { p.pop() }).Panic() //, "path.pop should panic when empty")
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}
