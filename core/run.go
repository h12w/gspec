// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
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
