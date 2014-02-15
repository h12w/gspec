package gspec

import "sync"

type runner struct {
	f       TestFunc
	wg      *sync.WaitGroup
	newSpec func(*group) S
}

func (r *runner) run(sequential bool) {
	if sequential {
		r.runSeq(path{})
	} else {
		r.runCon(path{})
	}
}

func (r *runner) runCon(dst path) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.runSpec(dst, r.runCon)
	}()
}

func (r *runner) runSeq(dst path) {
	r.runSpec(dst, r.runSeq)
}

func (r *runner) runSpec(dst path, run runFunc) {
	r.f(r.newSpec(newGrouper(dst, run)))
}

type group struct {
	dst    path
	cur    idStack
	done   bool
	runNew runFunc
}
type runFunc func(path)

func newGrouper(dst path, run runFunc) *group {
	return &group{dst: dst, runNew: run}
}

func (t *group) run(id funcID, f func()) {
	t.cur.push(id)
	defer t.cur.pop()
	if !t.cur.onPath(t.dst) {
		return
	} else if t.done {
		t.runNew(t.cur.clone())
		return
	}
	f()
	t.done = true
}

func (t *group) current() path {
	return t.cur.clone()
}

type idStack struct {
	path
}

func (p *idStack) push(i funcID) {
	p.path = append(p.path, i)
}
func (p *idStack) pop() (i funcID) {
	if len(p.path) == 0 {
		panic("call pop when idStack is empty.")
	}
	p.path, i = p.path[:len(p.path)-1], p.path[len(p.path)-1]
	return
}

type path []funcID

func (p path) clone() path {
	return append(path{}, p...)
}
func (p path) onPath(dst path) bool {
	// func id is unique, comparing the last should be enough
	if last := imin(len(p), len(dst)) - 1; last >= 0 {
		return p[last] == dst[last]
	}
	return true // initial idStack is empty
}
