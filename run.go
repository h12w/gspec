// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

type RootFunc func(g *G)

func Run(f RootFunc) {
	runner := concurrentRunner{baseRunner: baseRunner{f: f, collector: NewCollector()}}
	runner.run(path{})
	runner.Wait()
}

func RunSeq(f RootFunc) {
	runner := sequentialRunner{baseRunner{f: f, collector: NewCollector()}}
	runner.run(path{})
}

type concurrentRunner struct {
	baseRunner
	wg sync.WaitGroup
}

func (r *concurrentRunner) run(p path) {
	r.wg.Add(1) // no need to lock
	go func() {
		defer r.wg.Done()
		r.runPath(r, p)
	}()
}

func (r *concurrentRunner) Wait() {
	r.wg.Wait()
}

type sequentialRunner struct {
	baseRunner
}

func (r *sequentialRunner) run(p path) {
	r.runPath(r, p)
}

type baseRunner struct {
	f RootFunc
	collector
}

func (r *baseRunner) runPath(pr pathRunner, p path) {
	r.f(newG(p, pr, r.collector))
}
