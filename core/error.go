// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"

	ext "h12.me/gspec/extension"
)

type failNowError struct {
	error
}

// testError receives the error or panic sent within a test group and transfers
// it to the context when needed.
type testError struct {
	err error
	mu  sync.Mutex
}

func (t *testError) setErr(err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err == nil {
		t.err = err // only keeps the first failure.
	}
}

// get clears the err field so that it will not be repeatedly recorded by
// parent test groups
func (t *testError) getErr() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	defer func() { t.err = nil }()
	return t.err
}

// Fail marks that the test case has failed with an error. It can be called in
// another goroutine.
func (t *testError) Fail(err error) {
	t.setErr(err)
}

// FailNow marks that the test case has failed with an error, and stops the
// excution of the current goroutine and defer functions will get called. It
// must be called in the same goroutine of the test group.
func (t *testError) FailNow(err error) {
	t.Fail(err)
	panic(failNowError{err})
}

func (t *testError) capturePanic(f func()) {
	defer func() {
		if e := recover(); e != nil {
			switch err := e.(type) {
			case failNowError:
				t.setErr(err.error)
			default:
				t.setErr(ext.NewPanicError(e, 2))
			}
		}
	}()
	if f != nil {
		f()
	} else {
		t.setErr(ext.NewPendingError())
	}
}
