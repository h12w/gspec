// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
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
	ch := NewSS()
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
	Then the ordered execution sequence is: sa, sb, sc in any order.
*/
func TestBeforeEach(t *testing.T) {
	ch := NewSS()
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
	Then the ordered execution sequence is: at, bt in any order.
*/
func TestAfterEach(t *testing.T) {
	ch := NewSS()
	runGroup(func(g groupFunc) {
		s := ""
		g(func() {
			defer func() {
				s += "t"
			}()
			g(func() {
				s += "a"
			})
			g(func() {
				s += "b"
			})
		})
		ch.Send(s)
	})
	if exp := []string{"at", "bt"}; !ch.Equal(exp) {
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
	ch := NewSS()
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
	if exp := []string{"a", "d", "b", "e", "c", "f"}; !ch.Equal(exp) {
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
	Then the ordered execution sequence is: ab, acd, aef in any order.
*/
func TestNestedTestingContext(t *testing.T) {
	ch := NewSS()
	runGroup(func(g groupFunc) {
		s := ""
		g(func() {
			s = "a"
			defer func() {
				s += "A"
				ch.Send(s)
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
	if exp := []string{"abA", "acdCA", "aefA"}; !ch.Equal(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

// runGroup is a simplified implementation to run nested test group.
func runGroup(f func(g groupFunc)) {
	fifo := &pathQueue{}
	fifo.enqueue(Path{})
	for fifo.count() > 0 {
		dst := fifo.dequeue()
		runPath(dst, f, fifo)
	}
}

type groupFunc func(func())

func runPath(dst Path, f func(g groupFunc), fifo *pathQueue) {
	group := newGroup(
		dst,
		func(newDst Path) {
			fifo.enqueue(newDst)
		})
	f(func(ff func()) {
		group.visit(func(cur Path) {
			ff()
		})
	})
}
