// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
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

// toMeasurable decorates a TestFunc with time measurement.
func (f TestFunc) toMeasurable() TestFunc {
	return func(s S) {
		s.start()
		defer s.end()
		f(s)
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
	f  TestFunc
	q  pathQueue
	wg sync.WaitGroup
	c  *collector
}

func newRunner(f TestFunc, concurrent bool, c *collector) *runner {
	r := &runner{f: f.toMeasurable(), c: c}
	if concurrent {
		r.f = r.f.toConcurrent(&r.wg)
	}
	return r
}

func (r *runner) run(dst Path) {
	r.q.enqueue(dst)
	for r.q.count() > 0 {
		for r.q.count() > 0 {
			dst := r.q.dequeue()
			r.runOne(dst)
		}
		// make sure that there are no test groups running so that all groups
		// have been visited.
		r.wg.Wait()
	}
}

func (r *runner) runOne(dst Path) {
	r.f(newSpec(
		newGroup(dst, func(newDst Path) { r.q.enqueue(newDst) }),
		r.c,
	))
}
