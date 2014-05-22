// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"time"
)

// TestFunc is the type of the function prepared to run in a goroutine for each
// test case.
type TestFunc func(S)

// toConcurrent converts a TestFunc to its concurrent version.
func (f TestFunc) toConcurrent(wg *sync.WaitGroup) TestFunc {
	return func(s S) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f(s)
		}()
	}
}

// measureTime decorates a TestFunc with time measurement.
func (f TestFunc) measureTime() TestFunc {
	return func(s S) {
		t := time.Now()
		f(s)
		d := time.Now().Sub(t)
		s.setDuration(d)
	}
}

// runner is used to run nest test groups in a TestFunc concurrently or
// sequencially.
//
// A queue is used to store the targets. It is not mandatary because the closure
// for a test group can be executed directly rather than pushing into a queue,
// but a queue is needed to keep the execution order in sequencial mode in the
// order the nested groups are written, making the algorithm easier to understand.
type runner struct {
	f       TestFunc
	q       pathQueue
	wg      sync.WaitGroup
	newSpec func(*group) S
}

func newRunner(f TestFunc, concurrent bool, newSpec func(*group) S) *runner {
	r := &runner{f: f.measureTime(), newSpec: newSpec}
	if concurrent {
		r.f = r.f.toConcurrent(&r.wg)
	}
	return r
}

func (r *runner) run(dst path) {
	r.q.enqueue(dst)
	for r.q.count() > 0 {
		for r.q.count() > 0 {
			dst := r.q.dequeue()
			r.runOne(dst)
		}
		// make sure there are no running test groups so that all test
		// groups have been visited.
		r.wg.Wait()
	}
}

func (r *runner) runOne(dst path) {
	r.f(r.newSpec(
		newGroup(
			dst,
			func(newDst path) {
				r.q.enqueue(newDst)
			})))
}
