// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"errors"
	"fmt"
)

func (t *specImpl) run(f func()) (err error) {
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
	if t.err != nil {
		err = t.err
	}
	return
}

// Fail notify that the test case has failed with an error.
func (t *specImpl) Fail(err error) {
	t.lock()
	defer t.unlock()
	if t.err == nil {
		t.err = err // only keeps the first failure.
	}
}

func (t *specImpl) lock() {
	t.mu.Lock()
}

func (t *specImpl) unlock() {
	t.mu.Unlock()
}
