// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

type group struct {
	dst    path
	cur    serialStack
	next   serial
	done   bool
	runNew runFunc
}
type runFunc func(path)

func newGroup(dst path, run runFunc) *group {
	return &group{dst: dst, runNew: run}
}

func (g *group) visit(f func()) {
	g.cur.push(g.next)
	g.next = 0
	defer func() { g.next = g.cur.pop() + 1 }()
	if !g.cur.onPath(g.dst) {
		return
	} else if g.done {
		g.runNew(g.cur.path)
		return
	}
	defer func() { g.done = true }()
	f()
}

func (g *group) current() path {
	return g.cur.clone()
}
