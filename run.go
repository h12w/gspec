// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

// RootFunc is the type of the function called for each test case.
type RootFunc func(g G)

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
func (r *Scheduler) Start(sequential bool, funcs ...RootFunc) {
	defer func() {
		r.wg.Wait()
		r.Reporter.End(r.groups)
	}()
	r.Reporter.Start()
	for _, f := range funcs {
		if sequential {
			seq{r}.run(f, path{})
		} else {
			con{r}.run(f, path{})
		}
	}
}

func (r *Scheduler) runSeq(f RootFunc, p path, self scheduler) {
	f(newG(f, p, self))
}

func (r *Scheduler) runCon(f RootFunc, p path, self scheduler) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runSeq(f, p, self)
	}()
}

type con struct {
	*Scheduler
}

func (r con) run(f RootFunc, p path) {
	r.runCon(f, p, r)
}

type seq struct {
	*Scheduler
}

func (r seq) run(f RootFunc, p path) {
	r.runSeq(f, p, r)
}
