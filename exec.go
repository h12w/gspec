// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

type Desc func(func())

func Run(f func(do Desc)) {
	var wg sync.WaitGroup
	newT(f, path{}, &wg).runPath(path{})
	wg.Wait()
}

type T struct {
	f func(Desc)
	path
	cur       path
	skipRest  bool
	skipCount int
	wg        *sync.WaitGroup
}

func newT(f func(Desc), p path, wg *sync.WaitGroup) *T {
	return &T{f: f, path: p, wg: wg}
}

func (t *T) do(f func()) {
	id := getFuncId(f)
	t.cur.push(id)
	defer t.cur.pop()
	if !t.onPath(t.cur) {
		return
	} else if t.skipRest {
		t.runPath(t.cur.clone())
		t.skipCount++
		return
	}
	sc := t.skipCount
	f()
	if sc == t.skipCount {
		t.skipRest = true
	}
}

func (t *T) runPath(p path) {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.f(newT(t.f, p, t.wg).do) // Always use a new T
	}()
}

type path struct {
	a []funcId
}

func (p *path) push(i funcId) {
	p.a = append(p.a, i)
}

func (p *path) pop() (i funcId) {
	if len(p.a) == 0 {
		panic("call pop when path is empty.")
	}
	p.a, i = p.a[:len(p.a)-1], p.a[len(p.a)-1]
	return
}

func (p *path) clone() path {
	return path{append([]funcId{}, p.a...)}
}

func (p *path) onPath(cur path) bool {
	// func id is unique, comparing the last should be enough
	if last := imin(len(cur.a), len(p.a)) - 1; last >= 0 {
		return cur.a[last] == p.a[last]
	}
	return true // initial path is empty
}
