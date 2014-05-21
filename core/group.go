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

func (t *group) visit(f func()) {
	t.cur.push(t.next)
	t.next = 0
	defer func() { t.next = t.cur.pop() + 1 }()
	if !t.cur.onPath(t.dst) {
		return
	} else if t.done {
		t.runNew(t.cur.clone())
		return
	}
	defer func() { t.done = true }()
	f()
}

func (t *group) current() path {
	return t.cur.clone()
}
