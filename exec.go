// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

type RootFunc func(g *G)

func Run(f RootFunc) {
	runner := &concurrentRunner{f: f, collector: NewCollector()}
	runner.run(path{})
	runner.Wait()
}

func RunSeq(f RootFunc) {
	runner := &sequentialRunner{f: f, collector: NewCollector()}
	runner.run(path{})
	runner.Wait()
}

type runner interface {
	run(p path)
}

type concurrentRunner struct {
	f         RootFunc
	wg        sync.WaitGroup
	collector Collector
}

func (r *concurrentRunner) run(p path) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.f(newG(p, r, r.collector)) // Always use a new G
	}()
}

func (r *concurrentRunner) Wait() {
	r.wg.Wait()
}

type sequentialRunner struct {
	f         RootFunc
	collector Collector
}

func (r *sequentialRunner) run(p path) {
	r.f(newG(p, r, r.collector))
}

func (r *sequentialRunner) Wait() {
}

// G contains minimal context variables needed to implement nested testing
// group.
type G struct {
	dst       path
	cur       path
	skipRest  bool
	skipCount int
	collector Collector
	runner
}

func newG(p path, r runner, c Collector) *G {
	return &G{dst: p, runner: r, collector: c}
}

func (t *G) Group(f func()) bool {
	t.cur.push(getFuncId(f))
	defer t.cur.pop()
	if !t.cur.onPath(t.dst) {
		return false
	} else if t.skipRest {
		t.run(t.cur.clone())
		t.skipCount++
		return false
	}
	sc := t.skipCount
	f()
	if sc == t.skipCount { // true when f is a leaf node
		t.skipRest = true
	}
	return true
}

type path struct {
	a []FuncId
}

func (p *path) push(i FuncId) {
	p.a = append(p.a, i)
}

func (p *path) pop() (i FuncId) {
	if len(p.a) == 0 {
		panic("call pop when path is empty.")
	}
	p.a, i = p.a[:len(p.a)-1], p.a[len(p.a)-1]
	return
}

func (p *path) clone() path {
	return path{append([]FuncId{}, p.a...)}
}

func (p *path) onPath(dst path) bool {
	// func id is unique, comparing the last should be enough
	if last := imin(len(p.a), len(dst.a)) - 1; last >= 0 {
		return p.a[last] == dst.a[last]
	}
	return true // initial path is empty
}
