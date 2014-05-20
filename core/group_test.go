// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: A dveloper runs tests

As a developer
I want to write tests
So that I can get my test code run
*/

/*
Scenario: run a test defined in a closure
	Given a test defined in a closure
	When it is executed
	Then it should be executed once and only once
*/
func TestRunClosureTest(t *testing.T) {
	ch := NewSChan()
	runGroup(func(g groupFunc) {
		g(func() {
			ch.Send("a")
		})
	})
	if exp := []string{"a"}; !ch.Equal(exp) {
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
	runGroup(func(g groupFunc) {
		g(func() {
			s := "s"
			g(func() {
				s += "a"
				ch.Send(s)
			})
			g(func() {
				s += "b"
				ch.Send(s)
			})
			g(func() {
				s += "c"
				ch.Send(s)
			})
		})
	})
	if exp := []string{"sa", "sb", "sc"}; !ch.Equal(exp) {
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
	runGroup(func(g groupFunc) {
		g(func() {
			s := ""
			defer func() {
				s += "t"
				ch.Send(s)
			}()
			g(func() {
				s += "a"
			})
			g(func() {
				s += "b"
			})
		})
	})
	if exp := []string{"at", "bt"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

/*
Scenario: Table driven test
	Given test cases defined in a for loop
	When executed
	Each test case get run once
*/
func TestTableDriven(t *testing.T) {
	ch := NewSChan()
	runGroup(func(g groupFunc) {
		g(func() { // outer
			for i := 0; i < 3; i++ {
				g(func() { // loop a,b,c
					g(func() { // inner
						s := string('a' + i)
						ch.Send(s)
					})
				})
				g(func() { // loop d,e,f
					g(func() { // inner
						s := string('d' + i)
						ch.Send(s)
					})
				})
			}
		})
	})
	if exp := []string{"a", "b", "c", "d", "e", "f"}; !ch.EqualSorted(exp) {
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
	runGroup(func(g groupFunc) {
		g(func() {
			s := ""
			defer func() {
				ch.Send(s)
			}()
			s = "a"
			defer func() {
				s += "A"
			}()
			g(func() {
				s += "b"
			})
			g(func() {
				s += "c"
				defer func() {
					s += "C"
				}()
				g(func() {
					s += "d"
				})
			})
			g(func() {
				s += "e"
				g(func() {
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
Story: Internal Tests
	Test internal types/functions
*/

func TestPathSerialization(t *testing.T) {
	expect := exp.Alias(exp.TFailNow(t))

	var p path
	p.Set("0/1/2")
	expect(len(p)).Equal(3)
	expect(p[0]).Equal(funcID(0))
	expect(p[1]).Equal(funcID(1))
	expect(p[2]).Equal(funcID(2))

	err := p.Set("UVW")
	expect(err).NotEqual(nil)

	p = path{0, 1, 2}
	expect(p.String()).Equal("0/1/2")
}

func TestIDStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	p := idStack{}
	p.push(funcID(1))
	p.push(funcID(2))
	expect(p.path).Equal(path{1, 2})
	i := p.pop()
	expect(p.path).Equal(path{1})
	expect(i).Equal(funcID(2))
	i = p.pop()
	expect(p.path).Equal(path{})
	expect(func() { p.pop() }).Panic()
}

type groupFunc func(func())

func runGroup(f func(g groupFunc)) {
	runPath(path{}, f)
}

func runPath(dst path, f func(g groupFunc)) {
	group := newGroup(
		dst,
		func(dst path) {
			runPath(dst, f)
		})
	f(group.visit)
}
