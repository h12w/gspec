// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"

	ext "github.com/hailiang/gspec/extension"
)

// A Scheduler schedules test running.
type Scheduler struct {
	broadcaster
	*listener
	wg sync.WaitGroup
}

// NewScheduler creates and intialize a new Scheduler using r as the test
// reporter.
func NewScheduler(t ext.T, reporters ...ext.Reporter) *Scheduler {
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
	if !sequential {
		s.t.Parallel() // signal "go test" to allow concurrent testing.
	}
	s.Reporter.Start()
	for _, f := range funcs {
		(&runner{f, &s.wg, s.newSpec}).run(sequential)
	}
}

func (s *Scheduler) newSpec(g *group) S {
	return newSpec(g, s.listener)
}

// broadcaster syncs, receives and broadcasts all messages via Reporter interface
type broadcaster struct {
	t  ext.T
	a  []ext.Reporter
	mu sync.Mutex
}

func newBroadcaster(t ext.T, reporters []ext.Reporter) broadcaster {
	return broadcaster{t: t, a: reporters}
}

func (b *broadcaster) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.Start()
	}
}

func (b *broadcaster) End(groups ext.TestGroups) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.End(groups)
	}
}

func (b *broadcaster) Progress(g *ext.TestGroup, stats *ext.Stats) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if g.Error != nil {
		b.t.Fail()
	}
	for _, r := range b.a {
		r.Progress(g, stats)
	}
}
