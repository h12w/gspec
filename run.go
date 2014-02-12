package gspec

import (
	"sync"
)

type grouper interface {
	group(id FuncID, f func())
	current() []FuncID
}

type runner struct {
	f    TestFunc
	wg   *sync.WaitGroup
	newS func(grouper) S
}

func (r runner) runCon(dst path) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.f(r.newS(newGrouper(dst, r.runCon)))
	}()
}

func (r runner) runSeq(dst path) {
	r.f(r.newS(newGrouper(dst, r.runSeq)))
}

type grouperImpl struct {
	dst       path
	cur       path
	skipRest  bool
	skipCount int
	run       runFunc
}
type runFunc func(path)

func newGrouper(dst path, run runFunc) grouper {
	return &grouperImpl{dst: dst, run: run}
}

func (t *grouperImpl) group(id FuncID, f func()) {
	t.cur.push(id)
	defer t.cur.pop()
	if !t.cur.onPath(t.dst) {
		return
	} else if t.skipRest {
		t.run(t.cur.clone())
		t.skipCount++
		return
	}
	sc := t.skipCount
	f()
	if sc == t.skipCount { // true when f is a leaf node
		t.skipRest = true
	}
}

func (t *grouperImpl) current() []FuncID {
	return t.cur.slice()
}

type path struct {
	a []FuncID
}

func (p *path) push(i FuncID) {
	p.a = append(p.a, i)
}

func (p *path) pop() (i FuncID) {
	if len(p.a) == 0 {
		panic("call pop when path is empty.")
	}
	p.a, i = p.a[:len(p.a)-1], p.a[len(p.a)-1]
	return
}

func (p *path) slice() []FuncID {
	return append([]FuncID{}, p.a...)
}

func (p *path) clone() path {
	return path{p.slice()}
}

func (p *path) onPath(dst path) bool {
	// func id is unique, comparing the last should be enough
	if last := imin(len(p.a), len(dst.a)) - 1; last >= 0 {
		return p.a[last] == dst.a[last]
	}
	return true // initial path is empty
}
