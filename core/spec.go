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
	FailNow(err error)
}

// spec implements "S" interface.
type spec struct {
	*group
	*collector
	testError
}

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func newSpec(g *group, l *collector) S {
	return &spec{
		group:     g,
		collector: l,
	}
}

func (s *spec) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		s.visit(func() {
			s.groupStart(&ext.TestGroup{ID: s.current().String(), Description: name + description}, s.current())
			defer func() {
				s.groupEnd(s.getErr(), s.current())
			}()
			s.capturePanic(f)
		})
	}
}
