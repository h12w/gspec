// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"
	"time"
)

/*
Scenario: concurrent running tests
	Given N (N > 2) identical time consuming test cases
	When they are completed
	Then the time to run all should be 2 to 3 times of one test rather than N times
		(the first test must return first then the rest tests can be discovered and run simultaneously)
*/
func TestConcurrentRunningTime(t *testing.T) {
	delay := 10 * time.Millisecond
	tm := time.Now()
	runCon(func(s S) {
		g := aliasGroup(s)
		g(func() {
			time.Sleep(delay)
			g(func() {
			})
			g(func() {
			})
			g(func() {
			})
		})
		g(func() {
			time.Sleep(delay)
		})
		g(func() {
			time.Sleep(delay)
		})
	})
	d := time.Now().Sub(tm)
	if d > time.Duration(3*float64(delay)) {
		t.Fatalf("Tests are not run concurrently, duration: %v", d)
	}
}

func TestNestedTestingContextConcurrently(t *testing.T) {
	ch := NewSS()
	runCon(func(s S) {
		g := aliasGroup(s)
		g(func() {
			s := "a"
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
	if exp := []string{"abA", "acdCA", "aefA"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution sequence for nested group, expected: %v, got: %v", exp, ch.Slice())
	}
}

func TestRunSeq(t *testing.T) {
	ch := NewSS()
	runSeq(func(s S) {
		g := aliasGroup(s)
		g(func() {
			ch.Send("a")
		})
	})
	if exp := []string{"a"}; !ch.Equal(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
}

func aliasGroup(s S) func(func()) {
	return func(f func()) { s.Alias("")("", f) }
}
