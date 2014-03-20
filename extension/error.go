// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"errors"
	"fmt"
	"runtime"
)

// FuncPos represents a position in a stack trace.
type FuncPos struct {
	Name string
	File string
	Line int
	PC   uintptr
}

// PanicError wraps an error from panicking and its call stack trace.
type PanicError struct {
	Err   error
	Stack []FuncPos
	SS    []byte // TODO: if it is useful, parse it, otherwise, delete it.
}

// NewPanicError returns a new PanicError. Object o is the panicking error,
// skip is the magic number to skip noise entries of stack trace.
func NewPanicError(o interface{}, skip int) error {
	var err error
	switch v := o.(type) {
	case string:
		err = errors.New(v)
	case error:
		err = v
	default:
		err = fmt.Errorf("%v", o)
	}
	pe := &PanicError{err, newStackTrace(skip + 2), make([]byte, 4096)}
	runtime.Stack(pe.SS, true)
	return pe
}

// Error is the same as the panicking error.
func (e *PanicError) Error() string {
	return e.Err.Error()
}

func newStackTrace(skip int) []FuncPos {
	s := []FuncPos{}
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		s = append(s, FuncPos{runtime.FuncForPC(pc).Name(), file, line, pc})
		if !ok {
			break
		}
	}
	return s
}
