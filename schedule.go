// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

// A Scheduler schedules test running.
type Scheduler struct {
	wg sync.WaitGroup
	*treeListener
}

// NewScheduler creates and intialize a new Scheduler using r as the test
// reporter.
func NewScheduler(r Reporter) *Scheduler {
	return &Scheduler{treeListener: newTreeListener(r)}
}

// Start starts tests defined in funcs concurrently or sequentially.
func (s *Scheduler) Start(sequential bool, funcs ...TestFunc) {
	defer func() {
		s.wg.Wait()
		s.Reporter.End(s.groups)
	}()
	s.Reporter.Start()
	for _, f := range funcs {
		r := runner{f, &s.wg, s.newS}
		if sequential {
			r.runSeq(path{})
		} else {
			r.runCon(path{})
		}
	}
}

func (s *Scheduler) newS(g grouper) S {
	return newS(g, s)
}
