// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

type RootFunc func(g G)

type Scheduler struct {
	wg sync.WaitGroup
	*treeListener
}

func NewScheduler(r Reporter) *Scheduler {
	return &Scheduler{treeListener: newTreeListener(r)}
}

func (r *Scheduler) Start(sequential bool, fs ...RootFunc) {
	defer func() {
		r.wg.Wait()
		r.Reporter.End(r.groups)
	}()
	r.Reporter.Start()
	for _, f := range fs {
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
