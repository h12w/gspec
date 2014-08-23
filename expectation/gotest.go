// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type logger struct {
	mu  sync.Mutex // ensures atomic writes; protects the following fields
	out io.Writer  // destination for output
}

var std = &logger{out: os.Stdout}

// SetOutput sets the output destination for the standard logger, which is used
// by gotest utilities.
func SetOutput(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.out = w
}

// T is a subset of testing.T used in this package.
type T interface {
	Fail()
	FailNow()
}

// TFail return the FailFunc for testing.T.Fail
func TFail(fail func()) FailFunc {
	return func(err error) {
		fmt.Fprintln(std.out, err.Error())
		fail()
	}
}
