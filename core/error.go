// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"

	"github.com/hailiang/gspec/extension"
)

type testError struct {
	err error
	mu  sync.Mutex
}

func (t *testError) set(err error) {
	if t.err == nil {
		t.err = err // only keeps the first failure.
	}
}

func (t *testError) get() error {
	defer func() { t.err = nil }()
	return t.err
}

// Fail notify that the test case has failed with an error.
func (t *testError) Fail(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.set(err)
}

func (t *testError) run(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = extension.NewPanicError(e, 2)
		}
	}()
	f()
	err = t.get()
	return
}
