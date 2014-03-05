// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
	"time"

	exp "github.com/hailiang/gspec/expectation"
)

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
Scenario: Table driven test
	Given test cases defined in a for loop
	When executed
	Each test case get run once
*/
func TestTableDriven(t *testing.T) {
	ch := NewSChan()
	Run(func(s S) {
		do := s.Alias("")
		do("outer", func() {
			for i := 0; i < 3; i++ {
				do("loop a,b,c", func() {
					do("inner", func() {
						s := string('a' + i)
						ch.Send(s)
					})
				})
				do("loop d,e,f", func() {
					do("inner", func() {
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
