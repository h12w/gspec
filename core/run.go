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

func (t *group) visit(id funcID, f func()) {
	t.cur.push(id)
	defer t.cur.pop()
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
	// func id is unique, comparing the last should be enough
	if last := imin(len(p), len(dst)) - 1; last >= 0 {
		return p[last] == dst[last]
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

func (p path) valid() bool {
	for _, id := range p {
		if !id.valid() {
			return false
		}
	}
	return true
}
