// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

/*
TODO:
	generate start/end event
	default listener
	register listener
*/

type RootFunc func(g *G)

func Run(f RootFunc) {
	newConcurrentRunner(f, nil).start()
}

func RunSeq(f RootFunc) {
	newSequentialRunner(f, nil).start()
}

type concurrentRunner struct {
	*sequentialRunner
	wg sync.WaitGroup
}

func newConcurrentRunner(f RootFunc, l Listener) *concurrentRunner {
	r := &concurrentRunner{sequentialRunner: newSequentialRunner(f, l)}
	r.self = r
	return r
}

func (r *concurrentRunner) start() {
	defer r.wg.Wait()
	r.sequentialRunner.start()
}

func (r *concurrentRunner) run(p path) {
	r.wg.Add(1) // no need to lock
	go func() {
		defer r.wg.Done()
		r.sequentialRunner.run(p)
	}()
}

type sequentialRunner struct {
	f    RootFunc
	l    Listener
	tc   treeCollector
	self runner
	groupListeners
}

func newSequentialRunner(f RootFunc, l Listener) *sequentialRunner {
	r := &sequentialRunner{f: f, l: l, tc: newTreeCollector()}
	r.groupListeners.add(l)
	r.self = r
	return r
}

func (r *sequentialRunner) start() {
	r.run(path{})
}

func (r *sequentialRunner) run(p path) {
	r.f(newG(p, r.self))
}
