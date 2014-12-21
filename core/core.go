// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package core provides a minimal core for organizing, executing and reporting of
test cases nested in test groups.
*/
package core // import "h12.me/gspec/core"

import (
	"sync"

	ext "h12.me/gspec/extension"
)

// Controller is the "C" of MVC (Model View Controller). In a test framework,
// all the test cases form a model (but unchangable by the controller), the test
// reporter is the view, and the controller controls the test running and send
// the test result to the test reporter.
type Controller struct {
	*collector
	*broadcaster
}

// NewController creates and intialize a new Controller using r as the test
// reporter.
func NewController(reporters ...ext.Reporter) *Controller {
	c := &Controller{
		broadcaster: newBroadcaster(reporters),
	}
	c.collector = newCollector(c.broadcaster)
	return c
}

// Start starts tests defined in funcs concurrently or sequentially.
func (c *Controller) Start(path Path, concurrent bool, funcs ...TestFunc) error {
	c.broadcaster.Start()
	defer func() {
		c.collector.sort()
		c.broadcaster.End(c.group)
	}()

	newRunner(func(s S) {
		top := s.Alias("")
		top("", func() {
			for _, f := range funcs {
				f(s)
			}
		})
	}, concurrent, c.collector).run(path)

	return nil
}

// broadcaster syncs, receives and broadcasts all messages via Reporter interface
type broadcaster struct {
	a  []ext.Reporter
	mu sync.Mutex
}

func newBroadcaster(reporters []ext.Reporter) *broadcaster {
	return &broadcaster{a: reporters}
}

func (b *broadcaster) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.Start()
	}
}

func (b *broadcaster) End(group *ext.TestGroup) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.End(group)
	}
}

func (b *broadcaster) Progress(g *ext.TestGroup, stats *ext.Stats) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.Progress(g, stats)
	}
}
