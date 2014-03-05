// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	ext "github.com/hailiang/gspec/extension"
)

// TestFunc is the type of the function prepared to run in a goroutine for each
// test case.
type TestFunc func(S)

// S (short for "spec") provides the interface for writing tests and internally
// holds an object that contains minimal context needed to pass into a testing
// goroutine.
type S interface {
	Alias(name string) DescFunc
	Fail(err error)
}

// specImpl implements "S" interface.
type specImpl struct {
	*group
	*listener
	funcCounter
	testError
}

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func newSpec(g *group, l *listener) S {
	return &specImpl{
		group:       g,
		listener:    l,
		funcCounter: newFuncCounter(),
	}
}

func (t *specImpl) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		id := t.funcID(f)
		t.visit(id, func() {
			t.groupStart(&ext.TestGroup{ID: t.dst.String(), Description: name + description}, t.current())
			err := t.run(f)
			t.groupEnd(err, id)
		})
	}
}
