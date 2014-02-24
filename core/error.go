// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"sync"
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
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				err = fmt.Errorf("%v", v)
			}
			// TODO: print error, terminate all tests and exit
		}
	}()
	f()
	err = t.get()
	return
}
