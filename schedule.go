// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

// A Scheduler schedules test running.
type Scheduler struct {
	broadcaster
	*listener
	wg sync.WaitGroup
}

// NewScheduler creates and intialize a new Scheduler using r as the test
// reporter.
func NewScheduler(t T, reporters ...Reporter) *Scheduler {
	s := &Scheduler{broadcaster: newBroadcaster(t, reporters)}
	s.listener = newListener(&s.broadcaster)
	return s
}

// Start starts tests defined in funcs concurrently or sequentially.
func (s *Scheduler) Start(sequential bool, funcs ...TestFunc) {
	defer func() {
		s.wg.Wait()
		s.Reporter.End(s.groups)
	}()
	s.Reporter.Start()
	for _, f := range funcs {
		(&runner{f, &s.wg, s.newSpec}).run(sequential)
	}
}

func (s *Scheduler) newSpec(g *group) S {
	return newSpec(g, s.listener)
}
