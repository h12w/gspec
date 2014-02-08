package gspec

// The interface for G to call back
type scheduler interface {
	run(f RootFunc, p path)
	listener
}

// G contains minimal context needed to implement nested test group
type G struct {
	f         RootFunc
	dst       path
	cur       path
	skipRest  bool
	skipCount int
	scheduler
}

func newG(f RootFunc, p path, s scheduler) *G {
	return &G{f: f, dst: p, scheduler: s}
}

func (t *G) group(id FuncId, f func()) {
	t.cur.push(id)
	defer t.cur.pop()
	if !t.cur.onPath(t.dst) {
		return
	} else if t.skipRest {
		t.run(t.f, t.cur.clone())
		t.skipCount++
		return
	}
	sc := t.skipCount
	f()
	if sc == t.skipCount { // true when f is a leaf node
		t.skipRest = true
	}
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

func (p *path) slice() []FuncId {
	return append([]FuncId{}, p.a...)
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
