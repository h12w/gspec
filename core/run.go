// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"strings"
	"sync"
)

type runner struct {
	f       TestFunc
	wg      *sync.WaitGroup
	newSpec func(*group) S
}

func (r *runner) run(sequential bool, dst path) {
	if sequential {
		r.runSeq(dst)
	} else {
		r.runCon(dst)
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
	r.f(r.newSpec(newGroup(dst, run)))
}

type group struct {
	dst    path
	cur    idStack
	next   funcID
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
	defer func() {
		t.next = t.cur.pop() + 1
	}()
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

type path []funcID

func (p path) clone() path {
	return append(path{}, p...)
}

func (p path) onPath(dst path) bool {
	last := imin(len(p), len(dst)) - 1
	for i := 0; i <= last; i++ {
		if p[i] != dst[i] {
			return false
		}
	}
	return true // initial idStack is empty
}

func (p path) String() string {
	ss := make([]string, len(p))
	for i := range ss {
		ss[i] = p[i].String()
	}
	return strings.Join(ss, "/")
}

func (p *path) Set(s string) (err error) {
	ss := strings.Split(s, "/")
	*p = make(path, len(ss))
	for i := range ss {
		(*p)[i], err = parseFuncID(ss[i])
		if err != nil {
			*p = nil
			return
		}
	}
	return
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
