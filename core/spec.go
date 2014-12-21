// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	ext "h12.me/gspec/extension"
)

// S (short for "spec") provides the interface for writing tests and internally
// holds an object that contains minimal context needed to pass into a testing
// goroutine.
type S interface {
	Alias(name string) DescFunc // Define an alias of a DescFunc.
	Fail(err error)             // Report a failure and continue test execution.
	FailNow(err error)          // Report a failure and stop test execution immediately.
	start()                     // start time measurement.
	end()                       // stop time measurement.
}

// spec implements "S" interface.
type spec struct {
	*group
	*collector
	timer
	testError
}

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func newSpec(g *group, c *collector) S {
	return &spec{
		group:     g,
		collector: c,
		timer:     timer{setDuration: c.setDuration},
	}
}

func (s *spec) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		s.visit(func(cur Path) {
			s.leaf = cur.clone()
			s.groupStart(
				&ext.TestGroup{
					ID:          cur.String(),
					Description: name + description},
				cur,
			)
			defer func() {
				s.groupEnd(s.getErr(), cur)
			}()
			s.capturePanic(f)
		})
	}
}
