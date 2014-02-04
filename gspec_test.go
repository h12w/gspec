// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

/*
TODO:
* describe DSL
* handle panic
* assert
* formatter
* go test options (e.g. parallel)
*/

/*
Story: Dveloper write tests

As a developer
I want to write tests
So that I can get my code verified by running those tests
*/

/*
Scenario: run a test defined in a closure
	Given a test defined in a closure
	When it is executed
	Then it should be executed once and only once
*/
func TestRunClosureTest(t *testing.T) {
	ch := NewSChan()
	Run(func(do Desc) {
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
	Run(func(do Desc) {
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
		t.Fatalf("Wrong execution sequence for nested context, expected: %v, got: %v", exp, ch.Slice())
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
	Run(func(do Desc) {
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
		t.Fatalf("Wrong execution sequence for nested context, expected: %v, got: %v", exp, ch.Slice())
	}
}

/*
Scenario: nested testing context
	Given a nested testing context defined by closures like pseudo code below:

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
	Run(func(do Desc) {
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
		t.Fatalf("Wrong execution sequence for nested context, expected: %v, got: %v", exp, ch.Slice())
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
	Run(func(do Desc) {
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
	p := path{}
	p.push(1)
	p.push(2)
	if exp := []funcId{1, 2}; !idSliceEqual(p.a, exp) {
		t.Fatalf("path.push failed, expected: %v, got %v", exp, p.a)
	}
	i := p.pop()
	if exp := []funcId{1}; !idSliceEqual(p.a, exp) {
		t.Fatalf("path.pop failed, expected: %v, got %v", exp, p.a)
	}
	if i != 2 {
		t.Fatalf("path.pop failed, expected: %v, got %v", 2, i)
	}
	i = p.pop()
	if exp := []funcId{}; !idSliceEqual(p.a, exp) {
		t.Fatalf("path.pop failed, expected: %v, got %v", exp, p.a)
	}

	if !panicked(func() { p.pop() }) {
		t.Fatal("path.pop should panic when empty")
	}
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}

func panicked(f func()) (r bool) {
	defer func() {
		if err := recover(); err != nil {
			r = true
		}
	}()
	f()
	return false
}

type SChan struct {
	ch chan string
	ss []string
	wg sync.WaitGroup
}

func NewSChan() *SChan {
	return &SChan{ch: make(chan string)}
}

func (c *SChan) Send(s string) {
	c.wg.Add(1)
	go func() {
		c.ch <- s
		c.wg.Done()
	}()
}

func (c *SChan) Slice() []string {
	return c.ss
}

func (c *SChan) receiveAll() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for s := range c.ch {
			c.ss = append(c.ss, s)
		}
		wg.Done()
	}()
	c.wg.Wait()
	close(c.ch)
	wg.Wait()
}

func (c *SChan) EqualSorted(ss []string) bool {
	c.receiveAll()
	sort.Strings(c.ss)
	return c.equal(ss)
}

func (c *SChan) equal(ss []string) bool {
	if len(ss) != len(c.ss) {
		return false
	}
	for i := range ss {
		if ss[i] != c.ss[i] {
			return false
		}
	}
	return true
}

func idSliceEqual(a, b []funcId) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
