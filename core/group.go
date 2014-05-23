// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// group implements the core algorithm of nested test groups.
type group struct {
	dst    Path
	cur    serialStack
	next   Serial
	done   bool
	runNew runFunc
}
type runFunc func(Path)

func newGroup(dst Path, run runFunc) *group {
	return &group{dst: dst, runNew: run}
}

// visit method calls each test group closure along the path defined by
// group.dst, ignoring the node not on path. Once a leaf node is visited, it
// stops calling the rest closures on path, but provides the path of them
// through group.runNew.
func (g *group) visit(f func(cur Path)) {
	g.cur.push(g.next)
	g.next = 0
	defer func() { g.next = g.cur.pop() + 1 }()
	if !g.cur.onPath(g.dst) {
		return
	} else if g.done {
		g.runNew(g.cur.Path)
		return
	}
	defer func() { g.done = true }()
	f(g.cur.Path)
}
