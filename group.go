package gspec

// G holds an internal object that contains minimal context needed to implement
// nested test group.
type G interface {
	Alias(name string) DescFunc
}

// The interface for groupContext to call back.
type scheduler interface {
	run(f RootFunc, p path)
	listener
}

type groupContext struct {
	f         RootFunc
	dst       path
	cur       path
	skipRest  bool
	skipCount int
	scheduler
}

func newG(f RootFunc, p path, s scheduler) *groupContext {
	return &groupContext{f: f, dst: p, scheduler: s}
}

func (t *groupContext) group(id FuncID, f func()) {
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
