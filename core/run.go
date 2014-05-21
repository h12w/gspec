// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
)

// runner is used to run nested groups concurrently or sequencially.
//
// A queue is used to store the targets. It is not mandatary because the closure
// for a test group can be executed directly rather than pushing into a queue,
// but a queue is needed to keep the execution order in sequencial mode in the
// order the nested groups are written, making the algorithm easier to understand.
type runner struct {
	f       TestFunc
	q       pathQueue
	newSpec func(*group) S
	wg      sync.WaitGroup
}

func newRunner(f TestFunc, newSpec func(*group) S) *runner {
	return &runner{f: f, newSpec: newSpec}
}

func (r *runner) run(sequential bool, dst path) {
	r.q.enqueue(dst)
	for r.q.count() > 0 {
		for r.q.count() > 0 {
			dst := r.q.dequeue()
			r.runOne(sequential, dst)
		}
		// if the queue is empty, wait until all the current jobs are finished
		// and check again.
		r.wg.Wait()
	}
}

func (r *runner) runOne(sequential bool, dst path) {
	s := r.newSpec(
		newGroup(
			dst,
			func(newDst path) {
				r.q.enqueue(newDst)
			}))
	if sequential {
		r.f(s)
	} else {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			r.f(s)
		}()
	}
}
