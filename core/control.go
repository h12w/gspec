// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"

	ext "github.com/hailiang/gspec/extension"
)

// T is an interface that allows a testing.T to be passed to GSpec.
type T interface {
	Fail()
	Parallel()
}

// Controller is the "C" of MVC (Model View Controller). In a test framework,
// all the test cases form a model (but unchangable by the controller), the test
// reporter is the view, and the controller controls the test running and send
// the test result to the test reporter.
type Controller struct {
	*broadcaster
	*listener
	*config
}

// NewController creates and intialize a new Controller using r as the test
// reporter.
func NewController(t T, reporters ...ext.Reporter) *Controller {
	s := &Controller{
		broadcaster: newBroadcaster(t, reporters),
		config:      &globalConfig,
	}
	s.listener = newListener(s.broadcaster)
	return s
}

// Start starts tests defined in funcs concurrently or sequentially.
func (s *Controller) Start(sequential bool, funcs ...TestFunc) error {
	if !sequential {
		s.t.Parallel() // signal "go test" to allow concurrent testing.
	}

	s.broadcaster.Start()
	defer func() {
		s.broadcaster.End(s.groups)
	}()

	newRunner(func(s S) {
		for _, f := range funcs {
			f(s)
		}
	}, s.newSpec).run(sequential, s.focus)

	return nil
}

func (s *Controller) newSpec(g *group) S {
	return newSpec(g, s.listener)
}

// broadcaster syncs, receives and broadcasts all messages via Reporter interface
type broadcaster struct {
	t  T
	a  []ext.Reporter
	mu sync.Mutex
}

func newBroadcaster(t T, reporters []ext.Reporter) *broadcaster {
	return &broadcaster{t: t, a: reporters}
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
