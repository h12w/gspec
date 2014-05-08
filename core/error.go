// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"runtime"
	"sync"

	"github.com/hailiang/gspec/extension"
)

type testError struct {
	err error
	mu  sync.Mutex
}

func (t *testError) setErr(err error) {
	if t.err == nil {
		t.err = err // only keeps the first failure.
	}
}

// get clears the err field so that it will not be repeatedly recorded by
// parent test groups
func (t *testError) getErr() error {
	defer func() { t.err = nil }()
	return t.err
}

// Fail marks that the test case has failed with an error.
func (t *testError) Fail(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.setErr(err)
}

// FailNow marks that the test case has failed with an error, and stops the
// excution of the current goroutine and defer functions will get called.
func (t *testError) FailNow(err error) {
	t.Fail(err)
	runtime.Goexit()
}

func (t *testError) run(f func()) {
	defer func() {
		if e := recover(); e != nil {
			t.setErr(extension.NewPanicError(e, 2))
		}
	}()
	f()
}
