package gspec

// The interface for G to call back
type runner interface {
	run(p path)
	GroupListener
}

// G contains minimal context variables needed to implement nested test group
type G struct {
	dst       path
	cur       path
	skipRest  bool
	skipCount int
	runner
}

func newG(p path, r runner) *G {
	return &G{dst: p, runner: r}
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
